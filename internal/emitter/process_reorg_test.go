package emitter

import (
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

type testConfig struct {
	StartBlock uint64   `json:"startBlock"`
	ReorgedAt  uint64   `json:"reorgedAt"`
	FromBlock  uint64   `json:"fromBlock"`
	ToBlock    uint64   `json:"toBlock"`
	LogsFiles  []string `json:"logs"`
}

var testCases = []testConfig{
	{
		StartBlock: 15944390,
		ReorgedAt:  15944411,
		FromBlock:  15944400,
		ToBlock:    15944500,
		LogsFiles: []string{
			"./assets/logs_poolfactory.json",
			"./assets/logs_lp.json",
		},
	},
	{
		StartBlock: 15965710,
		ReorgedAt:  15965730,
		FromBlock:  15965717,
		ToBlock:    15965748,
		LogsFiles: []string{
			"./assets/logs_lp_2_1.json",
			"./assets/logs_lp_2_2.json",
		},
	},
	{
		StartBlock: 15965800,
		ReorgedAt:  15965811,
		FromBlock:  15965802,
		ToBlock:    15965835,
		LogsFiles: []string{
			"./assets/logs_lp_3_1.json",
			"./assets/logs_lp_3_2.json",
		},
	},
	{
		StartBlock: 15966455,
		ReorgedAt:  15966475,
		FromBlock:  15966460,
		ToBlock:    15966479,
		LogsFiles: []string{
			"./assets/logs_lp_4.json",
		},
	},
	{
		StartBlock: 15966490,
		ReorgedAt:  15966536,
		FromBlock:  15966500,
		ToBlock:    15966536,
		LogsFiles: []string{
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

func testProcessReorg(c testConfig) error {
	tracker := newTracker()
	hardcodedLogs := reorgsim.InitLogs(c.LogsFiles)
	oldChain, reorgedChain := reorgsim.NewBlockChain(hardcodedLogs, c.ReorgedAt)

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
		c.FromBlock,
		c.ToBlock,
		freshHashes,
		freshLogs,
		processLogs,
	)

	for blockNumber, reorged := range wasReorged {
		// Any blocks after c.reorgedAt should be reorged.
		if blockNumber >= c.ReorgedAt {
			if reorged {
				continue
			}

			return fmt.Errorf(
				"blockNumber %d is after reorg block at %d, but it was not tagged \"true\" in wasReorged",
				blockNumber, c.ReorgedAt,
			)
		}

		// And any block before c.reorgedAt should NOT be reorged.
		if reorged {
			return fmt.Errorf("blockNumber %d is before reorg block at %d, but it was not tagged \"false\" in wasReorged",
				blockNumber, c.ReorgedAt,
			)
		}
	}

	return nil
}
