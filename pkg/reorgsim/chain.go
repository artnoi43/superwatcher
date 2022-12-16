package reorgsim

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const NoReorg uint64 = 0

type blockChain map[uint64]*block

// NewBlockChain
func NewBlockChain(reorgedAt uint64, logs ...types.Log) (blockChain, blockChain) {
	return newBlockChain(mapLogsToNumber(logs), reorgedAt)
}

// newBlockChain returns a tuple of blockChain(s) for reorgSim. It takes in |reorgedAt|,
// and construct the chains based on that number.
func newBlockChain(mappedLogs map[uint64][]types.Log, reorgedAt uint64) (blockChain, blockChain) {
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
	NewBlock uint64
	TxHashes []common.Hash // txHashes of logs to be moved to newBlock
}

// NewBlockChainWithMovedLogs is similar to NewBlockChain, but the reorgedChain (the 2nd returned variable)
// will call blockChain.reorgMoveLogs with |movedLogs|.
// Calling NewBlockChainWithMovedLogs with nil |movedLogs| is the same as calling NewBlockChain
func NewBlockChainWithMovedLogs(
	mappedLogs map[uint64][]types.Log,
	reorgedAt uint64,
	movedLogs map[uint64][]MoveLogs,
) (
	blockChain,
	blockChain,
) {
	for blockNumber := range movedLogs {
		if blockNumber < reorgedAt {
			panic(fmt.Sprintf("blockNumber %d < reorgedAt %d", blockNumber, reorgedAt))
		}
	}

	chain, reorgedChain := newBlockChain(mappedLogs, reorgedAt)

	if movedLogs != nil || len(movedLogs) != 0 {
		reorgedChain.reorgMoveLogs(movedLogs)
	}

	return chain, reorgedChain
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
			targetBlock, ok := c[move.NewBlock]
			if !ok {
				panic("logsMoved to non-existent block")
			}

			// Save logsToMove before removing it from b
			var logsToMove []types.Log
			for _, log := range b.logs {
				if gslutils.Contains(move.TxHashes, log.TxHash) {
					logsToMove = append(logsToMove, log)
				}
			}

			// Remove logs from the old block
			b.removeLogs(move.TxHashes)

			// Change log.BlockHash to new BlockHash
			for i := range logsToMove {
				logsToMove[i].BlockNumber = targetBlock.blockNumber
				logsToMove[i].BlockHash = targetBlock.hash
			}

			targetBlock.logs = append(targetBlock.logs, logsToMove...)
		}
	}
}
