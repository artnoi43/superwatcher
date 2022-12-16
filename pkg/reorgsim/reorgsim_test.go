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

func initDefaultChains(reorgedAt uint64) (blockChain, blockChain) {
	return NewBlockChain(reorgedAt, InitLogsFromFiles(defaultLogs...)...)
}

func TestNewBlockChainNg(t *testing.T) {
	oldChain, reorgedChain := initDefaultChains(reorgedAt)
	if err := testBlockChain(t, oldChain, reorgedChain); err != nil {
		t.Fatal(err.Error())
	}
}

// Test if NewBlockChain works properly
func TestNewBlockChain(t *testing.T) {
	oldChain, reorgedChain := initDefaultChains(reorgedAt)
	if err := testBlockChain(t, oldChain, reorgedChain); err != nil {
		t.Fatal(err.Error())
	}
}

func testBlockChain(t *testing.T, oldChain, reorgedChain blockChain) error {
	for blockNumber, reorgedBlock := range reorgedChain {
		oldBlock := oldChain[blockNumber]

		oldLogs := oldBlock.Logs()
		reorgedLogs := reorgedBlock.Logs()

		if lo, lr := len(oldLogs), len(reorgedLogs); lo != lr {
			return fmt.Errorf("len(logs) not match on block %d", blockNumber)
		}

		if !reorgedBlock.toBeForked {
			continue
		}

		if oldBlock.Hash() == reorgedBlock.Hash() {
			return fmt.Errorf("old and reorg block hashes match on block %d:%s", blockNumber, oldBlock.Hash().String())
		}

		if blockNumber < reorgedAt && reorgedBlock.toBeForked {
			return fmt.Errorf("unreorged block %d from reorgedChain tagged with toBeForked", blockNumber)
		}

		if blockNumber > reorgedAt && !reorgedBlock.toBeForked {
			return fmt.Errorf("reorgedBlock %d not tagged with toBeForked", blockNumber)
		}

		for i, reorgedLog := range reorgedLogs {
			oldLog := oldLogs[i]

			// Uncomment to change txHash when reorg too
			// if reorgedLog.TxHash == oldLog.TxHash {
			// 	t.Fatal("old and reorg log txHash match")
			// }

			if reorgedLog.BlockHash == oldLog.BlockHash {
				return fmt.Errorf("old and reorg log blockHash match %d:%s", blockNumber, reorgedLog.BlockHash.String())
			}
		}
	}

	return nil
}

func TestFoo(t *testing.T) {
	reorgedAt := uint64(15944408)
	chain, reorgedChain := initDefaultChains(reorgedAt)

	fmt.Println("old chain")
	prontBlockChain(chain)

	fmt.Println("reorged chain")
	prontBlockChain(reorgedChain)

	param := ParamV1{
		BaseParam: BaseParam{
			StartBlock:    startBlock,
			BlockProgress: 3,
			ExitBlock:     reorgedAt + 100,
		},
		ReorgEvent: ReorgEvent{
			ReorgedBlock: reorgedAt,
		},
	}

	sim := NewReorgSimFromLogsFiles(param, defaultLogs, 3)
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
		fmt.Println("blockNumber", log.BlockNumber, "blockHash", log.BlockHash.String(), "txHash", log.TxHash.String())
	}
}

func prontBlockChain(chain blockChain) {
	for _, b := range chain {
		fmt.Println(
			"blockNumber", b.blockNumber,
			"blockhash", b.Hash().String(),
			"len(logs)", len(b.logs),
			"forked", b.toBeForked,
		)
		prontLogs(b.logs)
	}
}
