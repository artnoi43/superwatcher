package reorgsim

import (
	"fmt"
	"testing"
)

func TestReorg(t *testing.T) {
	var blockNumber uint64 = 15944408
	logs := InitLogs()
	blockLogs := logs[blockNumber]
	oldLogsByTxHash := mapLogsToTxHash(blockLogs)
	fmt.Println("oldLogs by TxHash")
	prontMapLen(oldLogsByTxHash, "txHash", "len(logs)")

	b := block{
		blockNumber: blockNumber,
		hash:        randomHash(70),
		logs:        blockLogs,
		reorgedHere: false,
		toBeForked:  true,
	}

	_b := b.reorg()
	newLogsByTxHash := mapLogsToTxHash(_b.logs)
	fmt.Println("newLogs by TxHash")
	prontMapLen(newLogsByTxHash, "txHash", "len(logs)")
}
