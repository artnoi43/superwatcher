package reorgsim

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Use this as ReorgEvent.ReorgBlock to disable chain reorg.
const NoReorg uint64 = 0

type BlockChain map[uint64]*Block

// MoveLogs represent a move of logs to a new blockNumber
type MoveLogs struct {
	NewBlock uint64
	TxHashes []common.Hash // txHashes of logs to be moved to newBlock
}

// reorg calls `*block.reorg` on every block whose blockNumber is greater than |reorgedBlock|.
// |reorgIndex| is used to generate different hash for the same block in different ReorgEvent.
// Unlike `*block.reorg`, which returns a `*block`, `c.reorg` modifies c in-place.
func (chain BlockChain) reorg(reorgedBlock uint64, reorgIndex int) {
	for number, block := range chain {
		if number >= reorgedBlock {
			chain[number] = block.reorg(reorgIndex)
		}
	}
}

// moveLogs moved will have their blockHash and blockNumber changed to destination blocks.
// If you are manually reorging and moving logs, call blockChain.reorg before blockChain.moveLogs.
// If you are creating new reorged chain with moved logs, use NewBlockChainReorgMoveLogs instead,
// as blockChain.moveLogs only moves the logs and changing log.BlockHash and log.BlockNumber.
// It also returns 2 slices of block numbers, 1st of which is a slice of blocks from which logs are moved,
// the 2nd of which is a slice of blocks to which logs are moved.
// NOTE: Do not use this function directly, since it only moves logs to new blocks and does not reorg blocks.
// It is meant to be used inside NewBlockChainReorgMoveLogs, and NewBlockChain
func (chain BlockChain) moveLogs(
	movedLogs map[uint64][]MoveLogs,
) (
	[]uint64, // Blocks from which logs are moved from
	[]uint64, // Blocks to which logs are moved to
) {
	// A slice of unique blockNumbers that logs will be moved from.
	// Might be useful to caller, maybe to create empty blocks (no logs) for the old chain.
	moveFromBlocks := make([]uint64, len(movedLogs))
	// A slice of unique blockNumbers that logs will be moved to.
	// Might be useful to caller, maybe to create empty blocks (no logs) for the old chain.
	var moveToBlocks []uint64

	var c int
	for moveFromBlock, moves := range movedLogs {
		b, ok := chain[moveFromBlock]
		if !ok {
			panic(fmt.Sprintf("logs moved from non-existent block %d", moveFromBlock))
		}

		moveFromBlocks[c] = moveFromBlock
		c++

		for _, move := range moves {
			targetBlock, ok := chain[move.NewBlock]
			if !ok {
				panic(fmt.Sprintf("logs moved to non-existent block %d", move.NewBlock))
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

			// Remove logs from b
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

// newBlockChain returns a new blockChain from |mappedLogs|.
// The parameter |reorgedBlock| // is used to determine block.reorgedHere and block.toBeForked.
func newBlockChain(
	mappedLogs map[uint64][]types.Log,
	reorgedBlock uint64,
) BlockChain {
	chain := make(BlockChain)

	var noReorg bool
	if reorgedBlock == NoReorg {
		noReorg = true
	}

	for blockNumber, logs := range mappedLogs {
		var toBeForked bool
		if noReorg {
			toBeForked = false
		} else {
			toBeForked = blockNumber >= reorgedBlock
		}

		chain[blockNumber] = &Block{
			blockNumber: blockNumber,
			hash:        logs[0].BlockHash,
			logs:        logs,
			reorgedHere: blockNumber == reorgedBlock,
			toBeForked:  toBeForked,
		}
	}

	return chain
}

// NewBlockChain is the preferred way to init reorgsim `blockChain`s. It accept a slice of `ReorgEvent` and
// uses each event to construct a reorged chain, which will be appended to the second return variable.
// Each ReorgEvent will result in its own blockChain, with the identical index.
func NewBlockChain(
	logs map[uint64][]types.Log,
	events []ReorgEvent,
) (
	BlockChain, //nolint:revive
	[]BlockChain, //nolint:revive
) {
	if len(events) == 0 {
		return newBlockChain(logs, NoReorg), nil
	}

	chain := newBlockChain(logs, events[0].ReorgBlock)
	reorgedChains := make([]BlockChain, len(events))

	for i, event := range events {
		var prevChain BlockChain
		if i == 0 {
			prevChain = chain
		} else {
			prevChain = reorgedChains[i-1]
		}

		// Reorg and move logs
		forkedChain := copyBlockChain(prevChain)
		forkedChain.reorg(event.ReorgBlock, i)
		moveFromBlocks, moveToBlocks := forkedChain.moveLogs(event.MovedLogs)

		// Ensure that all moveFromBlock exist in forkedChain.
		// If the forkedChain did not have moveToBlock, create one.
		// This created block will need to have non-deterministic blockHash via RandomHash()
		// because the block needs to have different blockHash vs the reorgedBlock's hash (PRandomHash()).
		for _, prevFrom := range moveFromBlocks {
			if _, ok := prevChain[prevFrom]; !ok {
				panic(fmt.Sprintf("moved from non-existent block %d in the old chain", prevFrom))
			}

			if b, ok := forkedChain[prevFrom]; !ok || b == nil {
				fromBlock := &Block{
					blockNumber: prevFrom,
					reorgedHere: prevFrom == event.ReorgBlock,
					toBeForked:  true,
				}

				fromBlock.reorg(i)
				forkedChain[prevFrom] = fromBlock
			}
		}

		// Ensure that all moveToBlocks exist in prevChain.
		// If the old chain did not have moveToBlock, create one.
		// This created block will need to have non-deterministic blockHash via RandomHash()
		// because the block needs to have different blockHash vs the reorgedBlock's hash (PRandomHash()).
		for _, forkedTo := range moveToBlocks {
			if _, ok := forkedChain[forkedTo]; !ok {
				panic(fmt.Sprintf("moved to non-existent block %d in the new chain", forkedTo))
			}

			if _, ok := prevChain[forkedTo]; !ok {
				toBlock := &Block{
					blockNumber: forkedTo,
					reorgedHere: forkedTo == event.ReorgBlock,
					toBeForked:  true,
				}

				toBlock.reorg(i)
				prevChain[forkedTo] = toBlock
			}
		}

		// Make sure the block from which the logs moved
		reorgedChains[i] = forkedChain
	}

	return chain, reorgedChains
}

// newBlockChainReorgSimple returns a tuple of 2 blockChain, the first being the original chain,
// and the second being the reorged chain based solely on 1 |reorgedBlock| (no logs will be moved to new blocks).
// It is unexported and is only here as a convenient function wrapped by others, and for reorg logic tests.
// NOTE: maybe deprecated in favor of NewBlockChain
func newBlockChainReorgSimple(
	mappedLogs map[uint64][]types.Log,
	reorgedBlock uint64,
) (
	BlockChain,
	BlockChain,
) {
	// The "good old chain"
	chain := newBlockChain(mappedLogs, reorgedBlock)

	// No reorg - use the same chain
	if reorgedBlock == NoReorg {
		return chain, chain
	}

	// |reorgedChain| will differ from |oldChain| after |reorgedBlock|
	reorgedChain := copyBlockChain(chain)
	reorgedChain.reorg(reorgedBlock, 0)

	return chain, reorgedChain
}

// NewBlockChainReorgMoveLogs calls newBlockChainReorgSimple, and uses |event.MovedLogs| to move logs.
// If |event.MovedLogs| is nil, the result chains are identical to those of newBlockChainReorgSimple.
// NOTE: maybe deprecated in favor of NewBlockChain
func NewBlockChainReorgMoveLogs(
	mappedLogs map[uint64][]types.Log,
	event ReorgEvent,
) (
	BlockChain, //nolint:revive
	BlockChain, //nolint:revive
) {
	for blockNumber := range event.MovedLogs {
		if blockNumber < event.ReorgBlock {
			panic(fmt.Sprintf("blockNumber %d < reorgedBlock %d", blockNumber, event.ReorgBlock))
		}
	}

	chain, reorgedChain := newBlockChainReorgSimple(mappedLogs, event.ReorgBlock)

	if len(event.MovedLogs) != 0 {
		_, moveToBlocks := reorgedChain.moveLogs(event.MovedLogs)

		// Ensure that all moveToBlocks exist in original chain
		for _, moveToBlock := range moveToBlocks {
			// If the old chain did not have moveToBlock, create one.
			// This created block will need to have non-deterministic blockHash via RandomHash()
			// because the block needs to have different blockHash vs the reorgedBlock's hash (PRandomHash()).
			if _, ok := chain[moveToBlock]; !ok {
				toBlock := &Block{
					blockNumber: moveToBlock,
					reorgedHere: moveToBlock == event.ReorgBlock,
					toBeForked:  true,
				}

				chain[moveToBlock] = toBlock
			}
		}
	}

	return chain, reorgedChain
}
