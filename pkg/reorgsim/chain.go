package reorgsim

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const NoReorg uint64 = 0

type blockChain map[uint64]*block

// MoveLogs represent a move of logs to a new blockNumber
type MoveLogs struct {
	NewBlock uint64
	TxHashes []common.Hash // txHashes of logs to be moved to newBlock
}

// reorg calls `*block.reorg` on every block whose blockNumber is greater than |reorgedAt|.
// Unlike `*block.reorg`, which returns a `*block`, c.reorg(reorgedAt) modifies c in-place.
func (c blockChain) reorg(reorgedBlock uint64) {
	for number, block := range c {
		if number >= reorgedBlock {
			c[number] = block.reorg()
		}
	}
}

// moveLogs moved will have their blockHash and blockNumber changed to destination blocks.
// If you are manually reorging and moving logs, call blockChain.reorg before blockChain.moveLogs.
// If you are creating new reorged chain with moved logs, use NewBlockChainReorgMoveLogs instead,
// as blockChain.moveLogs only moves the logs and changing log.BlockHash and log.BlockNumber.
// It also returns 2 slices of block numbers, 1st of which is a slice of blocks from which logs are moved,
// the 2nd of which is a slice of blocks to which logs are moved.
func (c blockChain) moveLogs(
	movedLogs map[uint64][]MoveLogs,
) (
	[]uint64, // Blocks from which logs are moved from
	[]uint64, // Blocks to which logs art moved to
) {
	// A slice of unique blockNumbers that logs will be moved to.
	// Might be useful to caller, maybe to create empty blocks (no logs) for the old chain
	var moveFromBlocks []uint64
	var moveToBlocks []uint64

	for moveFromBlock, moves := range movedLogs {
		b, ok := c[moveFromBlock]
		if !ok {
			panic("logsMoved from non-existent block")
		}

		moveFromBlocks = append(moveFromBlocks, moveFromBlock)

		for _, move := range moves {
			targetBlock, ok := c[move.NewBlock]
			if !ok {
				panic(fmt.Sprintf("logsMoved to non-existent block %d", move.NewBlock))
			}

			// Add unique moveToBlocks
			if !gslutils.Contains(moveToBlocks, move.NewBlock) {
				moveToBlocks = append(moveToBlocks, move.NewBlock)
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

	return moveFromBlocks, moveToBlocks
}

// NewBlockChain returns a new blockChain from |mappedLogs|. The parameter |reorgedAt|
// is used to deterine block.reorgedHere and block.toBeForked
func NewBlockChain(
	mappedLogs map[uint64][]types.Log,
	reorgedBlock uint64,
) blockChain {
	chain := make(blockChain)

	for blockNumber, logs := range mappedLogs {
		chain[blockNumber] = &block{
			blockNumber: blockNumber,
			hash:        logs[0].BlockHash,
			logs:        logs,
			reorgedHere: blockNumber == reorgedBlock,
			toBeForked:  blockNumber >= reorgedBlock,
		}
	}

	return chain
}

// NewBlockChainReorgV1 returns a tuple of blockChain(s) for reorgSim. It takes in |reorgedAt|,
// and construct the chains based on that number.
func NewBlockChainReorgV1(
	mappedLogs map[uint64][]types.Log,
	reorgedBlock uint64,
) (
	blockChain,
	blockChain,
) {
	// The "good old chain"
	chain := NewBlockChain(mappedLogs, reorgedBlock)

	// No reorg - use the same chain
	if reorgedBlock == NoReorg {
		return chain, chain
	}

	// |reorgedChain| will differ from |oldChain| after |reorgedAt|
	reorgedChain := copyBlockChain(chain)
	reorgedChain.reorg(reorgedBlock)

	return chain, reorgedChain
}

// NewBlockChainReorgMoveLogs is similar to NewBlockChain, but the reorgedChain (the 2nd returned variable)
// will call blockChain.reorgMoveLogs with |movedLogs|. NewBlockChainReorgMoveLogs will also add moveToBlock
// in the old original chain to ensure that ReorgSim does not skip the block.
// Calling NewBlockChainReorgMoveLogs with nil |movedLogs| is the same as calling NewBlockChain.
func NewBlockChainReorgMoveLogs(
	mappedLogs map[uint64][]types.Log,
	event ReorgEvent,
) (
	blockChain,
	blockChain,
) {
	for blockNumber := range event.MovedLogs {
		if blockNumber < event.ReorgBlock {
			panic(fmt.Sprintf("blockNumber %d < reorgedAt %d", blockNumber, event.ReorgBlock))
		}
	}

	chain := NewBlockChain(mappedLogs, event.ReorgBlock)

	reorgedChain := copyBlockChain(chain)
	reorgedChain.reorg(event.ReorgBlock)

	if len(event.MovedLogs) != 0 {
		_, moveToBlocks := reorgedChain.moveLogs(event.MovedLogs)

		// Ensure that all moveToBlocks exist in original chain
		for _, moveToBlock := range moveToBlocks {
			// If the old chain did not have moveToBlock, create one
			if _, ok := chain[moveToBlock]; !ok {
				chain[moveToBlock] = &block{
					blockNumber: moveToBlock,
					hash:        RandomHash(moveToBlock),
					reorgedHere: moveToBlock == event.ReorgBlock,
					toBeForked:  true,
				}
			}
		}
	}

	return chain, reorgedChain
}

// NewBlockChainReorgV2 can simulate multiple chain reorgs with moved logs.
// Each ReorgEvent will result in its own blockChain, with the same index.
func NewBlockChainReorgV2(
	logs map[uint64][]types.Log,
	events []ReorgEvent,
) (
	blockChain, // Original chain
	[]blockChain, // Reorged chains
) {
	if len(events) == 0 {
		return NewBlockChain(logs, NoReorg), nil
	}

	chain := NewBlockChain(logs, events[0].ReorgBlock)

	var reorgedChains = make([]blockChain, len(events))
	for i, event := range events {
		var prevChain blockChain
		if i == 0 {
			prevChain = chain
		} else {
			prevChain = reorgedChains[i-1]
		}

		forkedChain := copyBlockChain(prevChain)

		// Reorg and move logs
		forkedChain.reorg(event.ReorgBlock)
		moveFromBlocks, moveToBlocks := forkedChain.moveLogs(event.MovedLogs)

		// Make sure the movedFrom block is not nil in forkedChain
		for _, prevFrom := range moveFromBlocks {
			if _, ok := prevChain[prevFrom]; !ok {
				panic(fmt.Sprintf("moved from non-existent block %d in the old chain", prevFrom))
			}

			if b, ok := forkedChain[prevFrom]; !ok || b == nil {
				forkedChain[prevFrom] = &block{
					blockNumber: prevFrom,
					hash:        RandomHash(prevFrom), // Uses non-deterministic hash
					reorgedHere: prevFrom == event.ReorgBlock,
					toBeForked:  true,
				}
			}
		}

		for _, forkedTo := range moveToBlocks {
			if _, ok := forkedChain[forkedTo]; !ok {
				panic(fmt.Sprintf("moved to non-existent block %d in the new chain", forkedTo))
			}

			if _, ok := prevChain[forkedTo]; !ok {
				prevChain[forkedTo] = &block{
					blockNumber: forkedTo,
					hash:        RandomHash(forkedTo),
					reorgedHere: forkedTo == event.ReorgBlock,
					toBeForked:  true,
				}
			}
		}

		// Make sure the block from which the logs moved
		reorgedChains[i] = forkedChain
	}

	return chain, reorgedChains
}
