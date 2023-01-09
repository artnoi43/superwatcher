package poller

import (
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/emitter/emittertest"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

func TestMapLogsNg(t *testing.T) {
	tc := emittertest.TestCasesV1[0]
	if err := testMapLogsNg(tc); err != nil {
		t.Error(err.Error())
	}
}

func testMapLogsNg(tc emittertest.TestConfig) error {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		// superwatcher.PolicyExpensive,
	} {

		tracker := newTracker("testProcessReorg", 3)
		logs := reorgsim.InitMappedLogsFromFiles(tc.LogsFiles...)

		var reorgEvent *reorgsim.ReorgEvent
		if len(tc.Events) == 0 {
			reorgEvent = new(reorgsim.ReorgEvent)
		} else {
			reorgEvent = &tc.Events[0]
		}

		oldChain, reorgedChain := reorgsim.NewBlockChainReorgMoveLogs(logs, *reorgEvent)
		mockClient, err := reorgsim.NewReorgSim(tc.Param, []reorgsim.ReorgEvent{tc.Events[0]}, reorgedChain, nil, "", 4)
		if err != nil {
			return errors.Wrap(err, "cannot init ReorgSim for testMapLogsV1")
		}

		// concatLogs store all logs, so that we can **skip block with out any logs**, fresh or reorged
		concatLogs := make(map[uint64][]*types.Log)

		// Add oldChain's blocks to tracker
		for blockNumber, block := range oldChain {
			blockLogs := block.Logs()
			logs := gslutils.CollectPointers(blockLogs)
			concatLogs[blockNumber] = append(concatLogs[blockNumber], logs...)

			tracker.addTrackerBlock(&superwatcher.Block{
				Logs:   logs,
				Hash:   block.Hash(),
				Number: blockNumber,
			})
		}

		pollResults := make(map[uint64]superwatcher.Block)
		// Collect reorgedLogs for checking
		for blockNumber, block := range reorgedChain {
			if logs := block.Logs(); len(logs) != 0 {
				pollResults[blockNumber] = superwatcher.Block{
					Number: blockNumber,
					Hash:   logs[0].BlockHash,
					Logs:   gslutils.CollectPointers(logs),
				}

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
		mapResults, err := mapLogsNg(
			nil,
			tc.FromBlock,
			tc.ToBlock,
			policy,
			pollResults,
			mockClient,
			tracker,
			true,
		)
		if err != nil {
			return errors.Wrap(err, "error in mapLogs")
		}

		wasReorged := make(map[uint64]bool)
		for k, v := range mapResults {
			wasReorged[k] = v.reorged
		}

		for blockNumber := tc.FromBlock; blockNumber <= tc.ToBlock; blockNumber++ {
			// Skip blocks without logs
			if len(concatLogs[blockNumber]) == 0 {
				continue
			}

			mapResult, ok := mapResults[blockNumber]
			if !ok {
				mapResult = new(mapLogsResult)
			}

			reorged := mapResult.reorged
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
				return fmt.Errorf(
					"blockNumber %d is before reorg block at %d, but it was not tagged \"false\" in wasReorged: %v",
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
	}

	return nil
}
