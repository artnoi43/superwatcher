package reorgsim

import (
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const NoReorg uint64 = 0

type blockChain map[uint64]*block

// NewBlockChainNg returns a tuple of blockChain(s). It takes in |reorgedAt|,
// and construct the chains based on that number.
func NewBlockChainNg(logs []types.Log, reorgedAt uint64) (blockChain, blockChain) {
	return NewBlockChain(mapLogsToNumber(logs), reorgedAt)
}

// NewBlockChain returns a tuple of blockChain(s). It takes in |reorgedAt|,
// and construct the chains based on that number.
func NewBlockChain(mappedLogs map[uint64][]types.Log, reorgedAt uint64) (blockChain, blockChain) {
	// The "good old chain"
	chain := make(blockChain)
	for blockNumber, logs := range mappedLogs {
		chain[blockNumber] = &block{
			blockNumber: blockNumber,
			hash:        logs[0].BlockHash,
			logs:        logs,
			reorgedHere: blockNumber == reorgedAt,
			toBeForked:  blockNumber >= reorgedAt,
		}
	}

	// No reorg - use the same chain
	if reorgedAt == NoReorg {
		return chain, chain
	}

	// |reorgedChain| will differ from |oldChain| after |reorgedAt|
	reorgedChain := make(blockChain)
	for blockNumber, oldBlock := range chain {
		// Use old block for |reorgedChain| if |blockNumber| < |reorgedAt|
		if blockNumber < reorgedAt {
			reorgedChain[blockNumber] = oldBlock
			continue
		}

		reorgedChain[blockNumber] = oldBlock.reorg()
	}
	return chain, reorgedChain
}

// MoveLogs represent a move of logs to a new blockNumber
type MoveLogs struct {
	newBlock uint64
	txHashes []common.Hash // txHashes of logs to be moved to newBlock
}

func (c blockChain) reorgMoveLogs(
	// logsMoved maps old block to []moveLog
	logsMoved map[uint64][]MoveLogs,
) {
	for blockNumber, moves := range logsMoved {
		b, ok := c[blockNumber]
		if !ok {
			panic("logsMoved from non-existent block")
		}

		for _, move := range moves {
			targetBlock, ok := c[move.newBlock]
			if !ok {
				panic("logsMoved to non-existent block")
			}

			// Save logsToMove before removing it from b
			var logsToMove []types.Log
			for _, log := range b.logs {
				if gslutils.Contains(move.txHashes, log.TxHash) {
					logsToMove = append(logsToMove, log)
				}
			}

			// Remove logs from the old block
			b.removeLogs(move.txHashes)

			// Change log.BlockHash to new BlockHash
			for _, log := range logsToMove {
				log.BlockHash = targetBlock.hash
			}

			targetBlock.logs = append(targetBlock.logs, logsToMove...)
		}
	}
}
