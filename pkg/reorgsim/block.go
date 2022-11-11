package reorgsim

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// block represents the Ethereum block. It is also used
// by reorgSim as superwatcher.EmitterBlockHeader.
type block struct {
	blockNumber uint64
	hash        common.Hash
	logs        []types.Log
	reorgedHere bool
	toBeForked  bool
}

// Implements superwatcher.EmitterBlockHeader
// We'll use block in place of *types.Header,
// because *types.Header is too packed to mock.
func (b block) Hash() common.Hash {
	return b.hash
}

// reorg takes a block, and simulates chain reorg on that block
// by changing the hash, and changing the logs' block hashes.
func (b *block) reorg() block {
	// TODO: implement
	newBlockHash := randomHash(b.blockNumber)
	newTxHash := randomHash(b.blockNumber + 696969)
	var logs []types.Log
	copy(logs, b.logs)

	for _, log := range logs {
		log.BlockHash = newBlockHash
		log.TxHash = newTxHash
	}

	return block{
		blockNumber: b.blockNumber,
		hash:        newBlockHash,
		logs:        logs,
		reorgedHere: b.reorgedHere,
		toBeForked:  true,
	}
}
