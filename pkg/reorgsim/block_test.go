package reorgsim

import (
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
)

func TestReorg(t *testing.T) {
	var blockNumber uint64 = 15944408
	logs := InitMappedLogsFromFiles(defaultLogsFiles...)
	blockLogs := logs[blockNumber]
	oldLogsByTxHash := mapLogsToTxHash(blockLogs)
	fmt.Println("oldLogs by TxHash")
	prontMapLen(oldLogsByTxHash, "txHash", "len(logs)")

	b := block{
		blockNumber: blockNumber,
		hash:        RandomHash(70),
		logs:        blockLogs,
		reorgedHere: false,
		toBeForked:  true,
	}

	_b := b.reorg()
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

	var foundLogs bool
	for _, log := range b.logs {
		if gslutils.Contains(hashesToRemove, log.TxHash) {
			foundLogs = true
			break
		}
	}

	if !foundLogs {
		t.Skip("txHashes not found - probably bad hard-coded logs")
	}

	b.removeLogs(hashesToRemove)

	for _, log := range b.logs {
		if gslutils.Contains(hashesToRemove, log.TxHash) {
			t.Fatalf("removed log %s was not removed", log.TxHash.String())
		}
	}
}
