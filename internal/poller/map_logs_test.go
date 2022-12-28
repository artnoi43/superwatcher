package poller

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/emitter/emittertest"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

func TestMapLogs(t *testing.T) {
	for i, tc := range emittertest.TestCasesV1 {
		b, _ := json.Marshal(tc)
		t.Logf("testCase: %s", b)
		err := testMapLogsV1(&tc)
		if err != nil {
			t.Fatalf("Case %d: %s", i, err.Error())
		}
	}
}

// testMapLogsV1 tests function mapLogs with ReorgSimV1 (1 reorg)
func testMapLogsV1(tc *emittertest.TestConfig) error {
	tracker := newTracker("testProcessReorg", 3)
	logs := reorgsim.InitMappedLogsFromFiles(tc.LogsFiles...)

	var reorgEvent *reorgsim.ReorgEvent
	if len(tc.Events) == 0 {
		reorgEvent = new(reorgsim.ReorgEvent)
	} else {
		reorgEvent = &tc.Events[0]
	}

	oldChain, reorgedChain := reorgsim.NewBlockChainReorgMoveLogs(logs, *reorgEvent)

	// concatLogs store all logs, so that we can **skip block with out any logs**, fresh or reorged
	var concatLogs = make(map[uint64][]*types.Log)

	// Add oldChain's blocks to tracker
	for blockNumber, block := range oldChain {
		blockLogs := block.Logs()
		logs := gslutils.CollectPointers(blockLogs)
		concatLogs[blockNumber] = append(concatLogs[blockNumber], logs...)

		tracker.addTrackerBlockInfo(&superwatcher.BlockInfo{
			Logs:   logs,
			Hash:   block.Hash(),
			Number: blockNumber,
		})
	}

	// Collect reorgedLogs for checking
	var reorgedLogs []types.Log
	for blockNumber, block := range reorgedChain {
		if logs := block.Logs(); len(logs) != 0 {
			reorgedLogs = append(reorgedLogs, logs...)
			concatLogs[blockNumber] = append(concatLogs[blockNumber], gslutils.CollectPointers(logs)...)
		}
	}

	// Collect movedFrom blockNumbers
	var movedLogFromBlockNumbers []uint64
	var movedLogToBlockNumbers []uint64
	for blockNumber, moves := range tc.Events[0].MovedLogs {
		movedLogFromBlockNumbers = append(movedLogFromBlockNumbers, blockNumber)
		for _, move := range moves {
			movedLogToBlockNumbers = append(movedLogToBlockNumbers, move.NewBlock)
		}
	}

	// Call mapFreshLogs with reorgedLogs
	wasReorged, _, _, err := mapLogs(
		tc.FromBlock,
		tc.ToBlock,
		gslutils.CollectPointers(reorgedLogs),
		tracker,
		nil,
	)

	if err != nil {
		return errors.Wrap(err, "error in mapLogs")
	}

	for blockNumber := tc.FromBlock; blockNumber <= tc.ToBlock; blockNumber++ {
		// Skip blocks without logs
		if len(concatLogs[blockNumber]) == 0 {
			continue
		}

		reorged := wasReorged[blockNumber]

		// Any blocks after c.reorgedAt should be reorged.
		if blockNumber >= reorgEvent.ReorgBlock {
			if reorged {
				continue
			}

			return fmt.Errorf(
				"blockNumber %d is after reorg block at %d, but it was not tagged \"true\" in wasReorged: %v",
				blockNumber, reorgEvent.ReorgBlock, wasReorged,
			)
		}

		// And any block before c.reorgedAt should NOT be reorged.
		if reorged {
			return fmt.Errorf("blockNumber %d is before reorg block at %d, but it was not tagged \"false\" in wasReorged: %v",
				blockNumber, reorgEvent.ReorgBlock, wasReorged,
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
