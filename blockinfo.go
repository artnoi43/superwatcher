package superwatcher

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockInfo represents the minimum block info needed for superwatcher.
type BlockInfo struct {
	// LogsMigrated indicates whether all interesting logs were moved/migrated
	// _from_ this block after a chain reorg or not. The field is primarily used
	// by EmitterPoller to trigger the poller to get new, fresh block Hash for a block.
	// The field should always be false if the BlockInfo is in PollResult.GoodBlocks.
	LogsMigrated bool

	Number uint64
	Hash   common.Hash
	Logs   []*types.Log
}

// String returns the block hash with 0x prepended in all lowercase string.
func (b *BlockInfo) String() string {
	return gslutils.StringerToLowerString(b.Hash)
}

func (b *BlockInfo) BlockNumberString() string {
	return fmt.Sprintf("%d", b.Number)
}
