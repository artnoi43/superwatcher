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
	// param specifies how ReorgSim should behave
	param Param
	// events represents the multiple events
	events []ReorgEvent
	// currentReorgEvent is used to index events and reorgedChains
	currentReorgEvent int
	// triggered is the index of current triggered ReorgSim.reorgedChains
	triggered int
	// forked tracks whether reorgChains[i] was forked (used).
	forked int
	// chain is the original blockChain
	chain BlockChain
	// reorgedChains is the multiple reorged blockchains construct from `ReorgSim.chain` and `ReorgSim.param`
	reorgedChains []BlockChain
	// currentBlock tracks the current block for the fake blockChain and is used for exclusively in BlockByNumber
	currentBlock uint64

	debugger *debugger.Debugger
}

func NewReorgSim(
	param Param,
	events []ReorgEvent,
	chain BlockChain,
	reorgedChains []BlockChain,
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
		param:         param,
		events:        validatedEvents,
		chain:         chain,
		reorgedChains: reorgedChains,
		triggered:     -1, // can't be 0
		forked:        -1, // can't be 0
		debugger:      debugger.NewDebugger(name, logLevel),
	}, nil
}

func NewReorgSimFromLogs(
	param Param,
	events []ReorgEvent,
	logs map[uint64][]types.Log,
	debugName string,
	logLevel uint8,
) (
	superwatcher.EthClient,
	error,
) {
	chain, reorgedChains := NewBlockChain(logs, events)
	return NewReorgSim(param, events, chain, reorgedChains, debugName, logLevel)
}

func NewReorgSimFromLogsFiles(
	param Param,
	events []ReorgEvent,
	logsFiles []string,
	debugName string,
	logLevel uint8,
) (
	superwatcher.EthClient,
	error,
) {
	return NewReorgSimFromLogs(
		param,
		events,
		InitMappedLogsFromFiles(logsFiles...),
		debugName,
		logLevel,
	)
}

func (r *ReorgSim) Chain() BlockChain { //nolint:revive
	return r.chain
}

func (r *ReorgSim) ReorgedChains() []BlockChain { //nolint:revive
	return r.reorgedChains
}

func (r *ReorgSim) ReorgedChain(i int) BlockChain { //nolint:revive
	return r.reorgedChains[i]
}
