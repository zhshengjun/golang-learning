package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.yaml.in/yaml/v3"
)

var faucetsChainID = big.NewInt(11155111)
var utc8 = time.FixedZone("UTC+8", 8*60*60)

const twoDayWindowSeconds = uint64(2 * 24 * 60 * 60)
const sevenDayWindowSeconds = uint64(7 * 24 * 60 * 60)

var sessionOffsetSeconds uint64

var singleClaimWei = big.NewInt(2_500_000_000_000_000_000)
var twoDayLimitWei = big.NewInt(5_000_000_000_000_000_000)
var sevenDayLimitWei = new(big.Int).Mul(big.NewInt(10), big.NewInt(1_000_000_000_000_000_000))

type config struct {
	Raw struct {
		URL    string `yaml:"url"`
		APIKey string `yaml:"api_key"`
	} `yaml:"raw"`
	Account struct {
		From struct {
			Address string `yaml:"address"`
		} `yaml:"from"`
	} `yaml:"account"`
	Faucets struct {
		Accounts             []string `yaml:"accounts"`
		Observe              []string `yaml:"observe"`
		SessionOffsetSeconds uint64   `yaml:"session_offset_seconds"`
	} `yaml:"faucets"`
}

type accountStat struct {
	TotalWei  *big.Int
	Count     uint64
	Earliest  uint64
	Latest    uint64
	Transfers []claimTransfer
}

type claimTransfer struct {
	Value     *big.Int
	Timestamp uint64
}

type transferParams struct {
	FromBlock        string   `json:"fromBlock"`
	ToBlock          string   `json:"toBlock"`
	FromAddress      string   `json:"fromAddress"`
	ToAddress        string   `json:"toAddress"`
	Category         []string `json:"category"`
	ExcludeZeroValue bool     `json:"excludeZeroValue"`
	WithMetadata     bool     `json:"withMetadata"`
	MaxCount         string   `json:"maxCount"`
	Order            string   `json:"order"`
	PageKey          string   `json:"pageKey,omitempty"`
}

type transferResult struct {
	Transfers []assetTransfer `json:"transfers"`
	PageKey   string          `json:"pageKey"`
}

type assetTransfer struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Category    string `json:"category"`
	RawContract struct {
		Value string `json:"value"`
	} `json:"rawContract"`
	Metadata struct {
		BlockTimestamp string `json:"blockTimestamp"`
	} `json:"metadata"`
}

func main() {
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	if err := run(*configPath); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}

func run(configPath string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	sessionOffsetSeconds = cfg.Faucets.SessionOffsetSeconds

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := ethclient.DialContext(ctx, cfg.Raw.URL)
	if err != nil {
		return fmt.Errorf("连接 RPC: %w", err)
	}
	defer client.Close()

	actualChainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("查询 chain ID: %w", err)
	}
	if actualChainID.Cmp(faucetsChainID) != 0 {
		return fmt.Errorf("RPC 不是 Sepolia: chain ID = %s", actualChainID)
	}

	latest, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("查询最新区块: %w", err)
	}
	cutoff := uint64(0)
	if sevenDayWindowSeconds < latest.Time {
		cutoff = latest.Time - sevenDayWindowSeconds
	}
	start, err := firstBlockSince(ctx, client, latest.Number.Uint64(), cutoff)
	if err != nil {
		return err
	}

	sender := common.HexToAddress(cfg.Account.From.Address)
	stats := make(map[common.Address]*accountStat, len(cfg.Faucets.Accounts))
	balances := make(map[common.Address]*big.Int, len(cfg.Faucets.Accounts))
	for _, value := range cfg.Faucets.Accounts {
		address := common.HexToAddress(value)
		stats[address] = &accountStat{TotalWei: new(big.Int)}
		balance, err := client.BalanceAt(ctx, address, latest.Number)
		if err != nil {
			return fmt.Errorf("查询账户 %s 余额: %w", address, err)
		}
		balances[address] = balance
	}

	for _, value := range cfg.Faucets.Accounts {
		recipient := common.HexToAddress(value)
		if err := fetchTransfers(ctx, client, start, latest.Number.Uint64(), sender, recipient, stats[recipient]); err != nil {
			return err
		}
	}

	now := latest.Time
	fmt.Printf("链上到账时间按 session 创建时间提前 %d min 估算\n", sessionOffsetSeconds/60)
	fmt.Println("基于 Alchemy 链上数据的可挖矿顺序:")
	for index, value := range sortAccountsByAvailability(cfg.Faucets.Accounts, stats, now) {
		address := common.HexToAddress(value)
		stat := stats[address]
		availableAt, last48Hours, last7Days := claimAvailability(stat, now)
		fmt.Printf("%d. %s  余额=%s ETH  预计可开始时间=%s  2d=%s ETH  7d=%s ETH\n",
			index+1, address, weiToETH(balances[address]), formatAvailability(availableAt, now), weiToETH(last48Hours), weiToETH(last7Days))
	}
	if len(cfg.Faucets.Observe) > 0 {
		fmt.Println("观察钱包:")
	}
	for index, value := range cfg.Faucets.Observe {
		address := common.HexToAddress(value)
		balance, err := client.BalanceAt(ctx, address, latest.Number)
		if err != nil {
			return fmt.Errorf("查询观察钱包 %s 余额: %w", address, err)
		}
		fmt.Printf("%d. %s  余额=%s ETH\n", index+1, address, weiToETH(balance))
	}
	return nil
}

