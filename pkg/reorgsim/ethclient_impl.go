package reorgsim

import (
	"context"
	"fmt"
	"math/big"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
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

	fmt.Println("< chooseBlock:", blockNumber, caller, "toBeForked:", b.toBeForked, "seen:", r.seenFilterLogs[blockNumber])

	// Use reorg block for this blockNumber
	if b.toBeForked {
		var n int
		switch caller {
		case filterLogs:
			n = 1
		case headerByNumber:
			n = 2
		}

		if r.seenFilterLogs[blockNumber] >= n {
			b, found = r.reorgedChain[blockNumber]
			if !found {
				return nil
			}

			if !r.forked {
				fmt.Println("REORG!", blockNumber)
				r.chain = r.reorgedChain
				r.forked = true
			}
		}
	}

	if caller == filterLogs && b.toBeForked {
		r.seenFilterLogs[blockNumber]++
	}

	fmt.Println("> chooseBlock:", blockNumber, caller, "toBeForked:", b.toBeForked, "seen:", r.seenFilterLogs[blockNumber])

	return &b
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

	if r.ReorgParam.BlockProgress == 0 {
		panic("0 BlockProgress")
	}

	if r.ReorgParam.currentBlock == 0 {
		r.ReorgParam.currentBlock = r.ReorgParam.StartBlock
		return r.currentBlock, nil
	}

	currentBlock := r.ReorgParam.currentBlock
	if currentBlock >= r.ReorgParam.ExitBlock {
		return currentBlock, errors.Wrapf(ErrExitBlockReached, "exit block %d reached", r.ReorgParam.ExitBlock)
	}

	r.ReorgParam.currentBlock = currentBlock + r.ReorgParam.BlockProgress
	return currentBlock, nil
}

func (r *ReorgSim) HeaderByNumber(ctx context.Context, number *big.Int) (superwatcher.BlockHeader, error) {
	blockNumber := number.Uint64()

	b := r.chooseBlock(blockNumber, blockNumber, blockNumber, headerByNumber)
	if b != nil {
		return *b, nil
	}

	block := &block{
		// We only need hash here because caller will only call superwatcher.EthClient.Hash()
		hash: deterministicRandomHash(blockNumber),
	}

	return block, nil
}
