package emitter

import (
	"testing"

	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

func TestComputeFromBlockToBlock(t *testing.T) {
	type testCase struct {
		Name               string `json:"Name"`
		CurrentBlock       uint64 `json:"currentBlock"`
		LastRecordedBlocks uint64 `json:"lastRecordedBlocks"`
		FilterRange        uint64 `json:"filterRange"`
		MaxRetries         uint64 `json:"goBackRetries"`
		GoBackFirstStart   bool   `json:"goBackFirstStart"`
		Reorging           bool   `json:"reorging"`
	}

	newStatus := func(isReorging bool, fromBlock, toBlock uint64) *filterLogStatus {
		return &filterLogStatus{
			IsReorging:   isReorging,
			FromBlock:    fromBlock,
			ToBlock:      toBlock,
			RetriesCount: 1,
		}
	}

	reorging := true
	firstStart := true

	// Maps testCase to fromBlock, toBlock
	tests := map[testCase]struct{ fromBlock, toBlock uint64 }{
		{
			Name:               "normal_10",
			CurrentBlock:       200,
			LastRecordedBlocks: 10, // Base
			FilterRange:        10, // FilterRange
			MaxRetries:         2,
		}: {
			fromBlock: 1,  // (Base+1) + FilterRange
			toBlock:   20, // fromBlock + FilterRange
		},
		{
			Name:               "normal_100",
			CurrentBlock:       200,
			LastRecordedBlocks: 100, // Base
			FilterRange:        10,  // FilterRange
			MaxRetries:         2,
		}: {
			fromBlock: 91,  // (Base+1) - FilterRange
			toBlock:   110, // fromBlock + FilterRange
		},
		{
			Name:               "reorging_10",
			CurrentBlock:       200,
			LastRecordedBlocks: 10, // Base
			FilterRange:        10, // FilterRange
			MaxRetries:         2,
			Reorging:           reorging,
		}: {
			fromBlock: 0,  // (Base+1) - GoBack
			toBlock:   10, // fromBlock + FilterRange
		},
		{
			Name:               "reorging_100",
			CurrentBlock:       200,
			LastRecordedBlocks: 100, // Base
			FilterRange:        10,
			MaxRetries:         2,
			Reorging:           reorging,
		}: {
			fromBlock: 81,  // (Base+1) - GoBack
			toBlock:   100, // fromBlock + FilterRange
		},
		{
			Name:               "normal_first_100",
			CurrentBlock:       200,
			LastRecordedBlocks: 100, // Base
			FilterRange:        10,
			MaxRetries:         2,
			GoBackFirstStart:   firstStart, // GoBack
		}: {
			fromBlock: 81, // (Base+1) - GoBack
			toBlock:   90, // fromBlock + FilterRange
		},
		{
			Name:               "normal_first_50",
			CurrentBlock:       200,
			LastRecordedBlocks: 50, // Base
			FilterRange:        10,
			MaxRetries:         2,
			GoBackFirstStart:   firstStart, // GoBack
		}: {
			fromBlock: 31, // (Base+1) - GoBack
			toBlock:   40, // Base + FilterRange
		},
		{
			Name:               "reorg_first_100",
			CurrentBlock:       200,
			LastRecordedBlocks: 100, // Base
			FilterRange:        10,
			MaxRetries:         2,
			GoBackFirstStart:   firstStart, // GoBack
			Reorging:           reorging,
		}: {
			fromBlock: 81, // (Base+1) - GoBack
			toBlock:   90, // fromBlock + FilterRange
		},
		{
			Name:               "reorg_first_50",
			CurrentBlock:       200,
			LastRecordedBlocks: 50, // Base
			FilterRange:        10,
			MaxRetries:         2,
			GoBackFirstStart:   firstStart, // GoBack
			Reorging:           reorging,
		}: {
			fromBlock: 31, // (Base+1) - GoBack
			toBlock:   40, // fromBlock + FilterRange
		},
		{
			Name:               "normal_6999",
			CurrentBlock:       7000,
			LastRecordedBlocks: 6999, // Base
			FilterRange:        10,   // FilterRange
			MaxRetries:         2,
		}: {
			fromBlock: 6990, // (Base+1) - FilterRange
			toBlock:   7000, // fromBlock + FilterRange
		},
		{
			Name:               "normal_first_6999",
			CurrentBlock:       7000,
			LastRecordedBlocks: 6999, // Base
			FilterRange:        10,
			MaxRetries:         2,
			GoBackFirstStart:   firstStart,
		}: {
			fromBlock: 6980, // (Base+1) - GoBack
			toBlock:   6989, // fromBlock + FilterRange
		},

		{
			Name:               "normal_first_14999",
			CurrentBlock:       15000,
			LastRecordedBlocks: 14999, // Base
			FilterRange:        2000,
			MaxRetries:         2,
			GoBackFirstStart:   firstStart,
		}: {
			fromBlock: 11000, // (Base+1) - GoBack
			toBlock:   12999, // fromBlock + FilterRange
		},
	}

	for test, expected := range tests {
		prevFromBlock := 1 + test.LastRecordedBlocks - test.FilterRange

		fromBlock, toBlock, err := computeFromBlockToBlock(
			test.CurrentBlock,
			test.LastRecordedBlocks,
			test.FilterRange,
			test.MaxRetries,
			&test.GoBackFirstStart,
			newStatus(test.Reorging, prevFromBlock, test.LastRecordedBlocks),
			test.MaxRetries,
			debugger.NewDebugger("testComputeFromBlockToBlock", 4),
		)

		if err != nil {
			t.Error("error from computingFromBlockToBlock", err.Error())
		}

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
			t.Error("unexpected result")
		}
	}
}
