package superwatcher

// FilterResult is published by emitter
type FilterResult struct {
	FromBlock     uint64       // The emitter's `fromBlock`
	ToBlock       uint64       // The emitter's `toBlock`
	LastGoodBlock uint64       // This number should be used as LastRecordedBlock for the emitter
	GoodBlocks    []*BlockInfo // Can be either (1) fresh, new blocks, or (2) blocks whose hashes had not changed yet.
	ReorgedBlocks []*BlockInfo // Blocks that emitter marked as removed. A service should undo/revert its actions done on the blocks.
}
