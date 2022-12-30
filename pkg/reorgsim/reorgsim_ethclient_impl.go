package reorgsim

// See README.md for code documentation

// TODO: new exit strategy, and reimplement HeaderByNumber

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
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

	// FilterLogs triggers chain reorg with fromBlock, toBlock
	r.triggerForkChain(from, to)

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
			zap.Uint64("currentBlock", r.currentBlock),
			zap.Bool("toBeForked", b.toBeForked),
			zap.Bool("reorgedHere", b.reorgedHere),
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

	if r.currentBlock >= r.param.ExitBlock {
		return r.currentBlock, errors.Wrapf(ErrExitBlockReached, "exit block %d reached", r.param.ExitBlock)
	}

	currentBlock := r.currentBlock
	r.currentBlock += r.param.BlockProgress
	return currentBlock, nil
}

func (r *ReorgSim) HeaderByNumber(ctx context.Context, number *big.Int) (superwatcher.BlockHeader, error) {
	blockNumber := number.Uint64()
	b, ok := r.chain[blockNumber]
	if !ok {
		var h common.Hash
		if blockNumber >= r.events[r.currentReorgEvent].ReorgBlock {
			h = ReorgHash(blockNumber, r.currentReorgEvent)
		} else {
			h = PRandomHash(blockNumber)
		}

		return &Block{
			hash:        h,
			blockNumber: blockNumber,
		}, nil
	}

	return b, nil
}
