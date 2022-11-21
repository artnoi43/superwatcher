package reorgsim

import (
	"github.com/artnoi43/superwatcher"
)

type ReorgParam struct {
	StartBlock    uint64
	currentBlock  uint64 // currentBlock is hidden from outside for using exclusively in BlockByNumber
	BlockProgress uint64
	ReorgedAt     uint64
	ExitBlock     uint64 // reorgSim.HeaderByNumber will return ErrExitBlockReached once its currentBlock reach this
}

// reorgSim implements superwatcher.EthClient[block],
// and will be used in place of the normal Ethereum client
// to test the default emitter implementation.
type reorgSim struct {
	ReorgParam

	chain        blockChain
	reorgedChain blockChain

	seen map[uint64]int
}

// NewReorgSim returns a new reorgSim with hard-coded good and reorged chains.
func NewReorgSim(param ReorgParam, logFiles []string) superwatcher.EthClient {
	mappedLogs := InitLogs(logFiles)
	chain, reorgedChain := NewBlockChain(mappedLogs, param.ReorgedAt)

	return &reorgSim{
		ReorgParam:   param,
		chain:        chain,
		reorgedChain: reorgedChain,
		seen:         make(map[uint64]int),
	}
}
