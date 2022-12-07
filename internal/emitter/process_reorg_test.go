package emitter

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

func TestProcessReorg(t *testing.T) {
	for i, tc := range testCases {
		b, _ := json.Marshal(tc)
		t.Logf("testCase: %s", b)
		err := testProcessReorg(tc)
		if err != nil {
			t.Fatalf("Case %d: %s", i, err.Error())
		}
	}
}

func testProcessReorg(tc testConfig) error {
	tracker := newTracker("testProcessReorg", 3)
	logs := reorgsim.InitMappedLogsFromFiles(tc.LogsFiles)
	oldChain, reorgedChain := reorgsim.NewBlockChainV2(logs, tc.ReorgedAt, tc.MovedLogs)

	// Add oldChain's blocks to tracker
	for blockNumber, block := range oldChain {
		blockLogs := block.Logs()
		logs := gslutils.CollectPointers(blockLogs)

		tracker.addTrackerBlockInfo(&superwatcher.BlockInfo{
			Logs:   logs,
			Hash:   block.Hash(),
			Number: blockNumber,
		})
	}

	// Collect reorgedLogs for checking
	var reorgedLogs []types.Log
	for _, block := range reorgedChain {
		if logs := block.Logs(); len(logs) != 0 {
			reorgedLogs = append(reorgedLogs, logs...)
		}
	}

	// Collect movedFrom blockNumbers
	var movedLogFromBlockNumbers []uint64
	var movedLogToBlockNumbers []uint64
	for blockNumber, moves := range tc.MovedLogs {
		movedLogFromBlockNumbers = append(movedLogFromBlockNumbers, blockNumber)
		for _, move := range moves {
			movedLogToBlockNumbers = append(movedLogToBlockNumbers, move.NewBlock)
		}
	}

	// Call mapFreshLogs with reorgedLogs
	freshHashes, freshLogs, processLogs := mapFreshLogs(reorgedLogs)

	wasReorged, err := processReorg(
		tracker,
		tc.FromBlock,
		tc.ToBlock,
		freshHashes,
		freshLogs,
		processLogs,
	)

	if err != nil {
		return err
	}

	for blockNumber := tc.FromBlock; blockNumber <= tc.ToBlock; blockNumber++ {
		// Skip blocks without logs
		if len(logs[blockNumber]) == 0 {
			continue
		}

		reorged := wasReorged[blockNumber]

		// Any blocks after c.reorgedAt should be reorged.
		if blockNumber >= tc.ReorgedAt {
			if reorged {
				continue
			}

			return fmt.Errorf(
				"blockNumber %d is after reorg block at %d, but it was not tagged \"true\" in wasReorged: %v",
				blockNumber, tc.ReorgedAt, wasReorged,
			)
		}

		// And any block before c.reorgedAt should NOT be reorged.
		if reorged {
			return fmt.Errorf("blockNumber %d is before reorg block at %d, but it was not tagged \"false\" in wasReorged: %v",
				blockNumber, tc.ReorgedAt, wasReorged,
			)
		}
	}

	// Blocks from which logs were moved must be tagged as reorged
	for _, blockNumber := range movedLogFromBlockNumbers {
		reorged := wasReorged[blockNumber]
		if !reorged {
			return fmt.Errorf("movedLogFromBlock %d was not tagged as reorged", blockNumber)
		}
	}

	// Blocks to which logs were moved must be tagged as reorged
	for _, blockNumber := range movedLogToBlockNumbers {
		reorged := wasReorged[blockNumber]
		if !reorged {
			return fmt.Errorf("movedLogTo %d was not tagged as reorged", blockNumber)
		}
	}

	return nil
}
