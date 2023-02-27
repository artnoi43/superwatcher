package poller

import (
	"fmt"
	"testing"

	"github.com/artnoi43/gsl"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/testutils"
	"github.com/artnoi43/superwatcher/testlogs"
)

func init() {
	testlogs.SetLogsPath("../../testlogs")
}

func TestPollNg(t *testing.T) {
	err := testutils.RunTestCase(t, "TestPollNg", testlogs.TestCasesV1, testPollNg)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testPollNg(t *testing.T, caseNumber int) error {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		superwatcher.PolicyExpensive,
	} {
		tc := testlogs.TestCasesV1[caseNumber-1]
		tracker := newTracker("testPollNg", 3)
		logs := reorgsim.InitMappedLogsFromFiles(tc.LogsFiles...)

		var reorgEvent *reorgsim.ReorgEvent
		if len(tc.Events) == 0 {
			reorgEvent = new(reorgsim.ReorgEvent)
		} else {
			reorgEvent = &tc.Events[0]
		}

		oldChain, reorgedChains := reorgsim.NewBlockChain(logs, []reorgsim.ReorgEvent{*reorgEvent})
		reorgedChain := reorgedChains[0]

		mockClient, err := reorgsim.NewReorgSim(tc.Param, []reorgsim.ReorgEvent{tc.Events[0]}, reorgedChain, nil, "", 4)
		if err != nil {
			return errors.Wrap(err, "cannot init ReorgSim for testMapLogsV1")
		}

		// allLogs store all logs, so that we can **skip block with out any logs**, fresh or reorged
		allLogs := make(map[uint64][]*types.Log)

		// Add oldChain's blocks to tracker
		for blockNumber, block := range oldChain {
			blockLogs := block.Logs()
			logs := gsl.CollectPointers(blockLogs)
			allLogs[blockNumber] = append(allLogs[blockNumber], logs...)

			tracker.addTrackerBlock(&superwatcher.Block{
				Number: blockNumber,
				Header: block,
				Hash:   block.Hash(),
				Logs:   logs,
			})
		}
		// Add reorgChain's logs to allLogs
		for blockNumber, block := range reorgedChain {
			if logs := block.Logs(); len(logs) != 0 {
				allLogs[blockNumber] = append(allLogs[blockNumber], gsl.CollectPointers(logs)...)
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

		debugger := debugger.NewDebugger("poller-ng", 3)
		param := &param{
			fromBlock: tc.FromBlock,
			toBlock:   tc.ToBlock,
			policy:    policy,
		}

		blocksMissing := []uint64{}
		pollResults, err := poll(nil, param, nil, nil, mockClient, debugger)
		if err != nil {
			t.Fatal("pollCheap error", err.Error())
		}
		pollResults, blocksMissing, err = pollMissing(nil, param, mockClient, tracker, pollResults, debugger)
		if err != nil {
			t.Fatal("pollMissing error", err.Error())
		}
		pollResults, err = findReorg(param, blocksMissing, tracker, pollResults, debugger)
		if err != nil {
			t.Fatal("findReorg error")
		}

		wasReorged := make(map[uint64]bool)
		for k := range pollResults {
			v, ok := pollResults[k]
			if !ok {
				continue
			}
			wasReorged[k] = v.forked
		}

		for blockNumber := tc.FromBlock; blockNumber <= tc.ToBlock; blockNumber++ {
			// Skip blocks without logs
			if len(allLogs[blockNumber]) == 0 {
				continue
			}

			mapResult, ok := pollResults[blockNumber]
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
