package superwatcher

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockInfo represents the minimum block info needed for superwatcher.
type BlockInfo struct {
	Number       uint64
	Hash         common.Hash
	Logs         []*types.Log
	LogsMigrated bool // true means that after chain reorg, interesting logs were entirely moved from this block
}

// String returns the block hash with 0x prepended in all lowercase string.
func (b *BlockInfo) String() string {
	return gslutils.StringerToLowerString(b.Hash)
}

func (b *BlockInfo) BlockNumberString() string {
	return fmt.Sprintf("%d", b.Number)
}
