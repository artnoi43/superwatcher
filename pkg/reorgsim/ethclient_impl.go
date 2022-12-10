package reorgsim

import (
	"context"
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	filterLogs     = "filterLogs"
	headerByNumber = "headerByNumber"
)

// chooseBlock returns a block at that |blockNumber| from
// an appropriate chain. The "reorg" logic is defined here.
func (r *ReorgSim) chooseBlock(blockNumber, fromBlock, toBlock uint64, caller string) *block {
	r.Lock()
	defer r.Unlock()
	// If we see a "toBeForked" block more than once,
	// return the reorged block from reorged chain.

	b, found := r.chain[blockNumber]
	if !found {
		return nil
	}

	log := func(prefix string) {
		r.debugger.Debug(
			2,
			fmt.Sprintf("%s chooseBlock", prefix),
			zap.Uint64("blockNumber", b.blockNumber),
			zap.String("caller", caller),
			zap.Bool("toBeForked", b.toBeForked),
			zap.Int("seen", r.filterLogsCounter[b.blockNumber]),
		)
	}

	log("<")

	// Use reorg block for this blockNumber
	if b.toBeForked {
		// In emitter.FilterLogs, client.FilterLogs is called before client.HeaderByNumber,
		// so here we use call from FilterLogs to trigger a reorg by incrementing the filterLogsCounter
		// and using reorged block if the counter is > 1.
		// This is why reorgSim.HeaderByNumber should only use reorgedChain if FilterLogs already returned
		// the reorged logs, and thus a different |n| value.
		var n int
		switch caller {
		case filterLogs:
			n = 1 // reorgSim.FilterLogs returns reorged blocks first
			// case headerByNumber:
			// 	n = 2
		default:
			panic("unexpected call to chooseBlock by \"" + caller + "\"")
		}

		if r.filterLogsCounter[blockNumber] >= n {
			b, found = r.reorgedChain[blockNumber]
			if !found {
				return nil
			}

			if !r.wasForked {
				r.debugger.Debug(1, "!REORGED!", zap.Uint64("blockNumber", blockNumber))

				r.chain = r.reorgedChain
				r.wasForked = true
			}
		}
	}

	if caller == filterLogs && b.toBeForked {
		r.filterLogsCounter[blockNumber]++
	}

	log(">")
	return b
}

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

	var logs []types.Log
	for blockNumber := from; blockNumber <= to; blockNumber++ {
		// Choose a block from an appropriate chain
		b := r.chooseBlock(blockNumber, from, to, filterLogs)
		if b == nil {
			continue
		}

		for _, log := range b.logs {
			if query.Addresses == nil && query.Topics == nil {
				logs = append(logs, log)
				continue
			}

			if query.Addresses != nil {
				if query.Topics != nil {
					if gslutils.Contains(query.Topics[0], log.Topics[0]) {
						logs = append(logs, log)
						continue
					}
				}

				if gslutils.Contains(query.Addresses, log.Address) {
					logs = append(logs, log)
					continue
				}
			}
		}
	}

	return logs, nil
}

func (r *ReorgSim) BlockNumber(ctx context.Context) (uint64, error) {
	r.Lock()
	defer r.Unlock()

	if r.Param.BlockProgress == 0 {
		panic("0 BlockProgress")
	}

	if r.currentBlock == 0 {
		r.currentBlock = r.Param.StartBlock
		return r.currentBlock, nil
	}

	currentBlock := r.currentBlock
	if currentBlock >= r.Param.ExitBlock {
		return currentBlock, errors.Wrapf(ErrExitBlockReached, "exit block %d reached", r.Param.ExitBlock)
	}

	r.currentBlock = currentBlock + r.Param.BlockProgress
	return currentBlock, nil
}

// func (r *ReorgSim) HeaderByNumber(ctx context.Context, number *big.Int) (superwatcher.BlockHeader, error) {
// 	blockNumber := number.Uint64()
//
// 	b := r.chooseBlock(blockNumber, blockNumber, blockNumber, headerByNumber)
// 	if b != nil {
// 		return *b, nil
// 	}
//
// 	return &block{
// 		hash: PRandomHash(number.Uint64()),
// 	}, nil
// }
