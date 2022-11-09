package superwatcher

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockInfo represents the bare minimum info needed for superwatcher.
// It is used in WatcherEngine to detect chain reorgs, and is embedded
// in FilterResult.
type BlockInfo struct {
	Number uint64
	Hash   common.Hash
	Logs   []*types.Log
}

// NewBlankBlockInfo returns a new BlockInfo sans the event logs.
// Callers will have to populate the logs themselves.
func NewBlankBlockInfo(
	blockNumber uint64,
	blockHash common.Hash,
) *BlockInfo {
	return &BlockInfo{
		Number: blockNumber,
		Hash:   blockHash,
	}
}

// String returns the block hash with 0x prepended
func (b *BlockInfo) String() string {
	return b.Hash.String()
}

func (b *BlockInfo) BlockNumberString() string {
	return fmt.Sprintf("%d", b.Number)
}
