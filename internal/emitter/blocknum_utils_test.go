package emitter

import "testing"

func TestFromBlockToBlock(t *testing.T) {
	type blocknums struct {
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
			current:     1000,
			lastRec:     100,
			filterRange: 300,
		}: {
			fromBlock: 0,
			toBlock:   401,
		},
		{
			current:     6000,
			lastRec:     500,
			filterRange: 300,
		}: {
			fromBlock: 201,
			toBlock:   801,
		},
	}
	for input, expected := range tests {
		from, to := fromBlockToBlock(input.current, input.lastRec, input.filterRange)
		var failed bool
		if expected.fromBlock != from {
			t.Fatal("fromBlock not matched")
			failed = true
		}
		if expected.toBlock != to {
			t.Fatal("toBlock not matched")
			failed = true
		}

		if failed {
			t.Logf("from: %v, to: %v, expectedFrom: %v,  expectedTo: %v\n", from, to, expected.fromBlock, expected.toBlock)
		}
	}
}
