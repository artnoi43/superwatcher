package reorgsim

import (
	"github.com/artnoi43/gsl/gslutils"
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
func (b *block) Hash() common.Hash {
	return b.hash
}

func (b *block) Logs() []types.Log {
	return b.logs
}

// reorg takes a block, and simulates chain reorg on that block
// by changing the hash, and changing the logs' block hashes.
func (b *block) reorg() *block {
	// TODO: implement
	newBlockHash := PRandomHash(b.blockNumber)

	logs := make([]types.Log, len(b.logs))
	copy(logs, b.logs)

	// Use index to access logs so that the internal array members change value too.
	for i := range logs {
		logs[i].BlockHash = newBlockHash
		// logs[i].TxHash = PRandomHash(logs[i].TxHash.Big().Uint64() + 696969)
	}

	return &block{
		blockNumber: b.blockNumber,
		hash:        newBlockHash,
		logs:        logs,
		reorgedHere: b.reorgedHere,
		toBeForked:  true,
	}
}

func (b *block) removeLogs(txHashes []common.Hash) {
	for i, log := range b.logs {
		if gslutils.Contains(txHashes, log.TxHash) {
			b.removeLog(i)
		}
	}
}

func (b *block) removeLog(logIdx int) {
	b.logs = append(b.logs[:logIdx], b.logs[logIdx+1:]...)
}
