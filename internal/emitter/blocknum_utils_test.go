package emitter

import (
	"testing"
)

func TestFromBlockToBlockNormal(t *testing.T) {
	type param struct {
		start       uint64
		current     uint64
		lastRec     uint64
		filterRange uint64
	}
	type answer struct {
		fromBlock uint64
		toBlock   uint64
	}

	tests := map[param]answer{
		{
			start:       69,
			current:     1000,
			lastRec:     100,
			filterRange: 300,
		}: {
			fromBlock: 69,
			toBlock:   400,
		},
		{
			current:     6000,
			lastRec:     500,
			filterRange: 300,
		}: {
			fromBlock: 201,
			toBlock:   800,
		},
	}

	for input, expected := range tests {
		from, to := fromBlockToBlockNormal(input.start, input.current, input.lastRec, input.filterRange)

		var failed bool
		if expected.fromBlock != from {
			t.Fatalf("fromBlock not matched: expecting %d, got %d", expected.fromBlock, from)
			failed = true
		}
		if expected.toBlock != to {
			t.Fatalf("toBlock not matched: expecting %d, got %d", expected.toBlock, to)
			failed = true
		}

		if failed {
			t.Logf("from: %v, to: %v, expectedFrom: %v,  expectedTo: %v\n", from, to, expected.fromBlock, expected.toBlock)
		}
	}
}

func TestFromBlockToBlockIsReorging(t *testing.T) {
	type param struct {
		prevFromBlock uint64
		prevToBlock   uint64
		start         uint64
		current       uint64
		lastRec       uint64
		filterRange   uint64
		retriesCount  uint64
		maxRetries    uint64
	}
	type answer struct {
		fromBlock uint64
		toBlock   uint64
	}

	newReorgStatus := func(
		fromBlock, toBlock, startBlock, currentBlock, lasRecordedBlock, retriesCount uint64,
	) *emitterStatus {
		return &emitterStatus{
			FromBlock:         fromBlock,
			ToBlock:           toBlock,
			CurrentBlock:      currentBlock,
			LastRecordedBlock: lasRecordedBlock,
			RetriesCount:      retriesCount,
			GoBackFirstStart:  false,
			IsReorging:        true,
		}
	}

	tests := map[param]answer{
		{
			prevFromBlock: 61,
			prevToBlock:   80,
			start:         30,
			current:       86,
			lastRec:       80,
			filterRange:   10,
			retriesCount:  1,
			maxRetries:    3,
		}: {
			fromBlock: 51,
			toBlock:   80,
		},
		{
			prevFromBlock: 51,
			prevToBlock:   80,
			start:         30,
			current:       86,
			lastRec:       80,
			filterRange:   10,
			retriesCount:  2,
			maxRetries:    3,
		}: {
			fromBlock: 41,
			toBlock:   80,
		},
		{
			prevFromBlock: 41,
			prevToBlock:   80,
			start:         30,
			current:       86,
			lastRec:       80,
			filterRange:   10,
			retriesCount:  3,
			maxRetries:    3,
		}: {
			fromBlock: 31,
			toBlock:   80,
		},
		{
			prevFromBlock: 31,
			prevToBlock:   80,
			start:         30,
			current:       86,
			lastRec:       80,
			filterRange:   10,
			retriesCount:  4,
			maxRetries:    3,
		}: {
			fromBlock: 31,
			toBlock:   80,
		},
		{
			prevFromBlock: 31,
			prevToBlock:   80,
			start:         30,
			current:       86,
			lastRec:       80,
			filterRange:   10,
			retriesCount:  5,
			maxRetries:    3,
		}: {
			fromBlock: 31,
			toBlock:   80,
		},
		{
			prevFromBlock: 281,
			prevToBlock:   300,
			start:         100,
			current:       305,
			lastRec:       300,
			filterRange:   20,
			retriesCount:  2,
			maxRetries:    2,
		}: {
			fromBlock: 261,
			toBlock:   300,
		},
		{
			prevFromBlock: 271,
			prevToBlock:   300,
			start:         100,
			current:       305,
			lastRec:       300,
			filterRange:   10,
			retriesCount:  3,
			maxRetries:    2,
		}: {
			fromBlock: 271,
			toBlock:   300,
		},
	}

	var c int
	for input, expected := range tests {
		c++

		prevStatus := newReorgStatus(
			input.prevFromBlock,
			input.prevToBlock,
			input.start,
			input.current,
			input.lastRec,
			input.retriesCount,
		)

		fromBlock, toBlock, err := fromBlockToBlockIsReorging(
			input.start,
			input.current,
			input.lastRec,
			input.filterRange,
			input.maxRetries,
			prevStatus,
		)

		t.Log(c, "err == nil", err == nil)

		if err != nil {
			if !(input.retriesCount >= input.maxRetries-1) {
				t.Errorf("[%d] maxRetries unexpectedly reached: %s", c, err.Error())
			}
		}

		if expected.fromBlock != fromBlock || expected.toBlock != toBlock {
			t.Errorf(
				"[%d] unexpected results: expecting (from %d, to %d), got (from %d, to %d)",
				c, expected.fromBlock, expected.toBlock, fromBlock, toBlock,
			)
		}
	}
}
