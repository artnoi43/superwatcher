package superwatcher

// FilterResult is published by emitter
type FilterResult struct {
	FromBlock     uint64
	LastGoodBlock uint64
	GoodBlocks    []*BlockInfo
	ReorgedBlocks []*BlockInfo
}
