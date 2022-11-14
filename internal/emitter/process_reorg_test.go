package emitter

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

func TestProcessReorg(t *testing.T) {
	tracker := newTracker()
	hardcodedLogs := reorgsim.InitLogs()
	reorgedAt := uint64(15944444)
	oldChain, reorgedChain := reorgsim.NewBlockChain(hardcodedLogs, reorgedAt)

	// Add oldChain's blocks to tracker
	for blockNumber, block := range oldChain {
		blockLogs := block.Logs()
		logs := gslutils.CollectPointers(&blockLogs)

		b := new(superwatcher.BlockInfo)
		b.Logs = logs
		b.Hash = block.Hash()
		b.Number = blockNumber

		tracker.addTrackerBlock(b)
	}

	var reorgedLogs []types.Log
	var reorgedHeader = make(map[uint64]superwatcher.BlockHeader)
	for blockNumber, block := range reorgedChain {
		if logs := block.Logs(); len(logs) != 0 {
			reorgedLogs = append(reorgedLogs, logs...)
		}
		reorgedHeader[blockNumber] = block
	}

	freshHashes, freshLogs, processLogs := PopulateInitialMaps(reorgedLogs, reorgedHeader)

	wasReorged := processReorged(
		tracker,
		15944444,
		15950000,
		freshHashes,
		freshLogs,
		processLogs,
	)

	t.Log(wasReorged)
}
