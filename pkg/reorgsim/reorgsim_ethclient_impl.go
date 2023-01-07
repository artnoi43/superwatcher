package reorgsim

// See README.md for code documentation

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/batch"
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
		event := r.events[r.currentReorgEvent]
		toBeForked := blockNumber >= event.ReorgBlock
		reorgedHere := blockNumber == event.ReorgBlock

		var h common.Hash
		if toBeForked {
			h = ReorgHash(blockNumber, r.currentReorgEvent)
		} else {
			h = PRandomHash(blockNumber)
		}

		b = &Block{
			hash:        h,
			blockNumber: blockNumber,
			reorgedHere: reorgedHere,
			toBeForked:  toBeForked,
			logs:        nil,
		}
	}

	return b, nil
}

// BatchCallContext only processes `"eth_getBlockByNumber" RPC method calls.
// Each elem.Result in elems will be overwritten with *Block (implements superwatcher.BlockHeader)
// from the current chain.
func (r *ReorgSim) BatchCallContext(ctx context.Context, elems []rpc.BatchElem) error {
	r.RLock()
	defer r.RUnlock()

	for i, elem := range elems {
		if elem.Method != batch.MethodGetBlockByNumber {
			continue
		}

		// Get blockNumber from the first string argument
		bn, err := hexutil.DecodeBig(elem.Args[0].(string))
		if err != nil {
			return errors.Wrapf(err, "elems[%d] has invalid argument for method %s", i, batch.MethodGetBlockByNumber)
		}

		number := bn.Uint64()
		b, ok := r.chain[number]
		if !ok {
			eventIndex := r.currentReorgEvent
			if lenEvents := len(r.events); r.currentReorgEvent >= lenEvents {
				eventIndex = lenEvents - 1
			}
			event := r.events[eventIndex]
			toBeForked := number >= event.ReorgBlock

			var hash common.Hash
			if toBeForked {
				hash = ReorgHash(number, r.currentReorgEvent)
			} else {
				hash = PRandomHash(number)
			}

			b = &Block{
				blockNumber: number,
				hash:        hash,
				toBeForked:  toBeForked,
				reorgedHere: number == event.ReorgBlock,
			}
		}

		elems[i].Result = b
	}

	return nil
}
