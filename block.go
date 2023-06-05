package superwatcher

import (
	"github.com/soyart/gsl"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Block represents the minimum block info needed for superwatcher.
// Block data can be retrieved from Block itself or its Header field.
type Block struct {
	// LogsMigrated indicates whether all interesting logs were moved/migrated
	// _from_ this block after a chain reorg or not. The field is primarily used
	// by EmitterPoller to trigger the poller to get new, fresh block hash for a block.
	// The field should always be false if the Block is in PollerResult.GoodBlocks.
	LogsMigrated bool `json:"logsMigrated"`

	Number uint64       `json:"number"`
	Hash   common.Hash  `json:"hash"`
	Header BlockHeader  `json:"-"`
	Logs   []*types.Log `json:"logs"`
}

// String returns the block hash with 0x prepended in all lowercase string.
func (b *Block) String() string {
	return gsl.StringerToLowerString(b.Hash)
}

// BlockHeader is implemented by `blockHeaderWrapper` and `*reorgsim.Block`.
// It is used in place of *types.Header to make writing tests with reorgsim easier.
// More methods may be added as our needs for data from the headers grow,
// or we (i.e. you) can mock the actual *types.Header in reorgsim instead :)
type BlockHeader interface {
	Number() uint64
	Hash() common.Hash
	Nonce() types.BlockNonce
	Time() uint64
	GasLimit() uint64
	GasUsed() uint64
}

// BlockHeaderWrappers wrap *types.Header to implenent BlockHeader
type BlockHeaderWrapper struct {
	Header *types.Header
}

func (h BlockHeaderWrapper) Number() uint64 {
	return h.Header.Number.Uint64()
}

func (h BlockHeaderWrapper) Hash() common.Hash {
	return h.Header.Hash()
}

func (h BlockHeaderWrapper) Nonce() types.BlockNonce {
	return h.Header.Nonce
}

func (h BlockHeaderWrapper) Time() uint64 {
	return h.Header.Time
}

func (h BlockHeaderWrapper) GasLimit() uint64 {
	return h.Header.GasLimit
}

func (h BlockHeaderWrapper) GasUsed() uint64 {
	return h.Header.GasUsed
}
