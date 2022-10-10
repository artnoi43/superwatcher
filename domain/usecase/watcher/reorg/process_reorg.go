package reorg

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
)

// PopulateInitialMaps collects **fresh** hashes and logs into 3 maps
func PopulateInitialMaps(
	freshLogs []types.Log,
	freshHeaders map[uint64]*types.Header,
) (
	freshHashesByBlockNumber map[uint64]common.Hash,
	freshLogsByBlockNumber map[uint64][]*types.Log,
	processLogsByBlockNumber map[uint64][]*types.Log,
) {
	// freshHashesByBlockNumber maps blockNumber to fresh hash
	freshHashesByBlockNumber = make(map[uint64]common.Hash)
	// freshLogsByBlockNumber maps blockNumber to fresh logs
	freshLogsByBlockNumber = make(map[uint64][]*types.Log)
	// processLogsByBlockNumber maps blockNumber to all logs to be processed.
	// Since this function does not use data from tracker, processLogsByBlockNumber is not fully populated here.
	// it should be later passed to
	processLogsByBlockNumber = make(map[uint64][]*types.Log)

	for i := range freshLogs {
		freshLog := freshLogs[i]
		freshLogBlockNumber := freshLog.BlockNumber
		freshLogBlockHash := freshLog.BlockHash
		freshBlockHash := freshHeaders[freshLogBlockNumber].Hash()

		if !bytes.Equal(freshBlockHash[:], freshLogBlockHash[:]) {
			// @TODO: How to handle this?
			logger.Panic(
				"freshBlockHash and freshBlockHashFromLog differ",
				zap.String("freshBlockHash", freshBlockHash.String()),
				zap.String("freshBlockHashFromLog", freshLogBlockHash.String()),
			)
		}

		// Check if we saw the block hash
		if _, ok := freshHashesByBlockNumber[freshLogBlockNumber]; !ok {
			// We've never seen this block hash before
			freshHashesByBlockNumber[freshLogBlockNumber] = freshBlockHash
		} else {
			if freshHashesByBlockNumber[freshLogBlockNumber] != freshBlockHash {
				// If we saw this log's block hash before from the fresh eventLogs
				// from this loop's previous loop iteration(s), but the hashes are different:
				// Fatal, should not happen
				logger.Panic("fresh blockHash differs from tracker blockHash",
					zap.String("tag", "filterLogs bug"),
					zap.Uint64("freshBlockNumber", freshLogBlockNumber),
					zap.Any("known tracker blockHash", freshHashesByBlockNumber[freshLogBlockNumber]),
					zap.Any("fresh blockHash", freshBlockHash))
			}
		}

		// Collect this log fresh into freshLogs and processLogs
		thisLog := &freshLogs[i]
		freshLogsByBlockNumber[freshLogBlockNumber] = append(freshLogsByBlockNumber[freshLogBlockNumber], thisLog)
		processLogsByBlockNumber[freshLogBlockNumber] = append(processLogsByBlockNumber[freshLogBlockNumber], thisLog)
	}

	return freshHashesByBlockNumber, freshLogsByBlockNumber, processLogsByBlockNumber
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

// ProcessReorged checks hash fresh-tracker equality, and appends reorged blocks to processLogsByBlockNumber.
// Note: Go maps are passed by reference, so there's no need to return the map.
func ProcessReorged(
	tracker *Tracker,
	fromBlock, toBlock uint64,
	freshHashesByBlockNumber map[uint64]common.Hash,
	freshLogsByBlockNumber map[uint64][]*types.Log,
	processLogsByBlockNumber map[uint64][]*types.Log,
) map[uint64]bool {
	wasReorged := make(map[uint64]bool)
	// Detect and stamp removed/reverted event logs using tracker
	for blockNumber := fromBlock; blockNumber <= toBlock; blockNumber++ {
		// If the block had not been saved into w.tracker (new blocks), it's probably fresh blocks,
		// which are not yet 'reorged' at the execution time.
		trackerBlock, foundInTracker := tracker.GetTrackerBlockInfo(blockNumber)
		if !foundInTracker {
			continue
		}

		// If tracker's is the same from recently filtered hash, i.e. no reorg
		// logger.Info("found block in tracker, comparing hashes in tracker", zap.Uint64("blockNumber", blockNumber))
		if h := freshHashesByBlockNumber[blockNumber]; h == trackerBlock.Hash {
			// Mark blockNumber with identical hash (no reorg)
			if len(freshLogsByBlockNumber[blockNumber]) == len(trackerBlock.Logs) {
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
		processLogsByBlockNumber[blockNumber] = append(trackerBlock.Logs, processLogsByBlockNumber[blockNumber]...)
	}

	return wasReorged
}
