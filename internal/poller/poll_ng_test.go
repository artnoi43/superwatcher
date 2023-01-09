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
	"github.com/artnoi43/superwatcher/pkg/testutils"
)

func TestMapLogsNg(t *testing.T) {
	err := testutils.RunTestCase(t, "TestMapLogsNg", emittertest.TestCasesV1, testMapLogsNg)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testMapLogsNg(t *testing.T, caseNumber int) error {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		// superwatcher.PolicyExpensive,
	} {
		tc := emittertest.TestCasesV1[caseNumber-1]
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

			b := &superwatcher.Block{
				Number: blockNumber,
				Hash:   block.Hash(),
				Logs:   logs,
			}

			if policy >= superwatcher.PolicyExpensive {
				b.Header = block
			}

			tracker.addTrackerBlock(b)
		}

		pollResults := make(map[uint64]*mapLogsResult)
		// Collect reorgedLogs for checking
		for blockNumber, block := range reorgedChain {
			if logs := block.Logs(); len(logs) != 0 {
				b := superwatcher.Block{
					Number: blockNumber,
					Hash:   logs[0].BlockHash,
					Logs:   gslutils.CollectPointers(logs),
				}

				if policy >= superwatcher.PolicyExpensive {
					b.Header = block
				}

				pollResults[blockNumber] = &mapLogsResult{Block: b}
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
		mapResults, err := findMissing(
			nil,
			tc.FromBlock,
			tc.ToBlock,
			policy,
			pollResults,
			mockClient,
			tracker,
		)
		if err != nil {
			return errors.Wrap(err, "error in mapLogs")
		}

		wasReorged := make(map[uint64]bool)
		for k, v := range mapResults {
			wasReorged[k] = v.forked
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

			reorged := mapResult.forked
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
