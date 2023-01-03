package reorgsim

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
)

// PRandomHash returns a deterministic, pseudo-random hash for i
func PRandomHash(i uint64) common.Hash {
	return common.BigToHash(big.NewInt(int64(i)))
}

func ReorgHash(blockNumber uint64, reorgIndex int) common.Hash {
	reorgSeed := rand.NewSource(int64(reorgIndex)).Int63()
	return PRandomHash(blockNumber + uint64(reorgSeed))
}
