package superwatcher

// FilterResult is published by emitter
type FilterResult struct {
	FromBlock     uint64       // The emitter's `fromBlock`
	ToBlock       uint64       // The emitter's `toBlock`
	LastGoodBlock uint64       // This number should be saved to StateDataGateway with SetLastRecordedBlock for the emitter
	GoodBlocks    []*BlockInfo // Can be either (1) fresh, new blocks, or (2) blocks whose hashes had not changed yet.
	ReorgedBlocks []*BlockInfo // Blocks that poller marked as removed. A service should undo/revert its actions done on the blocks.
}

// LastGoodBlock computes `FilterResult.LastGoodBlock` based on |result|.
func LastGoodBlock(
	result *FilterResult,
) uint64 {
	if len(result.ReorgedBlocks) != 0 {
		// If there's also goodBlocks during reorg
		if l := len(result.GoodBlocks); l != 0 {
			// Use last good block's number as LastGoodBlock
			lastGood := result.GoodBlocks[l-1].Number
			firstReorg := result.ReorgedBlocks[0].Number

			// lastGood should be less than firstReorg
			if lastGood > firstReorg {
				lastGood = firstReorg - 1
			}

			return lastGood
		}

		// If there's no goodBlocks, then we should re-filter the whole range
		return result.FromBlock - 1
	}

	return result.ToBlock
}
