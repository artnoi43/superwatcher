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
	tests := map[uint64]int{
		69:  1,
		70:  0,
		1:   69,
		0:   70,
		150: 7,
		151: 6,
		7:   151,
		6:   150,
		5:   70,
	}

	seen := make(map[common.Hash]bool)
	for blockNumber, reorgIndex := range tests {
		h := ReorgHash(blockNumber, reorgIndex)
		if seen[h] {
			t.Fatal("got hash collisions")
		}

		seen[h] = true

		h1 := ReorgHash(blockNumber, reorgIndex)
		if h != h1 {
			t.Fatalf(
				"ReorgHash(%d, %d) is not deterministic: expecting %s, got %s",
				blockNumber, reorgIndex, h.String(), h1.String(),
			)
		}
	}
}
