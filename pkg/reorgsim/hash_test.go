package reorgsim

import "testing"

func TestDeterministicRandomHash(t *testing.T) {
	// We should get the same hash for the same input number

	hash69 := deterministicRandomHash(69)
	_hash69 := deterministicRandomHash(69)

	if hash69 != _hash69 {
		t.Logf("%s vs %s\n", hash69, _hash69)
		t.Log("hashes not matches\n")
		t.Fatal()
	}
}
