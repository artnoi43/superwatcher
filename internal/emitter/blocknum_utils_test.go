package emitter

import "testing"

func TestFromBlockToBlockNormal(t *testing.T) {
	type blocknums struct {
		start       uint64
		current     uint64
		lastRec     uint64
		filterRange uint64
	}
	type answers struct {
		fromBlock uint64
		toBlock   uint64
	}

	tests := map[blocknums]answers{
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