func loadConfig(path string) (config, error) {
	var cfg config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("读取配置: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("解析配置: %w", err)
	}
	if cfg.Raw.URL == "" {
		return cfg, errors.New("raw.url 不能为空")
	}
	if cfg.Raw.APIKey == "" {
		return cfg, errors.New("raw.api_key 不能为空")
	}
	cfg.Raw.URL = strings.TrimRight(cfg.Raw.URL, "/") + "/" + cfg.Raw.APIKey
	if !common.IsHexAddress(cfg.Account.From.Address) {
		return cfg, errors.New("account.from.address 不是有效地址")
	}
	if len(cfg.Faucets.Accounts) == 0 {
		return cfg, errors.New("faucets.accounts 不能为空")
	}
	if cfg.Faucets.SessionOffsetSeconds == 0 {
		return cfg, errors.New("faucets.session_offset_seconds 必须大于 0")
	}
	for _, account := range cfg.Faucets.Accounts {
		if !common.IsHexAddress(account) {
			return cfg, fmt.Errorf("faucets.accounts 包含无效地址: %q", account)
		}
	}
	for _, account := range cfg.Faucets.Observe {
		if !common.IsHexAddress(account) {
			return cfg, fmt.Errorf("faucets.observe 包含无效地址: %q", account)
		}
	}
	return cfg, nil
}

func firstBlockSince(ctx context.Context, client *ethclient.Client, latest, cutoff uint64) (uint64, error) {
	low, high := uint64(0), latest
	for low < high {
		middle := low + (high-low)/2
		header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(middle))
		if err != nil {
			return 0, fmt.Errorf("查询区块 %d 时间: %w", middle, err)
		}
		if header.Time < cutoff {
			low = middle + 1
		} else {
			high = middle
		}
	}
	return low, nil
}

func fetchTransfers(ctx context.Context, client *ethclient.Client, start, end uint64, sender, recipient common.Address, stat *accountStat) error {
	params := transferParams{
		FromBlock: hexutil.EncodeUint64(start), ToBlock: hexutil.EncodeUint64(end),
		FromAddress: sender.Hex(), ToAddress: recipient.Hex(), Category: []string{"external"},
		ExcludeZeroValue: true, WithMetadata: true, MaxCount: "0x3e8", Order: "asc",
	}
	for {
		var result transferResult
		if err := client.Client().CallContext(ctx, &result, "alchemy_getAssetTransfers", params); err != nil {
			return fmt.Errorf("查询账户 %s 的 Alchemy 转账: %w", recipient, err)
		}
		for _, transfer := range result.Transfers {
			if err := recordTransfer(stat, transfer, sender, recipient); err != nil {
				return err
			}
		}
		if result.PageKey == "" {
			return nil
		}
		params.PageKey = result.PageKey
	}
}

