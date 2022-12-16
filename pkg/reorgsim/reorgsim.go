package reorgsim

import (
	"sync"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// Param is embedded into ReorgSim, and represents the fake blockchain client parameters
type Param struct {
	// StartBlock will be used as initial ReorgSim.currentBlock. ReorgSim.currentBlock increases by `BlockProgess`
	// after each call to ReorgSim.BlockNumber.
	StartBlock uint64
	// BlockProgress represents how many block numbers should ReorgSim.currentBlock
	// should increase during each call to ReorgSim.BlockNumber.
	BlockProgress uint64
	// ReorgedBlock is the pivot block after which ReorgSim should use ReorgSim.reorgedChain.
	ReorgedBlock uint64
	// ExitBlock is used in reorgSim.BlockNumber to return ErrExitBlockReached once its currentBlock reaches ExitBlock.
	ExitBlock uint64
	// MovedLogs represents all of the moved logs after a chain reorg event. The map key is the source blockNumber.
	MovedLogs map[uint64][]MoveLogs

	Debug bool
}

// ReorgSim implements superwatcher.EthClient[block], and will be used in place of the normal Ethereum client
// to test the default emitter implementation.
// ReorgSim stores old chain and reorged chain, and it will not change the chain internal data. This means that
// it can support `MoveLogs`, even though the reorgedChain has to be reorged before being sent to NewReorgSim.
// To create a new ReorgSim with MoveLogs functionality, populate `Param.MovedLogs` with the desired values
// before calling NewReorgSimFromLogsFiles.
type ReorgSim struct {
	sync.RWMutex
	Param

	// chain is source for all blocks before Param.ReorgedBlock
	chain blockChain
	// reorgedChain is source for logs after Param.ReorgedBlock. The logic for doing this is in method ReorgSim.chooseBlock
	reorgedChain blockChain
	// currentBlock tracks the current block for the fake blockChain and is hidden from outside for using exclusively in BlockByNumber
	currentBlock uint64
	// filterLogsCounter is used to switch chain for a blockNumber, after certain number of calls to FilterLogs
	filterLogsCounter map[uint64]int
	// wasForked tracks if the fake chain was forked (i.e. if ReorgChain has returned any logs from a reorgedChain)
	wasForked bool

	debugger *debugger.Debugger
}

// NewReorgSim returns the mocked client.EthClient. After the internal state ReorgSim.currentBlock exceeds |reorgedAt|,
// ReorgSim methods returns data from |reorgedChain|.
func NewReorgSim(
	param Param,
	chain blockChain,
	reorgedChain blockChain,
	logLevel uint8,
) superwatcher.EthClient {
	return &ReorgSim{
		Param:             param,
		chain:             chain,
		reorgedChain:      reorgedChain,
		filterLogsCounter: make(map[uint64]int),
		debugger:          debugger.NewDebugger("reorgSim", logLevel),
	}
}

// NewReorgSimFromLogsFiles returns a new reorgSim with good and reorged chains mocked using the files.
func NewReorgSimFromLogsFiles(
	param Param,
	logFiles []string,
	logLevel uint8,
) superwatcher.EthClient {
	logs := InitMappedLogsFromFiles(logFiles...)
	chain, reorgedChain := NewBlockChainWithMovedLogs(logs, param.ReorgedBlock, param.MovedLogs)

	return NewReorgSim(param, chain, reorgedChain, logLevel)
}

func (s *ReorgSim) Chain() blockChain {
	return s.chain
}

func (s *ReorgSim) ReorgedChain() blockChain {
	return s.reorgedChain
}
