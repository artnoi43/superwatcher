package reorgsim

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// chooseBlock returns a block at that |blockNumber| from
// an appropriate chain. The "reorg" logic is defined here.
func (r *reorgSim) chooseBlock(blockNumber uint64) *block {

	// If we see a "toBeForked" block more than once,
	// return the reorged block from reorged chain.

	b, found := r.chain[blockNumber]
	if !found {
		return nil
	}

	// Use reorg block for this blockNumber
	if b.toBeForked && r.seen[blockNumber] > 0 {
		b, found = r.reorgedChain[blockNumber]
		if !found {
			return nil
		}
	}

	return &b
}

func (r *reorgSim) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
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
		// Choose appropriate block
		b := r.chooseBlock(blockNumber)
		if b == nil {
			continue
		}

		logs = append(logs, b.logs...)
		r.seen[blockNumber]++
	}

	return logs, nil
}

func (r *reorgSim) BlockNumber(ctx context.Context) (uint64, error) {
	return 20000000, nil
}

func (r *reorgSim) HeaderByNumber(ctx context.Context, number *big.Int) (block, error) {
	blockNumber := number.Uint64()
	b := r.chooseBlock(blockNumber)

	return *b, nil
}
