package superwatcher

// FilterResult is published by emitter
type FilterResult struct {
	FromBlock     uint64       // Emitter's `fromBlock`
	ToBlock       uint64       // Emitter's `toBlock`
	LastGoodBlock uint64       // Block number of the last GoodBlocks
	GoodBlocks    []*BlockInfo // Can be either (1) fresh blocks, or (2) blocks whose hashes had not changed yet.
	ReorgedBlocks []*BlockInfo // Blocks that emitter marked as removed. A service should undo/revert its action done on the blocks.
}
