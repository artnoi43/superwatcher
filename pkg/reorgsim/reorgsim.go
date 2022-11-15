package reorgsim

import (
	"github.com/artnoi43/superwatcher"
)

// reorgSim implements superwatcher.EthClient[block],
// and will be used in place of the normal Ethereum client
// to test the default emitter implementation.
type reorgSim struct {
	lastRecord   uint64
	lookBack     uint64
	chain        blockChain
	reorgedChain blockChain

	seen map[uint64]int
}

// NewReorgSim returns a new reorgSim with hard-coded good and reorged chains.
func NewReorgSim(lookBack, lastRecord, reorgedAt uint64, logFiles []string) superwatcher.EthClient {
	mappedLogs := InitLogs(logFiles)
	chain, reorgedChain := NewBlockChain(mappedLogs, reorgedAt)

	return &reorgSim{
		lastRecord:   lastRecord,
		lookBack:     lookBack,
		chain:        chain,
		reorgedChain: reorgedChain,
		seen:         make(map[uint64]int),
	}
}
