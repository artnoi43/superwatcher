package emitter

import "github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"

// FilterResult is published by emitter
type FilterResult struct {
	FromBlock     uint64
	LastGoodBlock uint64
	GoodBlocks    []*reorg.BlockInfo
	ReorgedBlocks []*reorg.BlockInfo
}
