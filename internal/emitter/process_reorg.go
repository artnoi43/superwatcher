package emitter

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
)

// mapFreshLogs collects and maps information from the logs filtered into 3 hashmaps, with blockNumber as key.
func mapFreshLogs(
	freshLogs []types.Log,
) (
	mapFreshHashes map[uint64]common.Hash,
	mapFreshLogs map[uint64][]*types.Log,
	mapProcessLogs map[uint64][]*types.Log,
) {
	// freshHashesByBlockNumber maps blockNumber to fresh hash
	mapFreshHashes = make(map[uint64]common.Hash)
	// freshLogsByBlockNumber maps blockNumber to fresh logs
	mapFreshLogs = make(map[uint64][]*types.Log)
	// processLogsByBlockNumber maps blockNumber to all logs to be processed.
	// Since this function does not use data from tracker, processLogsByBlockNumber is not fully populated here.
	// it should be later passed to
	mapProcessLogs = make(map[uint64][]*types.Log)

	for i := range freshLogs {
		freshLog := freshLogs[i]
		freshLogBlockNumber := freshLog.BlockNumber
		freshLogBlockHash := freshLog.BlockHash

		// Check if we saw the block hash
		if _, ok := mapFreshHashes[freshLogBlockNumber]; !ok {
			// If the blockNumber is not found, it means we've never seen this block hash before.
			mapFreshHashes[freshLogBlockNumber] = freshLogBlockHash

		} else {
			if mapFreshHashes[freshLogBlockNumber] != freshLogBlockHash {
				// If we saw this log's block hash before from the fresh eventLogs
				// from this loop's previous loop iteration(s), but the hashes are different:
				// Fatal, should not happen
				logger.Info("fresh blockHash differs from tracker blockHash",
					zap.String("tag", "filterLogs bug"),
					zap.Uint64("freshBlockNumber", freshLogBlockNumber),
					zap.Any("known tracker blockHash", mapFreshHashes[freshLogBlockNumber]),
					zap.Any("fresh log blockHash", freshLogBlockHash))
			}
		}

		// Collect the (fresh) log into mapFreshLogs and mapProcessLogs
		log := &freshLogs[i]
		mapFreshLogs[freshLogBlockNumber] = append(mapFreshLogs[freshLogBlockNumber], log)
		mapProcessLogs[freshLogBlockNumber] = append(mapProcessLogs[freshLogBlockNumber], log)
	}

	return mapFreshHashes, mapFreshLogs, mapProcessLogs
}

/*
	[PopulateProcessLogs]: Check if block hash saved in tracker matches the fresh block hash.
	If they are different, old logs from w.tracker will be tagged as Removed and
	PREPENDED in processLogs[blockNumber]

	Let's say we have these logs in the tracker:

	{block:68, hash:"0x68"}, {block: 69, hash:"0x69"}, {block:70, hash:"0x70"}

	And then we have these fresh logs:

	{block:68, hash:"0x68"}, {block: 69, hash:"0x112"}, {block:70, hash:"0x70"}

	The result processLogs will look like this map:
	{
		68: [{block:68, hash:"0x68"}]
		69: [{block: 69, hash:"0x69", removed: true}, {block: 69, hash:"0x112"}]
		70: [{block:70, hash:"0x70"}]
	}
*/

// processReorg compares fresh hashes and hashes saved in tracker, and prepends reorged blocks to mapProcessLogs.
// processReorg returns a map of blockNumber and a bool indicating if chain reorg was detected on that blockNumber.
// Note: Go maps are passed by reference, so there's no need to return the map.
func processReorg(
	tracker *blockInfoTracker,
	fromBlock, toBlock uint64,
	// New hashes from *ethclient.Client.HeaderByNumber
	mapFreshHashes map[uint64]common.Hash,
	// New logs from *ethclient.Client.FilterLogs
	mapFreshLogs map[uint64][]*types.Log,
	// mapProcessLogs initially contains only fresh logs from mapFreshlogs,
	// but this function will modify (concat) it with tracker logs if it detects a reorg.
	mapProcessLogs map[uint64][]*types.Log,
) (
	map[uint64]bool,
	error,
) {
	// This map will be returned to caller.
	reorgedChain := make(map[uint64]bool)

	for blockNumber := toBlock; blockNumber >= fromBlock; blockNumber-- {
		// If the block had not been saved into w.tracker (new blocks), it's probably fresh blocks,
		// which are not yet 'reorged' at the execution time.
		trackerBlock, foundInTracker := tracker.getTrackerBlockInfo(blockNumber)
		if !foundInTracker {
			continue
		}

		// If tracker's is the same from recently filtered hash.
		// i.e. No reorg
		if h := mapFreshHashes[blockNumber]; h == trackerBlock.Hash {
			// If number of logs did not match - we're really screwed.
			if len(mapFreshLogs[blockNumber]) == len(trackerBlock.Logs) {
				continue
			}

			// TODO: Should we panic or return an error after POC passed?
			logger.Panic(
				"tracker has different number of logs for identical blockHash",
				zap.String("blockHash", trackerBlock.String()),
			)
		}

		// REORG HAPPENED!
		reorgedChain[blockNumber] = true
		// Mark every log in this block as removed
		for _, oldLog := range trackerBlock.Logs {
			oldLog.Removed = true
		}

		// Concat logs from the same block, old logs first, into freshLogs
		mapProcessLogs[blockNumber] = append(trackerBlock.Logs, mapProcessLogs[blockNumber]...)
	}

	return reorgedChain, nil
}
