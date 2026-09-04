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
}

func TestClaimAvailability(t *testing.T) {
	previousOffset := sessionOffsetSeconds
	sessionOffsetSeconds = 60 * 60
	t.Cleanup(func() { sessionOffsetSeconds = previousOffset })

	now := uint64(1_000_000)
	empty := &accountStat{TotalWei: new(big.Int)}
	if availableAt, total48Hours, total7Days := claimAvailability(empty, now); availableAt != now || total48Hours.Sign() != 0 || total7Days.Sign() != 0 {
		t.Fatalf("empty account = %d, %s, %s; want now, 0, 0", availableAt, total48Hours, total7Days)
	}

	blockedAddress := common.HexToAddress("0x1111111111111111111111111111111111111111")
	readyAddress := common.HexToAddress("0x2222222222222222222222222222222222222222")
	blocked := &accountStat{TotalWei: new(big.Int)}
	blocked.record(new(big.Int).Set(singleClaimWei), now-40*60*60)
	blocked.record(new(big.Int).Set(singleClaimWei), now-20*60*60)
	ready := &accountStat{TotalWei: new(big.Int)}
	ready.record(new(big.Int).Set(singleClaimWei), now-10*60*60)

	if availableAt, total48Hours, _ := claimAvailability(blocked, now); availableAt != now+7*60*60 || total48Hours.Cmp(twoDayLimitWei) != 0 {
		t.Fatalf("blocked account = %d, %s", availableAt, total48Hours)
	}
	weeklyBlocked := &accountStat{TotalWei: new(big.Int)}
	for daysAgo := uint64(3); daysAgo <= 6; daysAgo++ {
		weeklyBlocked.record(new(big.Int).Set(singleClaimWei), now-daysAgo*24*60*60)
	}
	if availableAt, total48Hours, total7Days := claimAvailability(weeklyBlocked, now); availableAt != now+23*60*60 || total48Hours.Sign() != 0 || total7Days.Cmp(sevenDayLimitWei) != 0 {
		t.Fatalf("weekly blocked account = %d, %s, %s", availableAt, total48Hours, total7Days)
	}
	stats := map[common.Address]*accountStat{blockedAddress: blocked, readyAddress: ready}
	accounts := sortAccountsByAvailability([]string{blockedAddress.Hex(), readyAddress.Hex()}, stats, now)
	if common.HexToAddress(accounts[0]) != readyAddress || common.HexToAddress(accounts[1]) != blockedAddress {
		t.Fatalf("sortAccountsByAvailability() = %v", accounts)
	}
	weeklyAddress := common.HexToAddress("0x3333333333333333333333333333333333333333")
	stats[weeklyAddress] = weeklyBlocked
	accounts = sortAccountsByAvailability([]string{readyAddress.Hex(), weeklyAddress.Hex()}, stats, now)
	if common.HexToAddress(accounts[0]) != readyAddress {
		t.Fatalf("earlier account was not first: %v", accounts)
	}
	emptyAddress := common.HexToAddress("0x4444444444444444444444444444444444444444")
	stats[emptyAddress] = empty
	accounts = sortAccountsByAvailability([]string{readyAddress.Hex(), emptyAddress.Hex()}, stats, now)
	if common.HexToAddress(accounts[0]) != emptyAddress {
		t.Fatalf("48h unclaimed account did not win equal-time tie: %v", accounts)
	}
	shifted := &accountStat{TotalWei: new(big.Int)}
	shifted.record(new(big.Int).Set(singleClaimWei), now-twoDayWindowSeconds+30*60)
	_, total48Hours := transfersWithin(shifted, now, twoDayWindowSeconds)
	if total48Hours.Sign() != 0 {
		t.Fatalf("session offset did not expire 2d transfer: %s", total48Hours)
	}
}
