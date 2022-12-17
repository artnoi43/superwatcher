package reorgsim

import (
	"sync"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// ReorgEvent is parameters for chain reorg events.
type ReorgEvent struct {
	// ReorgBlock is the pivot block after which ReorgSim should use ReorgSim.reorgedChain.
	ReorgBlock uint64 `json:"reorgBlock"`
	// MovedLogs represents all of the moved logs after a chain reorg event. The map key is the source blockNumber.
	MovedLogs map[uint64][]MoveLogs `json:"movedLogs"`
}

// BaseParam is the basic parameter for the mock client. Chain reorg parameters are NOT included here, but in ChainParam.
// It will be embedded in either ParamV1 or ParamV2, for ReorgSimV1 and ReorgSimV2 respectively.
type BaseParam struct {
	// StartBlock will be used as initial ReorgSim.currentBlock. ReorgSim.currentBlock increases by `BlockProgess`
	// after each call to ReorgSim.BlockNumber.
	StartBlock uint64 `json:"startBlock"`
	// BlockProgress represents how many block numbers should ReorgSim.currentBlock
	// should increase during each call to ReorgSim.BlockNumber.
	BlockProgress uint64 `json:"blockProgress"`
	// ExitBlock is used in ReorgSim.BlockNumber to return ErrExitBlockReached once its currentBlock reaches ExitBlock.
	ExitBlock uint64 `json:"exitBlock"`

	Debug bool `json:"-"`
}

// ParamV1 is embedded into ReorgSim (V1), and represents the fake blockchain client parameters
type ParamV1 struct {
	BaseParam
	ReorgEvent ReorgEvent // V1 can only perform 1 reorg event
}

// ReorgSim implements superwatcher.EthClient[block], and will be used in place of the normal Ethereum client
// to test the default emitter implementation.
// ReorgSim stores old chain and reorged chain, and it will not change the chain internal data. This means that
// it can support `MoveLogs`, even though the reorgedChain has to be reorged before being sent to NewReorgSim.
// To create a new ReorgSim with MoveLogs functionality, populate `Param.MovedLogs` with the desired values
// before calling NewReorgSimFromLogsFiles.
type ReorgSim struct {
	sync.RWMutex
	ParamV1

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
	param ParamV1,
	chain blockChain,
	reorgedChain blockChain,
	logLevel uint8,
) superwatcher.EthClient {
	if param.BlockProgress == 0 {
		panic("0 param.BlockProgress")
	}
	return &ReorgSim{
		ParamV1:           param,
		chain:             chain,
		reorgedChain:      reorgedChain,
		filterLogsCounter: make(map[uint64]int),
		debugger:          debugger.NewDebugger("ReorgSim", logLevel),
	}
}

// NewReorgSimFromLogsFiles returns a new ReorgSim with good and reorged chains mocked using the files.
func NewReorgSimFromLogsFiles(
	param ParamV1,
	logFiles []string,
	logLevel uint8,
) superwatcher.EthClient {
	chain, reorgedChain := NewBlockChainReorgMoveLogs(
		InitMappedLogsFromFiles(logFiles...),
		param.ReorgEvent,
	)

	return NewReorgSim(param, chain, reorgedChain, logLevel)
}

func (s *ReorgSim) Chain() blockChain {
	return s.chain
}

func (s *ReorgSim) ReorgedChain() blockChain {
	return s.reorgedChain
}
