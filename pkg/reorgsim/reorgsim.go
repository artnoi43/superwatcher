package reorgsim

import (
	"sync"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// BaseParam is the basic parameters for the mock client. Chain reorg parameters are NOT included here.
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

var DefaultParam = BaseParam{
	BlockProgress: 20,
	Debug:         true,
}

// ReorgEvent is parameters for chain reorg events.
type ReorgEvent struct {
	// ReorgBlock is the pivot block after which ReorgSim should use ReorgSim.reorgedChain.
	ReorgBlock uint64 `json:"reorgBlock"`
	// MovedLogs represents all of the moved logs after a chain reorg event. The map key is the source blockNumber.
	MovedLogs map[uint64][]MoveLogs `json:"movedLogs"`
}

var errInvalidReorgEvents = errors.New("invalid reorg events")

// ReorgSim is a mock superwatcher.EthClient can simulate multiple on-the-fly chain reorganizations.
type ReorgSim struct {
	sync.RWMutex
	// param is only BaseParam, which specifies how ReorgSim should behave
	param BaseParam
	// events represents the multiple events
	events []ReorgEvent
	// currentReorgEvent is used to index events and reorgedChains
	currentReorgEvent int
	// chain is the original blockChain
	chain blockChain
	// reorgedChains is the multiple reorged blockchains construct from `ReorgSim.chain` and `ReorgSim.param`
	reorgedChains []blockChain
	// forked tracks whether reorgChains[i] was forked (used)
	forked []bool
	// currentBlock tracks the current block for the fake blockChain and is used for exclusively in BlockByNumber
	currentBlock uint64
	// filterLogsCounter is used to switch chain for a blockNumber, after certain number of calls to FilterLogs
	filterLogsCounter map[uint64]int

	debugger *debugger.Debugger
}

func newReorgSim(
	param BaseParam,
	events []ReorgEvent,
	chain blockChain,
	reorgedChains []blockChain,
	debugName string,
	logLevel uint8,
) (
	*ReorgSim,
	error,
) {
	if err := validateReorgEvents(events); err != nil {
		return nil, errors.Wrap(err, "invalid events")
	}

	var name string
	if debugName == "" {
		name = "ReorgSim"
	} else {
		name = "ReorgSim " + debugName
	}

	return &ReorgSim{
		param:             param,
		events:            events,
		chain:             chain,
		reorgedChains:     reorgedChains,
		currentReorgEvent: 0,
		forked:            make([]bool, len(events)),
		filterLogsCounter: make(map[uint64]int),
		debugger:          debugger.NewDebugger(name, logLevel),
	}, nil
}

func NewReorgSim(
	param BaseParam,
	events []ReorgEvent,
	logs map[uint64][]types.Log,
	debugName string,
	logLevel uint8,
) (
	superwatcher.EthClient,
	error,
) {
	chain, reorgedChains := NewBlockChain(logs, events)
	return newReorgSim(param, events, chain, reorgedChains, debugName, logLevel)
}

func NewReorgSimFromLogsFiles(
	param BaseParam,
	events []ReorgEvent,
	logsFiles []string,
	debugName string,
	logLevel uint8,
) (
	superwatcher.EthClient,
	error,
) {
	return NewReorgSim(
		param,
		events,
		InitMappedLogsFromFiles(logsFiles...),
		debugName,
		logLevel,
	)
}

func (r *ReorgSim) Chain() blockChain {
	return r.chain
}

func (r *ReorgSim) ReorgedChains() []blockChain {
	return r.reorgedChains
}

func (r *ReorgSim) ReorgedChain(i int) blockChain {
	return r.reorgedChains[i]
}

// Subsequent ReorgEvent.ReorgBlock should be larger than the previous ones
func validateReorgEvents(events []ReorgEvent) error {
	reorgBlocks := gslutils.Map(events, func(event ReorgEvent) (uint64, bool) {
		return event.ReorgBlock, true
	})

	for i, reorgBlock := range reorgBlocks {
		var prevReorgBlock uint64
		if i == 0 {
			continue
		} else {
			prevReorgBlock = reorgBlocks[i-1]
		}

		if prevReorgBlock > reorgBlock {
			return errors.Wrapf(
				errInvalidReorgEvents, "event at index %d has smaller value than index %d (%d > %d)",
				i, i-1, reorgBlock, prevReorgBlock,
			)
		}
	}

	return nil
}
