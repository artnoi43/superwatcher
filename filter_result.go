package superwatcher

// FilterResult is published by emitter
type FilterResult struct {
	FromBlock     uint64       // Emitter's `fromBlock`
	LastGoodBlock uint64       // Block number of the last GoodBlocks
	GoodBlocks    []*BlockInfo // Can be either (1) fresh blocks (2) blocks whose hashes had not changed yet.
	ReorgedBlocks []*BlockInfo // Blocks that emitter noticed hash difference
}
