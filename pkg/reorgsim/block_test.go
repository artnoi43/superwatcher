package reorgsim

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/artnoi43/gsl"
	"github.com/ethereum/go-ethereum/common"
)

func TestReorg(t *testing.T) {
	var blockNumber uint64 = 15944408
	logs := InitMappedLogsFromFiles(defaultLogsFiles...)
	blockLogs := logs[blockNumber]
	oldLogsByTxHash := mapLogsToTxHash(blockLogs)
	fmt.Println("oldLogs by TxHash")
	prontMapLen(oldLogsByTxHash, "txHash", "len(logs)")

	b := Block{
		blockNumber: blockNumber,
		hash:        common.BigToHash(big.NewInt(69)),
		logs:        blockLogs,
		reorgedHere: false,
		toBeForked:  true,
	}

	_b := b.reorg(0)
	newLogsByTxHash := mapLogsToTxHash(_b.logs)
	fmt.Println("newLogs by TxHash")
	prontMapLen(newLogsByTxHash, "txHash", "len(logs)")
}

func TestRemoveLogs(t *testing.T) {
	var reorgedAt uint64 = 15944450
	var blockNumber uint64 = 15944455 // 0xf34b07
	// var moveTo uint64 = 15944498   // 0xf34b32
	hashesToRemove := []common.Hash{
		common.HexToHash("0x620be69b041f986127322985854d3bc785abe1dc9f4df49173409f15b7515164"),
	}

	logs := InitMappedLogsFromFiles(defaultLogsFiles...)
	chain := newBlockChain(logs, reorgedAt)

	b := chain[blockNumber]

	// Test in test, because test log files may have been changed
	{
		foundLogs := make([]bool, len(hashesToRemove))
		var c int
		for _, log := range b.logs {
			if gsl.Contains(hashesToRemove, log.TxHash) {
				foundLogs[c] = true
				c++
			}
		}

		for _, foundLog := range foundLogs {
			if !foundLog {
				t.Skip("txHashes not found - probably bad hard-coded logs")
			}
		}
	}

	b.removeLogs(hashesToRemove)

	for _, log := range b.logs {
		if gsl.Contains(hashesToRemove, log.TxHash) {
			t.Fatalf("removed log %s was not removed", log.TxHash.String())
		}
	}
}
