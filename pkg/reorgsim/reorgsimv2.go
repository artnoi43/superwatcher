package reorgsim

import (
	"sync"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

var errInvalidReorgEvents = errors.New("invalid reorg events")

// ReorgSimV2 can perform multiple on-the-fly chain reorganizations.
type ReorgSimV2 struct {
	sync.RWMutex
	// param is only BaseParam, which specifies how ReorgSimV2 should behave
	param BaseParam
	// events represents the multiple events
	events []ReorgEvent
	// currentReorgEvent is used to index events and reorgedChains
	currentReorgEvent int
	// chain is the original blockChain
	chain blockChain
	// reorgedChains is the multiple reorged blockchains construct from `ReorgSimV2.chain` and `ReorgSimV2.param`
	reorgedChains []blockChain
	// forked tracks whether reorgChains[i] was forked (used)
	forked []bool
	// currentBlock tracks the current block for the fake blockChain and is used for exclusively in BlockByNumber
	currentBlock uint64
	// filterLogsCounter is used to switch chain for a blockNumber, after certain number of calls to FilterLogs
	filterLogsCounter map[uint64]int
	// forked (forked[i]) tracks if the param[i] was forked (i.e. if ReorgChain has returned any logs from the reorgedEvents[i])

	debugger *debugger.Debugger
}

func newReorgSimV2(
	param BaseParam,
	events []ReorgEvent,
	chain blockChain,
	reorgedChains []blockChain,
	logLevel uint8,
) (
	*ReorgSimV2,
	error,
) {
	if err := validateReorgEvents(events); err != nil {
		return nil, errors.Wrap(err, "invalid events")
	}

	return &ReorgSimV2{
		param:             param,
		events:            events,
		chain:             chain,
		reorgedChains:     reorgedChains,
		currentReorgEvent: 0,
		forked:            make([]bool, len(events)),
		filterLogsCounter: make(map[uint64]int),
		debugger:          debugger.NewDebugger("ReorgSimV2", logLevel),
	}, nil
}

// NewReorgSimV2 constructs blockChains using NewBlockChainV2 to call newReorgSimV2.
func NewReorgSimV2(
	param BaseParam,
	events []ReorgEvent,
	logs map[uint64][]types.Log,
	logLevel uint8,
) (
	superwatcher.EthClient,
	error,
) {
	chain, reorgedChains := NewBlockChainReorgV2(logs, events)
	return newReorgSimV2(param, events, chain, reorgedChains, logLevel)
}

func NewReorgSimV2FromLogsFiles(
	param BaseParam,
	events []ReorgEvent,
	logsFiles []string,
	logLevel uint8,
) (
	superwatcher.EthClient,
	error,
) {
	return NewReorgSimV2(
		param,
		events,
		InitMappedLogsFromFiles(logsFiles...),
		logLevel,
	)
}

func (r *ReorgSimV2) Chain() blockChain {
	return r.chain
}

func (r *ReorgSimV2) ReorgedChains() []blockChain {
	return r.reorgedChains
}

func (r *ReorgSimV2) ReorgedChain(i int) blockChain {
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
