package reorgsim

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Block represents the Ethereum Block.
// It is also used as superwatcher.BlockHeader.
type Block struct {
	blockNumber uint64
	hash        common.Hash
	logs        []types.Log

	reorgedHere bool // reorgedHere marks if this block is where an ReorgEvent begins
	toBeForked  bool // toBeForked marks if this block will later be forked from the old chain according to ReorgEvent
}

func (b *Block) Number() uint64 {
	return b.blockNumber
}

// Implements superwatcher.BlockHeader
// We'll use block in place of *types.Header,
// because *types.Header is too packed to mock.
func (b *Block) Hash() common.Hash {
	return b.hash
}

func (b *Block) Logs() []types.Log {
	return b.logs
}

// Nonce mocks field *types.Header.Nonce
func (b *Block) Nonce() types.BlockNonce {
	return types.EncodeNonce(b.blockNumber)
}

// Time mocks field *types.Header.Time
func (b *Block) Time() uint64 {
	return b.blockNumber
}

// GasLimit mocks field *types.Header.GasLimit
func (b *Block) GasLimit() uint64 {
	return b.blockNumber
}

// GasUsed mocks field *types.Header.GasUsed
func (b *Block) GasUsed() uint64 {
	return b.blockNumber
}

// reorg takes a block, and simulates chain reorg on that block
// by changing the hash, and changing the logs' block hashes.
// math.RandInt(seed) is mixed with b.blockNumber to produce different
// block hash for the same block across different chains created by []ReorgEvent.
func (b *Block) reorg(reorgIndex int) *Block {
	reorgedHash := ReorgHash(b.blockNumber, reorgIndex)

	logs := make([]types.Log, len(b.logs))
	copy(logs, b.logs)

	// Use index to access logs so that the internal array members change value too.
	for i := range logs {
		logs[i].BlockHash = reorgedHash
	}

	return &Block{
		blockNumber: b.blockNumber,
		hash:        reorgedHash,
		logs:        logs,
		reorgedHere: b.reorgedHere,
		toBeForked:  true,
	}
}

func (b *Block) removeLogs(txHashes []common.Hash) {
	if len(b.logs) == 0 {
		panic(fmt.Sprintf("block %d has no logs", b.blockNumber))
	}

	// Only keep log whose TxHash is not in |txHashes|.
	remaining := gslutils.FilterSlice(b.logs, func(log types.Log) bool {
		return !gslutils.Contains(txHashes, log.TxHash)
	})

	b.logs = remaining
}
