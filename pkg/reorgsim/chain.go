package reorgsim

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type blockChain map[uint64]block

// NewBlockChain returns a tuple of blockChain(s). It takes in |reorgedAt|,
// and construct the chains based on that number.
func NewBlockChain(mappedLogs map[uint64][]types.Log, reorgedAt uint64) (blockChain, blockChain) {
	var found bool
	for blockNumber := range mappedLogs {
		if blockNumber == reorgedAt {
			found = true
		}
	}

	if !found {
		panic("reorgedAt block not found in any logs")
	}

	// The "good old chain"
	chain := make(blockChain)
	for blockNumber, logs := range mappedLogs {
		b := new(block)
		if blockNumber == reorgedAt {
			b.reorgedHere = true
		}

		b.blockNumber = blockNumber
		b.hash = logs[0].BlockHash
		b.logs = logs
		chain[blockNumber] = *b
	}

	// The "reorged chain" will only contain blocks after |reorgedAt|
	reorgedChain := make(blockChain)
	for blockNumber, oldBlock := range chain {
		if blockNumber >= reorgedAt {
			// Reorg this block and saves back to reorgedChain
			reorgedBlock := oldBlock.reorg()
			reorgedChain[blockNumber] = reorgedBlock
		}
	}

	return chain, reorgedChain
}