func recordTransfer(stat *accountStat, transfer assetTransfer, sender, recipient common.Address) error {
	if transfer.Category != "external" || !strings.EqualFold(transfer.From, sender.Hex()) || !strings.EqualFold(transfer.To, recipient.Hex()) {
		return errors.New("alchemy 返回了不符合查询条件的转账")
	}
	wei, err := hexutil.DecodeBig(transfer.RawContract.Value)
	if err != nil || wei.Sign() <= 0 {
		return fmt.Errorf("alchemy 返回了无效的转账金额 %q", transfer.RawContract.Value)
	}
	timestamp, err := time.Parse(time.RFC3339, transfer.Metadata.BlockTimestamp)
	if err != nil {
		return fmt.Errorf("alchemy 返回了无效的区块时间 %q: %w", transfer.Metadata.BlockTimestamp, err)
	}
	stat.record(wei, uint64(timestamp.Unix()))
	return nil
}

func (stat *accountStat) record(value *big.Int, timestamp uint64) {
	stat.TotalWei.Add(stat.TotalWei, value)
	stat.Count++
	stat.Transfers = append(stat.Transfers, claimTransfer{Value: new(big.Int).Set(value), Timestamp: timestamp})
	if stat.Earliest == 0 || timestamp < stat.Earliest {
		stat.Earliest = timestamp
	}
	if timestamp > stat.Latest {
		stat.Latest = timestamp
	}
}

func sortAccountsByAvailability(accounts []string, stats map[common.Address]*accountStat, now uint64) []string {
	sorted := append([]string(nil), accounts...)
	sort.SliceStable(sorted, func(i, j int) bool {
		left, left48Hours, _ := claimAvailability(stats[common.HexToAddress(sorted[i])], now)
		right, right48Hours, _ := claimAvailability(stats[common.HexToAddress(sorted[j])], now)
		if left != right {
			return left < right
		}
		if left48Hours.Sign() == 0 && right48Hours.Sign() != 0 {
			return true
		}
		if left48Hours.Sign() != 0 && right48Hours.Sign() == 0 {
			return false
		}
		return false
	})
	return sorted
}

func claimAvailability(stat *accountStat, now uint64) (uint64, *big.Int, *big.Int) {
	available48Hours, total48Hours := windowAvailability(stat, now, twoDayWindowSeconds, twoDayLimitWei)
	available7Days, total7Days := windowAvailability(stat, now, sevenDayWindowSeconds, sevenDayLimitWei)
	return max(available48Hours, available7Days), total48Hours, total7Days
}

func windowAvailability(stat *accountStat, now, window uint64, limit *big.Int) (uint64, *big.Int) {
	recent, total := transfersWithin(stat, now, window)
	if canClaim(total, limit) {
		return now, total
	}
	sort.Slice(recent, func(i, j int) bool { return recent[i].Timestamp < recent[j].Timestamp })
	remaining := new(big.Int).Set(total)
	for index := 0; index < len(recent); {
		timestamp := recent[index].Timestamp
		for index < len(recent) && recent[index].Timestamp == timestamp {
			remaining.Sub(remaining, recent[index].Value)
			index++
		}
		if canClaim(remaining, limit) {
			return timestamp + window, total
		}
	}
	return now, total
}

func transfersWithin(stat *accountStat, now, window uint64) ([]claimTransfer, *big.Int) {
	cutoff := uint64(0)
	if now > window {
		cutoff = now - window
	}
	recent := make([]claimTransfer, 0, len(stat.Transfers))
	total := new(big.Int)
	for _, transfer := range stat.Transfers {
		if transfer.Timestamp > sessionOffsetSeconds {
			timestamp := transfer.Timestamp - sessionOffsetSeconds
			if timestamp <= cutoff {
				continue
			}
			recent = append(recent, claimTransfer{Value: transfer.Value, Timestamp: timestamp})
			total.Add(total, transfer.Value)
		}
	}
	return recent, total
}

func canClaim(current, limit *big.Int) bool {
	return new(big.Int).Add(new(big.Int).Set(current), singleClaimWei).Cmp(limit) <= 0
}

func formatTime(timestamp uint64) string {
	if timestamp == 0 {
		return "无"
	}
	return time.Unix(int64(timestamp), 0).In(utc8).Format("2006-01-02 15:04:05")
}

func formatAvailability(timestamp, now uint64) string {
	if timestamp <= now {
		return "现在"
	}
	return formatTime(timestamp)
}

func weiToETH(wei *big.Int) string {
	whole, fraction := new(big.Int), new(big.Int)
	whole.QuoRem(wei, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), fraction)
	if fraction.Sign() == 0 {
		return whole.String()
	}
	digits := fraction.String()
	return whole.String() + "." + strings.TrimRight(strings.Repeat("0", 18-len(digits))+digits, "0")
}
