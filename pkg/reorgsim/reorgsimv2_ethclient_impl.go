package reorgsim

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (r *ReorgSimV2) chooseBlock(number uint64) *block {

	return nil
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

	for number := from; number <= to; number++ {

	}

	return nil, nil
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
