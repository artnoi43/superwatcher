package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/artnoi43/gsl"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/emitter/emittertest"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

func TestDeleteUnusable(t *testing.T) {
	mapResults := make(map[uint64]*mapLogsResult)
	var lastGood uint64 = 15
	to := lastGood + 10
	for i := uint64(0); i < to; i++ {
		mapResults[i] = &mapLogsResult{
			Block: superwatcher.Block{},
		}
	}

	deleteUnusableResult(mapResults, lastGood)
	for i := lastGood; i <= to; i++ {
		v, ok := mapResults[i]
		if !ok {
			continue
		}

		t.Errorf("deleted %d but value still in map %v", i, v)
	}
}

func TestMapLogs(t *testing.T) {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		superwatcher.PolicyExpensive,
	} {
		for i, tc := range emittertest.TestCasesV1 {
			b, _ := json.Marshal(tc)
			t.Logf("testCase: %s (policy %d)", b, policy)
			err := testMapLogsV1(&tc, policy)
			if err != nil {
				t.Fatalf("Case %d (policy %d): %s", i, policy, err.Error())
			}
		}
	}
}

// testMapLogsV1 tests function mapLogs with ReorgSimV1 (1 reorg)
func testMapLogsV1(tc *emittertest.TestConfig, policy superwatcher.Policy) error {
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
		logs := gsl.CollectPointers(blockLogs)
		concatLogs[blockNumber] = append(concatLogs[blockNumber], logs...)

		tracker.addTrackerBlock(&superwatcher.Block{
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
			concatLogs[blockNumber] = append(concatLogs[blockNumber], gsl.CollectPointers(logs)...)
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
	mapResults, err := mapLogs(
		nil,
		tc.FromBlock,
		tc.ToBlock,
		gsl.CollectPointers(reorgedLogs),
		true,
		tracker,
		mockClient,
		policy,
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

type mockGetHeader struct {
	chain reorgsim.BlockChain
}

func (h *mockGetHeader) HeaderByNumber(ctx context.Context, number *big.Int) (superwatcher.BlockHeader, error) {
	return h.chain[number.Uint64()], nil
}
