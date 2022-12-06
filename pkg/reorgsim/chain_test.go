package reorgsim

import (
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
)

func TestReorgMoveLogs(t *testing.T) {
	var reorgedAt uint64 = 15944450
	var moveFrom uint64 = 15944455 // 0xf34b07
	var moveTo uint64 = 15944498   // 0xf34b32

	hashesToRemove := []common.Hash{
		common.HexToHash("0x620be69b041f986127322985854d3bc785abe1dc9f4df49173409f15b7515164"),
	}

	logs := InitMappedLogsFromFiles(defaultLogs)
	_, reorgedChain := NewBlockChain(logs, reorgedAt)

	reorgedChain.reorgMoveLogs(map[uint64][]MoveLogs{
		moveFrom: {
			{
				newBlock: moveTo,
				txHashes: hashesToRemove,
			},
		},
	})

	moveFromBlock := reorgedChain[moveFrom]
	for _, log := range moveFromBlock.logs {
		if gslutils.Contains(hashesToRemove, log.TxHash) {
			t.Log("moveFromBlock.logs", moveFromBlock.logs)
			t.Fatalf("log %s was not removed", log.TxHash.String())
		}
	}

	moveToBlock := reorgedChain[moveTo]
	var foundMovedLog bool
	for _, log := range moveToBlock.logs {
		if gslutils.Contains(hashesToRemove, log.TxHash) {
			foundMovedLog = true
		}
	}

	if !foundMovedLog {
		t.Log("moveToBlock.logs", moveToBlock.logs)
		t.Fatal("moved logs not found in moveToBlock")
	}
}
