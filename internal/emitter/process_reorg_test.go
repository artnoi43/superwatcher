package emitter

import (
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

var (
	fromBlock   uint64 = 15944400
	toBlock     uint64 = 15944500
	reorgedAt   uint64 = 15944444
	defaultLogs        = []string{
		"./assets/logs_poolfactory.json",
		"./assets/logs_lp.json",
	}
)

func TestProcessReorg(t *testing.T) {
	tracker := newTracker()
	hardcodedLogs := reorgsim.InitLogs(defaultLogs)
	oldChain, reorgedChain := reorgsim.NewBlockChain(hardcodedLogs, reorgedAt)

	// Add oldChain's blocks to tracker
	for blockNumber, block := range oldChain {
		blockLogs := block.Logs()
		logs := gslutils.CollectPointers(&blockLogs)

		tracker.addTrackerBlock(&superwatcher.BlockInfo{
			Logs:   logs,
			Hash:   block.Hash(),
			Number: blockNumber,
		})
	}

	var reorgedLogs []types.Log
	var reorgedHeader = make(map[uint64]superwatcher.BlockHeader)
	for blockNumber, block := range reorgedChain {
		if logs := block.Logs(); len(logs) != 0 {
			reorgedLogs = append(reorgedLogs, logs...)
		}
		reorgedHeader[blockNumber] = block
	}

	freshHashes, freshLogs, processLogs := populateInitialMaps(reorgedLogs, reorgedHeader)

	wasReorged := processReorged(
		tracker,
		fromBlock,
		toBlock,
		freshHashes,
		freshLogs,
		processLogs,
	)

	for blockNumber, reorged := range wasReorged {
		if blockNumber >= reorgedAt {
			if !reorged {
				t.Fatalf(
					"blockNumber %d is after reorg block at %d, but it was not tagged \"true\" in wasReorged",
					blockNumber, reorgedAt,
				)
			}
		} else {
			if reorged {
				t.Fatalf(
					"blockNumber %d is before reorg block at %d, but it was not tagged \"false\" in wasReorged",
					blockNumber, reorgedAt,
				)
			}
		}
	}
}
