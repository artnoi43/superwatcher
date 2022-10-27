package emitter

import "github.com/artnoi43/superwatcher/lib"

// FilterResult is published by emitter
type FilterResult struct {
	FromBlock     uint64
	LastGoodBlock uint64
	GoodBlocks    []*lib.BlockInfo
	ReorgedBlocks []*lib.BlockInfo
}
