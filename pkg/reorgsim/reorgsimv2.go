package reorgsim

import (
	"fmt"
	"sync"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
	"github.com/ethereum/go-ethereum/core/types"
)

// ReorgSimV2 can perform multiple on-the-fly chain reorganizations.
type ReorgSimV2 struct {
	sync.RWMutex
	// param is only BaseParam, which specifies how ReorgSimV2 should behave
	param BaseParam
	// events represents the multiple events
	events []ReorgEvent
	// chain is the original blockChain
	chain blockChain
	// reorgedChains is the multiple reorged blockchains construct from `ReorgSimV2.chain` and `ReorgSimV2.param`
	reorgedChains []blockChain
	// currentBlock tracks the current block for the fake blockChain and is used for exclusively in BlockByNumber
	currentBlock uint64
	// filterLogsCounter is used to switch chain for a blockNumber, after certain number of calls to FilterLogs
	filterLogsCounter map[uint64]int
	// forked (forked[i]) tracks if the param[i] was forked (i.e. if ReorgChain has returned any logs from the reorgedEvents[i])
	forked []bool

	debugger *debugger.Debugger
}

func newReorgSimV2(
	param BaseParam,
	reorgEvents []ReorgEvent,
	chain blockChain,
	reorgedChains []blockChain,
	logLevel uint8,
) *ReorgSimV2 {
	return &ReorgSimV2{
		param:             param,
		events:            reorgEvents,
		chain:             chain,
		reorgedChains:     reorgedChains,
		filterLogsCounter: make(map[uint64]int),
		forked:            make([]bool, len(reorgEvents)),
		debugger:          debugger.NewDebugger("ReorgSimV2", logLevel),
	}
}

// NewReorgSimV2 uses params to construct multiple reorged chains. It uses `params[0]`.StartBlock as ReorgSimV2.currentBlock
func NewReorgSimV2(
	param BaseParam,
	events []ReorgEvent,
	logs map[uint64][]types.Log,
	logLevel uint8,
) superwatcher.EthClient {
	chain, _ := NewBlockChainWithMovedLogs(logs, events[0])

	var reorgedChains = make([]blockChain, len(events))
	for i, event := range events {
		var prevChain blockChain
		if i == 0 {
			prevChain = chain
		} else {
			prevChain = reorgedChains[i-1]
		}

		forkedChain := copyBlockChain(prevChain)
		prevChainFromBlocks := forkedChain.reorgMoveLogs(event.MovedLogs)

		for _, prevFrom := range prevChainFromBlocks {
			if _, ok := prevChain[prevFrom]; !ok {
				panic(fmt.Sprintf("moved from non-existent block %d", prevFrom))
			}

			// Make sure the movedFrom block is not nil
			if b, ok := forkedChain[prevFrom]; !ok || b == nil {
				forkedChain[prevFrom] = &block{
					blockNumber: prevFrom,
					hash:        RandomHash(prevFrom),
					reorgedHere: prevFrom == event.ReorgBlock,
					toBeForked:  true,
				}

			}
		}

		// Make sure the block from which the logs moved
		reorgedChains[i] = forkedChain
	}

	return newReorgSimV2(param, events, chain, reorgedChains, logLevel)
}

func NewReorgSimV2FromLogsFiles(
	param BaseParam,
	events []ReorgEvent,
	logsFiles []string,
	logLevel uint8,
) superwatcher.EthClient {
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
