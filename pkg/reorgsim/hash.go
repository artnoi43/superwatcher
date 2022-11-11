package reorgsim

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
)

// randomHash returns random common.Hash.
// If i > 0, rand.Intn will be used, otherwise rand.Int.
// This allows us to generate duplicate p-random numbers.
func randomHash(i uint64) common.Hash {
	var b *big.Int
	if i < 0 {
		b = big.NewInt(rand.Int63())
	} else {
		b = big.NewInt(rand.Int63n(int64(i)))
	}
	return common.BigToHash(b)
}

func deterministicRandomHash(i uint64) common.Hash {
	return common.BigToHash(big.NewInt(int64(i)))
}
