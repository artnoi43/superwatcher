package reorgsim

import (
	"sync"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

type Param struct {
	StartBlock    uint64
	BlockProgress uint64
	ReorgedBlock  uint64
	ExitBlock     uint64 // reorgSim.HeaderByNumber will return ErrExitBlockReached once its currentBlock reach this
	Debug         bool
}

// ReorgSim implements superwatcher.EthClient[block],
// and will be used in place of the normal Ethereum client
// to test the default emitter implementation.
type ReorgSim struct {
	sync.RWMutex
	Param

	chain             blockChain
	reorgedChain      blockChain
	currentBlock      uint64 // currentBlock is hidden from outside for using exclusively in BlockByNumber
	wasForked         bool
	filterLogsCounter map[uint64]int

	debugger debugger.Debugger
}

// NewReorgSimFromLogsFiles returns a new reorgSim with good and reorged chains mocked using the files.
func NewReorgSimFromLogsFiles(param Param, logFiles []string, logLevel uint8) superwatcher.EthClient {
	logs := InitLogsFromFiles(logFiles)
	chain, reorgedChain := NewBlockChainNg(logs, param.ReorgedBlock)

	return NewReorgSim(param, chain, reorgedChain, logLevel)
}

func NewReorgSim(param Param, chain, reorgedChain blockChain, logLevel uint8) superwatcher.EthClient {
	return &ReorgSim{
		Param:             param,
		chain:             chain,
		reorgedChain:      reorgedChain,
		filterLogsCounter: make(map[uint64]int),
		debugger:          *debugger.NewDebugger("reorgSim", logLevel),
	}
}

func (s *ReorgSim) Chain() blockChain {
	return s.chain
}

func (s *ReorgSim) ReorgedChain() blockChain {
	return s.reorgedChain
}
