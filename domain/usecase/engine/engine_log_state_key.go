package engine

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
)

// engineLogState key is used to identify a log across multiple
// filterLogs loop. This means that it must be specific enough to capture an event,
// while at the same time
type engineLogStateKey struct {
	address     string
	txHash      string
	topic0      string
	blockNumber uint64
}

func (k engineLogStateKey) BlockNumber() uint64 {
	// TODO: Here for debugging
	if k.blockNumber == 0 {
		panic("got blockNumber 0 from a serviceLogStateKey")
	}
	return k.blockNumber
}

func (k engineLogStateKey) String() string {
	return fmt.Sprintf("%s:%s:%s:%d", k.address, k.txHash, k.topic0, k.blockNumber)
}

func engineLogStateKeyFromLog(l *types.Log) engineLogStateKey {
	return engineLogStateKey{
		address:     gslutils.ToLower(l.Address.String()),
		txHash:      l.TxHash.String(),
		topic0:      l.Topics[0].String(),
		blockNumber: l.BlockNumber,
	}
}
