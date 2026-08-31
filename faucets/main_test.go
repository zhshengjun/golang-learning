package main

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestAlchemyTransferStatistics(t *testing.T) {
	sender := common.HexToAddress("0x1111111111111111111111111111111111111111")
	recipient := common.HexToAddress("0x2222222222222222222222222222222222222222")
	stats := map[common.Address]*accountStat{recipient: {TotalWei: new(big.Int)}}
	transfer := assetTransfer{From: sender.Hex(), To: recipient.Hex(), Category: "external"}
	transfer.RawContract.Value = "0x14d1120d7b160000"
	transfer.Metadata.BlockTimestamp = "1970-01-01T00:03:20Z"
	if err := recordTransfer(stats[recipient], transfer, sender, recipient); err != nil {
		t.Fatal(err)
	}
	transfer.RawContract.Value = "0x6f05b59d3b20000"
	transfer.Metadata.BlockTimestamp = "1970-01-01T00:01:40Z"
	if err := recordTransfer(stats[recipient], transfer, sender, recipient); err != nil {
		t.Fatal(err)
	}
	if got := weiToETH(stats[recipient].TotalWei); got != "2" {
		t.Fatalf("weiToETH() = %q, want 2", got)
	}
	if stats[recipient].Earliest != 100 || stats[recipient].Latest != 200 {
		t.Fatalf("time range = %d..%d, want 100..200", stats[recipient].Earliest, stats[recipient].Latest)
	}
	if got := formatTime(100); got != "1970-01-01 08:01:40" {
		t.Fatalf("formatTime() = %q", got)
	}
	if got := estimatedMiningTime(10_000, 5_000); got != 8_200 {
		t.Fatalf("estimatedMiningTime() = %d, want 8200", got)
	}
	if got := estimatedMiningTime(6_000, 5_000); got != 5_000 {
		t.Fatalf("estimatedMiningTime() = %d, want now", got)
	}
}

func TestClaimAvailability(t *testing.T) {
	now := uint64(1_000_000)
	empty := &accountStat{TotalWei: new(big.Int)}
	if availableAt, total := claimAvailability(empty, now); availableAt != now || total.Sign() != 0 {
		t.Fatalf("empty account = %d, %s; want now, 0", availableAt, total)
	}

	blockedAddress := common.HexToAddress("0x1111111111111111111111111111111111111111")
	readyAddress := common.HexToAddress("0x2222222222222222222222222222222222222222")
	blocked := &accountStat{TotalWei: new(big.Int)}
	blocked.record(new(big.Int).Set(singleClaimWei), now-40*60*60)
	blocked.record(new(big.Int).Set(singleClaimWei), now-20*60*60)
	ready := &accountStat{TotalWei: new(big.Int)}
	ready.record(new(big.Int).Set(singleClaimWei), now-10*60*60)

	if availableAt, total := claimAvailability(blocked, now); availableAt != now+8*60*60 || total.Cmp(claimLimitWei) != 0 {
		t.Fatalf("blocked account = %d, %s", availableAt, total)
	}
	stats := map[common.Address]*accountStat{blockedAddress: blocked, readyAddress: ready}
	accounts := sortAccountsByAvailability([]string{blockedAddress.Hex(), readyAddress.Hex()}, stats, now)
	if common.HexToAddress(accounts[0]) != readyAddress || common.HexToAddress(accounts[1]) != blockedAddress {
		t.Fatalf("sortAccountsByAvailability() = %v", accounts)
	}
}
