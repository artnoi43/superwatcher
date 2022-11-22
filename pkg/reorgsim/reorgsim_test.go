package reorgsim

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	startBlock  uint64 = 15900000
	reorgedAt   uint64 = 15944444
	defaultLogs        = []string{
		"./assets/logs_poolfactory.json",
		"./assets/logs_lp.json",
	}
)

func initChains(reorgedAt uint64) (blockChain, blockChain) {
	mappedLogs := InitLogs(defaultLogs)
	return NewBlockChain(mappedLogs, reorgedAt)
}

// Test if NewBlockChain works properly
func TestNewBlockChain(t *testing.T) {
	oldChain, reorgedChain := initChains(reorgedAt)

	for blockNumber, reorgedBlock := range reorgedChain {
		oldBlock := oldChain[blockNumber]

		oldLogs := oldBlock.Logs()
		reorgedLogs := reorgedBlock.Logs()

		if lo, lr := len(oldLogs), len(reorgedLogs); lo != lr {
			t.Fatalf("len(logs) not match on block %d", blockNumber)
		}

		if !reorgedBlock.toBeForked {
			continue
		}

		if oldBlock.Hash() == reorgedBlock.Hash() {
			t.Fatal("old and reorg block hashes match")
		}

		if blockNumber < reorgedAt && reorgedBlock.toBeForked {
			t.Fatal("unreorged block from reorgedChain tagged with toBeForked")
		}

		if blockNumber > reorgedAt && !reorgedBlock.toBeForked {
			t.Fatal("reorgedBlock not tagged with toBeForked")
		}

		for i, reorgedLog := range reorgedLogs {
			oldLog := oldLogs[i]

			if reorgedLog.TxHash == oldLog.TxHash {
				t.Fatal("old and reorg log txHash match")
			}

			if reorgedLog.BlockHash == oldLog.BlockHash {
				t.Fatal("old and reorg log blockHash match")
			}
		}
	}
}

func TestFoo(t *testing.T) {
	chain, reorgedChain := initChains(reorgedAt)

	fmt.Println("old chain")
	prontBlockChain(chain)

	fmt.Println("reorged chain")
	prontBlockChain(reorgedChain)

	param := ReorgParam{
		StartBlock:    startBlock,
		BlockProgress: 3,
		ReorgedAt:     reorgedAt,
		ExitBlock:     reorgedAt + 100,
	}

	sim := NewReorgSim(param, defaultLogs)
	filterLogs, err := sim.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(15944401),
		ToBlock:   big.NewInt(15944500),
	})

	if err != nil {
		t.Fatal(err.Error())
	}
	filterLogsMapped := mapLogsToNumber(filterLogs)
	fmt.Println("filterLogs")
	prontMapLen(filterLogsMapped, "blockNumber", "len(logs)")
}

func prontMapLen[T comparable, U any](m map[T][]U, keyString, lenString string) {
	for k, arr := range m {
		fmt.Println(keyString, k, lenString, len(arr))
	}
}

func prontLogs(logs []types.Log) {
	for _, log := range logs {
		fmt.Println("blockNumber", log.BlockNumber, "addr", log.Address)
	}
}

func prontBlockChain(chain blockChain) {
	for _, b := range chain {
		fmt.Println(
			"blockNumber", b.blockNumber,
			"blockhash", b.Hash().String(),
			"len(logs)", len(b.logs),
			"forked", b.toBeForked,
			"reorgedHere", b.reorgedHere,
		)
	}
}
