package emitter

import (
	"testing"
)

func TestComputeFromBlockToBlock(t *testing.T) {
	type testCase struct {
		Name               string           `json:"Name"`
		CurrentBlock       uint64           `json:"currentBlock"`
		LastRecordedBlocks uint64           `json:"lastRecordedBlocks"`
		LookBackBlocks     uint64           `json:"lookBackBlocks"`
		LookBackRetries    uint64           `json:"lookBackRetries"`
		LookBackFirstStart bool             `json:"lookbackFirstStart"`
		Status             *filterLogStatus `json:"status"`
	}

	normalStatus := &filterLogStatus{IsReorging: false}
	reorgingStatus := &filterLogStatus{IsReorging: true}

	// Maps testCase to fromBlock, toBlock
	tests := map[testCase]struct{ fromBlock, toBlock uint64 }{
		{
			Name:               "normal_10",
			CurrentBlock:       200,
			LastRecordedBlocks: 10, // Base
			LookBackBlocks:     10, // LookBack
			LookBackRetries:    2,
			LookBackFirstStart: false,
			Status:             normalStatus,
		}: {
			fromBlock: 1,  // (Base+1) + LookBack
			toBlock:   21, // fromBlock + LookBack
		},
		{
			Name:               "normal_100",
			CurrentBlock:       200,
			LastRecordedBlocks: 100, // Base
			LookBackBlocks:     10,  // LookBack
			LookBackRetries:    2,
			LookBackFirstStart: false,
			Status:             normalStatus,
		}: {
			fromBlock: 91,  // (Base+1) - LookBack
			toBlock:   111, // fromBlock + LookBack
		},
		{
			Name:               "reorging_10",
			CurrentBlock:       200,
			LastRecordedBlocks: 10, // Base
			LookBackBlocks:     10, // LookBack
			LookBackRetries:    2,
			LookBackFirstStart: false,
			Status:             reorgingStatus, // GoBack
		}: {
			fromBlock: 1,  // (Base+1) - GoBack
			toBlock:   11, // fromBlock + LookBack
		},
		{
			Name:               "reorging_100",
			CurrentBlock:       200,
			LastRecordedBlocks: 100, // Base
			LookBackBlocks:     10,
			LookBackRetries:    2,
			LookBackFirstStart: false,
			Status:             reorgingStatus, // GoBack
		}: {
			fromBlock: 81, // (Base+1) - GoBack
			toBlock:   91, // fromBlock + LookBack
		},
		{
			Name:               "normal_first_100",
			CurrentBlock:       200,
			LastRecordedBlocks: 100, // Base
			LookBackBlocks:     10,
			LookBackRetries:    2,
			LookBackFirstStart: true, // GoBack
			Status:             normalStatus,
		}: {
			fromBlock: 81, // (Base+1) - GoBack
			toBlock:   91, // fromBlock + LookBack
		},
		{
			Name:               "normal_first_50",
			CurrentBlock:       200,
			LastRecordedBlocks: 50, // Base
			LookBackBlocks:     10,
			LookBackRetries:    2,
			LookBackFirstStart: true, // GoBack
			Status:             normalStatus,
		}: {
			fromBlock: 31, // (Base+1) - GoBack
			toBlock:   41, // fromBlock + LookBack
		},
		{
			Name:               "reorg_first_100",
			CurrentBlock:       200,
			LastRecordedBlocks: 100, // Base
			LookBackBlocks:     10,
			LookBackRetries:    2,
			LookBackFirstStart: true, // GoBack
			Status:             normalStatus,
		}: {
			fromBlock: 81, // (Base+1) - GoBack
			toBlock:   91, // fromBlock + LookBack
		},
		{
			Name:               "reorg_first_50",
			CurrentBlock:       200,
			LastRecordedBlocks: 50, // Base
			LookBackBlocks:     10,
			LookBackRetries:    2,
			LookBackFirstStart: true, // GoBack
			Status:             normalStatus,
		}: {
			fromBlock: 31, // (Base+1) - GoBack
			toBlock:   41, // fromBlock + LookBack
		},
		{
			Name:               "normal_6999",
			CurrentBlock:       7000,
			LastRecordedBlocks: 6999, // Base
			LookBackBlocks:     10,   // LookBack
			LookBackRetries:    2,
			LookBackFirstStart: false,
			Status:             normalStatus,
		}: {
			fromBlock: 6990, // (Base+1) - LookBack
			toBlock:   7000, // fromBlock + LookBack
		},
		{
			Name:               "normal_first_6999",
			CurrentBlock:       7000,
			LastRecordedBlocks: 6999, // Base
			LookBackBlocks:     10,
			LookBackRetries:    2,
			LookBackFirstStart: true, // GoBack
			Status:             normalStatus,
		}: {
			fromBlock: 6980, // (Base+1) - GoBack
			toBlock:   6990, // fromBlock + LookBack
		},

		{
			Name:               "normal_first_14999",
			CurrentBlock:       15000,
			LastRecordedBlocks: 14999, // Base
			LookBackBlocks:     2000,
			LookBackRetries:    2,
			LookBackFirstStart: true, // GoBack
			Status:             normalStatus,
		}: {
			fromBlock: 11000, // (Base+1) - GoBack
			toBlock:   13000, // fromBlock + LookBack
		},
	}

	for test, expected := range tests {
		fromBlock, toBlock := computeFromBlockToBlock(
			test.CurrentBlock,
			test.LastRecordedBlocks,
			test.LookBackBlocks,
			test.LookBackRetries,
			&test.LookBackFirstStart,
			test.Status,
			false,
		)

		var failed bool
		if fromBlock != expected.fromBlock {
			t.Logf("[%s] unexpected fromBlock: expecting %d, got %d", test.Name, expected.fromBlock, fromBlock)
			failed = true
		}
		if toBlock != expected.toBlock {
			t.Logf("[%s] unexpected toBlock: expecting %d, got %d", test.Name, expected.toBlock, toBlock)
			failed = true
		}

		if failed {
			t.Fatalf("unexpected result")
		}
	}
}
