package reorgsim

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

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
	// triggered tracks whether a corresponding ReorgSim.reorgedChains was triggered
	triggered []bool
	// forked tracks whether reorgChains[i] was forked (used)
	forked []bool
	// seen tracks ReorgEvent.ReorgBlock
	seen map[uint64]int
	// currentBlock tracks the current block for the fake blockChain and is used for exclusively in BlockByNumber
	currentBlock uint64

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
	validatedEvents, err := validateReorgEvent(events)
	if err != nil {
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
		events:            validatedEvents,
		chain:             chain,
		reorgedChains:     reorgedChains,
		currentReorgEvent: 0,
		triggered:         make([]bool, len(events)),
		forked:            make([]bool, len(events)),
		seen:              make(map[uint64]int),
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

func (r *ReorgSim) Chain() blockChain { //nolint:revive
	return r.chain
}

func (r *ReorgSim) ReorgedChains() []blockChain { //nolint:revive
	return r.reorgedChains
}

func (r *ReorgSim) ReorgedChain(i int) blockChain { //nolint:revive
	return r.reorgedChains[i]
}
