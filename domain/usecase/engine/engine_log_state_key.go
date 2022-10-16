package engine

import (
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
)

type engineLogStateKey struct {
	address     string
	blockNumber uint64
}

func (k engineLogStateKey) BlockNumber() uint64 {
	// TODO: Here for debugging
	if k.blockNumber == 0 {
		panic("got blockNumber 0 from a serviceLogStateKey")
	}
	return k.blockNumber
}

func engineLogStateKeyFromLog(l *types.Log) engineLogStateKey {
	return engineLogStateKey{
		address:     gslutils.ToLower(l.Address.String()),
		blockNumber: l.BlockNumber,
	}
}
