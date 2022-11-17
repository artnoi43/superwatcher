package emitter

import (
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

type TestConfig struct {
	fromBlock uint64
	toBlock   uint64
	reorgedAt uint64
	logs      []string
}

var testCases = []TestConfig{
	{
		fromBlock: 15944400,
		toBlock:   15944500,
		reorgedAt: 15944444,
		logs: []string{
			"./assets/logs_poolfactory.json",
			"./assets/logs_lp.json",
		},
	},
	{
		fromBlock: 15965717,
		toBlock:   15965748,
		reorgedAt: 15965730,
		logs: []string{
			"./assets/logs_lp_2_1.json",
			"./assets/logs_lp_2_2.json",
		},
	},
	{
		fromBlock: 15965802,
		toBlock:   15965835,
		reorgedAt: 15965803,
		logs: []string{
			"./assets/logs_lp_3_1.json",
			"./assets/logs_lp_3_2.json",
		},
	},
	{
		fromBlock: 15966460,
		toBlock:   15966479,
		reorgedAt: 15966475,
		logs: []string{
			"./assets/logs_lp_4.json",
		},
	},
	{
		fromBlock: 15966500,
		toBlock:   15966536,
		reorgedAt: 15966536,
		logs: []string{
			"./assets/logs_lp_5.json",
		},
	},
}

func TestProcessReorg(t *testing.T) {
	for i, tc := range testCases {
		err := testProcessReorg(tc)
		if err != nil {
			t.Fatalf("Case %d: %s", i, err.Error())
		}
	}
}

func testProcessReorg(c TestConfig) error {
	tracker := newTracker()
	hardcodedLogs := reorgsim.InitLogs(c.logs)
	oldChain, reorgedChain := reorgsim.NewBlockChain(hardcodedLogs, c.reorgedAt)

	// Add oldChain's blocks to tracker
	for blockNumber, block := range oldChain {
		blockLogs := block.Logs()
		logs := gslutils.CollectPointers(&blockLogs)

		tracker.addTrackerBlockInfo(&superwatcher.BlockInfo{
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

	freshHashes, freshLogs, processLogs := mapFreshLogsByHashes(reorgedLogs, reorgedHeader)

	wasReorged := processReorged(
		tracker,
		c.fromBlock,
		c.toBlock,
		freshHashes,
		freshLogs,
		processLogs,
	)

	for blockNumber, reorged := range wasReorged {
		// Any blocks after c.reorgedAt should be reorged.
		if blockNumber >= c.reorgedAt {
			if reorged {
				continue
			}

			return fmt.Errorf(
				"blockNumber %d is after reorg block at %d, but it was not tagged \"true\" in wasReorged",
				blockNumber, c.reorgedAt,
			)
		}

		// And any block before c.reorgedAt should NOT be reorged.
		if reorged {
			return fmt.Errorf("blockNumber %d is before reorg block at %d, but it was not tagged \"false\" in wasReorged",
				blockNumber, c.reorgedAt,
			)
		}
	}

	return nil
}
