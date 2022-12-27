package reorgsim

// See README.md for code documentation

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (r *ReorgSim) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	if query.FromBlock == nil {
		return nil, errors.New("nil query.FromBlock")
	}
	if query.ToBlock == nil {
		return nil, errors.New("nil query.ToBlock")
	}

	from := query.FromBlock.Uint64()
	to := query.ToBlock.Uint64()

	if from > to {
		to = from
	}

	r.triggerForkChain(from, to)

	// See if there's unforked r.reorgedChains before forking
	for _, chainForked := range r.forked {
		if !chainForked {
			if r.triggered == r.currentReorgEvent {
				r.forkChain(from, to)
			}
		}
	}

	var logs []types.Log
	for number := from; number <= to; number++ {
		b := r.chain[number]
		if b == nil || len(b.logs) == 0 {
			continue
		}

		r.debugger.Debug(
			3, "FilterLogs block",
			zap.Uint64("blockNumber", b.blockNumber),
			zap.String("blockHash", b.hash.String()),
			zap.Bool("toBeForked", b.toBeForked),
			zap.Bool("forked", b.reorgedHere),
		)

		appendFilterLogs(&b.logs, &logs, query.Addresses, query.Topics)
	}

	return logs, nil
}

func (r *ReorgSim) BlockNumber(ctx context.Context) (uint64, error) {
	r.Lock()
	defer r.Unlock()

	if r.currentBlock == 0 {
		r.currentBlock = r.param.StartBlock
		return r.currentBlock, nil
	}

	currentBlock := r.currentBlock // currentBlock will be returned
	if currentBlock >= r.param.ExitBlock {
		return currentBlock, errors.Wrapf(ErrExitBlockReached, "exit block %d reached", r.param.ExitBlock)
	}

	r.currentBlock = currentBlock + r.param.BlockProgress
	return currentBlock, nil
}

// triggerForkChain updates `r.triggered` to true
// if the current ReorgEvent.ReorgTrigger is within range [from, to]
func (r *ReorgSim) triggerForkChain(rangeStart, rangeEnd uint64) {
	if len(r.events) == 0 {
		return
	}

	if r.currentReorgEvent >= len(r.events) {
		return
	}

	event := r.events[r.currentReorgEvent]
	// If event.ReorgBlock is within range, then mark as current event as triggered
	if rangeStart <= event.ReorgTrigger && rangeEnd >= event.ReorgTrigger {
		r.triggered = r.currentReorgEvent
	}
}

// forkChain performs chain reorg logic if the current ReorgEvent.ReorgBlock is within range [from, to]
// and if r.seen[ReorgEvent.ReorgBlock] is >= 1. The latter check allows for the poller/emitter to see
// pre-fork block hash once, so that we can test the poller/emitter logic.
func (r *ReorgSim) forkChain(fromBlock, toBlock uint64) {
	event := r.events[r.currentReorgEvent]

	if fromBlock <= event.ReorgBlock && toBlock >= event.ReorgBlock {
		r.seen[event.ReorgBlock]++

		if r.seen[event.ReorgBlock] < 1 {
			return
		}

		r.debugger.Debug(
			1, "REORG!",
			zap.Int("eventIndex", r.currentReorgEvent),
			zap.Uint64("currentBlock", r.currentBlock),
			zap.Uint64("reorgTrigger", event.ReorgTrigger),
			zap.Uint64("reorgBlock", event.ReorgBlock),
		)

		var currentChain blockChain
		var lastReorg bool

		if r.currentReorgEvent < len(r.events) {
			currentChain = r.reorgedChains[r.currentReorgEvent]
		} else {
			lastReorg = true
			currentChain = r.reorgedChains[len(r.reorgedChains)-1]
		}

		if !lastReorg {
			if !r.forked[r.currentReorgEvent] {
				r.chain = currentChain
				r.forked[r.currentReorgEvent] = true
				r.currentBlock = event.ReorgBlock
				r.seen = make(map[uint64]int)

				r.currentReorgEvent++
			}
		}

		r.debugger.Debug(
			1, "forkChain",
			zap.Uint64("reorgTrigger blockNumber", event.ReorgTrigger),
			zap.Uint64("reorgedBlock blockNumber", event.ReorgBlock),
			zap.Int("currentReorgEvent", r.currentReorgEvent),
			zap.Bools("forked", r.forked),
		)
	}
}
