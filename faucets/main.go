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

const claimWindowSeconds = uint64(48 * 60 * 60)

const miningLeadMins = uint64(0)
const miningLeadSeconds = miningLeadMins * 60

var singleClaimWei = big.NewInt(2_500_000_000_000_000_000)
var claimLimitWei = big.NewInt(5_000_000_000_000_000_000)

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
		Accounts []string `yaml:"accounts"`
		Interval uint64   `yaml:"interval"`
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
	if cfg.Faucets.Interval < latest.Time {
		cutoff = latest.Time - cfg.Faucets.Interval
	}
	start, err := firstBlockSince(ctx, client, latest.Number.Uint64(), cutoff)
	if err != nil {
		return err
	}

	sender := common.HexToAddress(cfg.Account.From.Address)
	stats := make(map[common.Address]*accountStat, len(cfg.Faucets.Accounts))
	for _, value := range cfg.Faucets.Accounts {
		stats[common.HexToAddress(value)] = &accountStat{TotalWei: new(big.Int)}
	}

	for _, value := range cfg.Faucets.Accounts {
		recipient := common.HexToAddress(value)
		if err := fetchTransfers(ctx, client, start, latest.Number.Uint64(), sender, recipient, stats[recipient]); err != nil {
			return err
		}
	}

	now := latest.Time
	fmt.Printf("滚动时间以 session 创建为准，统计会累积在领取的 IP 上，当前时间是提前 %d min\n", miningLeadMins)
	fmt.Println("下面是基于 Alchemy 链上数据估算可挖矿时间:")
	for index, value := range sortAccountsByAvailability(cfg.Faucets.Accounts, stats, now) {
		address := common.HexToAddress(value)
		stat := stats[address]
		availableAt, last48Hours := claimAvailability(stat, now)
		fmt.Printf("%d. %s  预计可开始时间=%s  48h=%s ETH  72h=%s ETH\n",
			index+1, address, formatAvailability(availableAt, now), weiToETH(last48Hours), weiToETH(stat.TotalWei))
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
	if cfg.Faucets.Interval == 0 {
		return cfg, errors.New("faucets.interval 必须大于 0")
	}
	if cfg.Faucets.Interval < claimWindowSeconds {
		return cfg, errors.New("faucets.interval 不能小于 172800 秒")
	}
	if len(cfg.Faucets.Accounts) == 0 {
		return cfg, errors.New("faucets.accounts 不能为空")
	}
	for _, account := range cfg.Faucets.Accounts {
		if !common.IsHexAddress(account) {
			return cfg, fmt.Errorf("faucets.accounts 包含无效地址: %q", account)
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
		left, _ := claimAvailability(stats[common.HexToAddress(sorted[i])], now)
		right, _ := claimAvailability(stats[common.HexToAddress(sorted[j])], now)
		return estimatedMiningTime(left, now) < estimatedMiningTime(right, now)
	})
	return sorted
}

func claimAvailability(stat *accountStat, now uint64) (uint64, *big.Int) {
	cutoff := uint64(0)
	if now > claimWindowSeconds {
		cutoff = now - claimWindowSeconds
	}
	recent := make([]claimTransfer, 0, len(stat.Transfers))
	total := new(big.Int)
	for _, transfer := range stat.Transfers {
		if transfer.Timestamp > cutoff {
			recent = append(recent, transfer)
			total.Add(total, transfer.Value)
		}
	}
	if canClaim(total) {
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
		if canClaim(remaining) {
			return timestamp + claimWindowSeconds, total
		}
	}
	return now, total
}

func canClaim(current *big.Int) bool {
	return new(big.Int).Add(new(big.Int).Set(current), singleClaimWei).Cmp(claimLimitWei) <= 0
}

func formatTime(timestamp uint64) string {
	if timestamp == 0 {
		return "无"
	}
	return time.Unix(int64(timestamp), 0).In(utc8).Format("2006-01-02 15:04:05")
}

func formatAvailability(timestamp, now uint64) string {
	timestamp = estimatedMiningTime(timestamp, now)
	if timestamp == now {
		return "现在"
	}
	return formatTime(timestamp)
}

func estimatedMiningTime(timestamp, now uint64) uint64 {
	if timestamp <= now+miningLeadSeconds {
		return now
	}
	return timestamp - miningLeadSeconds
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
