package reorgsim

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestPRandomHash(t *testing.T) {
	// We should get the same hash for the same input number
	hash69 := PRandomHash(69)
	_hash69 := PRandomHash(69)

	if hash69 != _hash69 {
		t.Logf("%s vs %s\n", hash69, _hash69)
		t.Log("hashes not matches\n")
		t.Fatal()
	}
}

func TestReorgHash(t *testing.T) {
	// tests map blockNumber and reorgIndex for func ReorgHash.
	// No combinations of these should cause collisions.
	type param struct {
		blockNumber uint64
		reorgIndex  int
	}

	tests := []param{
		{blockNumber: 69, reorgIndex: 1},
		{blockNumber: 69, reorgIndex: -1},
		{blockNumber: 70, reorgIndex: 0},
		{blockNumber: 1, reorgIndex: 69},
		{blockNumber: 0, reorgIndex: 70},
		{blockNumber: 150, reorgIndex: 7},
		{blockNumber: 151, reorgIndex: 6},
		{blockNumber: 7, reorgIndex: 151},
		{blockNumber: 6, reorgIndex: 150},
		{blockNumber: 6, reorgIndex: 151},
		{blockNumber: 500, reorgIndex: 151},
	}

	seen := make(map[common.Hash]bool)
	for _, test := range tests {
		// Check for collision
		h := ReorgHash(test.blockNumber, test.reorgIndex)
		if seen[h] {
			t.Fatal("got collisions from ReorgHash")
		}

		seen[h] = true

		// Check for non-deterministic behaviour
		h1 := ReorgHash(test.blockNumber, test.reorgIndex)
		if h != h1 {
			t.Fatalf(
				"ReorgHash(%d, %d) is not deterministic: expecting %s, got %s",
				test.blockNumber, test.reorgIndex, h.String(), h1.String(),
			)
		}
	}
}
