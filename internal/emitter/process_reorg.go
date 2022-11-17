package emitter

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

// mapFreshLogsByHashes collects **fresh** hashes and logs into 3 maps
func mapFreshLogsByHashes(
	freshLogs []types.Log,
	freshHeaders map[uint64]superwatcher.BlockHeader,
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
		freshBlockHash := freshHeaders[freshLogBlockNumber].Hash()

		// If the fresh block hash from client.HeaderByNumber differs from client.FilterLogs
		// Fatal, should not happen
		if !bytes.Equal(freshBlockHash[:], freshLogBlockHash[:]) {
			// TODO: How to handle this?
			logger.Panic(
				"freshBlockHash and freshBlockHashFromLog differ",
				zap.String("freshBlockHash", freshBlockHash.String()),
				zap.String("freshBlockHashFromLog", freshLogBlockHash.String()),
			)
		}

		// Check if we saw the block hash
		if _, ok := mapFreshHashes[freshLogBlockNumber]; !ok {
			// If the blockNumber is not found, it means we've never seen this block hash before.
			mapFreshHashes[freshLogBlockNumber] = freshBlockHash

		} else {
			if mapFreshHashes[freshLogBlockNumber] != freshBlockHash {
				// If we saw this log's block hash before from the fresh eventLogs
				// from this loop's previous loop iteration(s), but the hashes are different:
				// Fatal, should not happen
				logger.Panic("fresh blockHash differs from tracker blockHash",
					zap.String("tag", "filterLogs bug"),
					zap.Uint64("freshBlockNumber", freshLogBlockNumber),
					zap.Any("known tracker blockHash", mapFreshHashes[freshLogBlockNumber]),
					zap.Any("fresh blockHash", freshBlockHash))
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

// processReorged compares fresh hashes and hashes saved in tracker, and appends reorged blocks to processLogsByBlockNumber.
// Note: Go maps are passed by reference, so there's no need to return the map.
func processReorged(
	tracker *blockInfoTracker,
	fromBlock, toBlock uint64,
	mapFreshHashes map[uint64]common.Hash, // New hashes from *ethclient.Client.HeaderByNumber
	mapFreshLogs map[uint64][]*types.Log, // New logs from *ethclient.Client.FilterLogs
	mapProcessLogs map[uint64][]*types.Log, // Concatenated logs from both old and reorged chains
) map[uint64]bool {
	// This map will be returned to caller. True means the block was reorged and had different hashes.
	wasReorged := make(map[uint64]bool)

	for blockNumber := fromBlock; blockNumber <= toBlock; blockNumber++ {
		// If the block had not been saved into w.tracker (new blocks), it's probably fresh blocks,
		// which are not yet 'reorged' at the execution time.
		trackerBlock, foundInTracker := tracker.getTrackerBlockInfo(blockNumber)
		if !foundInTracker {
			continue
		}

		// If tracker's is the same from recently filtered hash, i.e. no reorg
		// logger.Info("found block in tracker, comparing hashes in tracker", zap.Uint64("blockNumber", blockNumber))
		if h := mapFreshHashes[blockNumber]; h == trackerBlock.Hash {
			// Mark blockNumber with identical hash (no reorg)
			if len(mapFreshLogs[blockNumber]) == len(trackerBlock.Logs) {
				continue
			}
		}

		// REORG HAPPENED!
		wasReorged[blockNumber] = true
		// Mark every log in this block as removed
		for _, oldLog := range trackerBlock.Logs {
			oldLog.Removed = true
		}
		// Concat logs from the same block, old logs first, into freshLogs
		mapProcessLogs[blockNumber] = append(trackerBlock.Logs, mapProcessLogs[blockNumber]...)
	}

	return wasReorged
}
