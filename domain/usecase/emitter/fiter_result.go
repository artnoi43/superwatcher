package emitter

import "github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"

// FilterResult is published by emitter
type FilterResult struct {
	GoodBlocks    []*reorg.BlockInfo
	ReorgedBlocks []*reorg.BlockInfo
}
