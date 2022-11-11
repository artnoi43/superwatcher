package reorgsim

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type block struct {
	blockNumber uint64
	hash        common.Hash
	logs        []types.Log
	reorgedHere bool
	forked      bool
}

type blockChain map[uint64]block

func newBlockChain(mappedLogs map[uint64][]types.Log, reorgedAt uint64) (blockChain, blockChain) {
	var found bool
	for blockNumber := range mappedLogs {
		if blockNumber == reorgedAt {
			found = true
		}
	}

	if !found {
		panic("reorgedAt block not found in any logs")
	}

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

	reorgedChain := make(blockChain)
	for blockNumber, oldBlock := range chain {
		if blockNumber >= reorgedAt {
			b := reorgBlock(oldBlock)
			reorgedChain[blockNumber] = b
		}
	}

	return chain, reorgedChain
}

// reorgBlock takes a block, and simulate chain reorg on that block
// by changing the hash, and changing the logs' block hashes
func reorgBlock(b block) block {
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
		forked:      true,
	}
}
