package reorgsim

import (
	"sync"

	"github.com/artnoi43/superwatcher"
)

type ReorgParam struct {
	StartBlock    uint64
	currentBlock  uint64 // currentBlock is hidden from outside for using exclusively in BlockByNumber
	BlockProgress uint64
	ReorgedAt     uint64
	ExitBlock     uint64 // reorgSim.HeaderByNumber will return ErrExitBlockReached once its currentBlock reach this
}

// ReorgSim implements superwatcher.EthClient[block],
// and will be used in place of the normal Ethereum client
// to test the default emitter implementation.
type ReorgSim struct {
	sync.RWMutex
	ReorgParam

	chain        blockChain
	reorgedChain blockChain
	forked       bool

	seenFilterLogs map[uint64]int
}

// NewReorgSim returns a new reorgSim with hard-coded good and reorged chains.
func NewReorgSim(param ReorgParam, logFiles []string) superwatcher.EthClient {
	mappedLogs := InitLogs(logFiles)
	chain, reorgedChain := NewBlockChain(mappedLogs, param.ReorgedAt)

	return &ReorgSim{
		ReorgParam:     param,
		chain:          chain,
		reorgedChain:   reorgedChain,
		seenFilterLogs: make(map[uint64]int),
	}
}

func (s *ReorgSim) Chain() blockChain {
	return s.chain
}

func (s *ReorgSim) ReorgedChain() blockChain {
	return s.reorgedChain
}
