package reorgsim

import (
	"github.com/ethereum/go-ethereum/core/types"
)

const NoReorg uint64 = 0

type blockChain map[uint64]block

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
		chain[blockNumber] = block{
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
