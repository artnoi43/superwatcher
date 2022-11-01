package superwatcher

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockInfo is saved to Tracker and is used
// to store block information for determining if chain reorg occured.
type BlockInfo struct {
	Number uint64
	Hash   common.Hash
	Logs   []*types.Log // will be removed?
}

func NewBlockInfo(
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
