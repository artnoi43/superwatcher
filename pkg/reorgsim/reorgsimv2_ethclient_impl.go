package reorgsim

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (r *ReorgSimV2) chooseBlock(blockNumber uint64, caller string) *block {
	var currentReorgEvent int
	for i, forked := range r.forked {
		if !forked {
			currentReorgEvent = i
			break
		}
	}

	b, found := r.chain[blockNumber]
	if !found || b == nil {
		return nil
	}

	logFunc := func(prefix string) {
		r.debugger.Debug(
			2,
			fmt.Sprintf("%s chooseBlock", prefix),
			zap.Int("currentReorgEvent", currentReorgEvent),
			zap.Uint64("blockNumber", b.blockNumber),
			zap.String("caller", caller),
			zap.Bool("toBeForked", b.toBeForked),
			zap.Int("seen", r.filterLogsCounter[b.blockNumber]),
		)
	}

	logFunc("<")

	if b.toBeForked {
		var n int
		switch caller {
		case methodFilterLogs:
			n = 1 // reorgSim.FilterLogs returns reorged blocks first
			// case headerByNumber:
			// 	n = 2
		default:
			panic("unexpected call to chooseBlock by \"" + caller + "\"")
		}

		// Only trigger new reorg if filterLogsCounter is >= n
		if r.filterLogsCounter[blockNumber] >= n {
			currentChain := r.reorgedChains[currentReorgEvent]
			b, found = currentChain[blockNumber]
			if !found {
				return nil
			}

			r.debugger.Debug(
				1, "REORGED!",
				zap.Uint64("blockNumber", blockNumber),
				zap.Bools("forked", r.forked),
			)

			r.forked[currentReorgEvent] = true
		}
	}

	if caller == methodFilterLogs && b.toBeForked {
		r.filterLogsCounter[blockNumber]++
	}

	logFunc(">")
	return b
}

func (r *ReorgSimV2) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
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

	var logs []types.Log
	for number := from; number <= to; number++ {
		b := r.chooseBlock(number, methodFilterLogs)
		if b == nil || len(b.logs) == 0 {
			continue
		}

		appendFilterLogs(&b.logs, &logs, query.Addresses, query.Topics)
	}

	return logs, nil
}

func (r *ReorgSimV2) BlockNumber(ctx context.Context) (uint64, error) {
	r.Lock()
	defer r.Unlock()

	if r.currentBlock == 0 {
		r.currentBlock = r.param.StartBlock
		return r.currentBlock, nil
	}

	currentBlock := r.currentBlock
	if currentBlock >= r.param.ExitBlock {
		return currentBlock, errors.Wrapf(ErrExitBlockReached, "exit block %d reached", r.param.ExitBlock)
	}

	r.currentBlock = currentBlock + r.param.BlockProgress
	return currentBlock, nil
}
