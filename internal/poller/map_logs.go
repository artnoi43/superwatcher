package poller

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
)

// mapLogs compare |logs| and their blockHashes to known blockHashes in |tracker.
// If a block has different tracker blockHash, the blockNumber will be marked true in |mapRemovedBlocks|.
// The emitter can then use the information in |mapRemovedBlocks| to get the removed `superwatcher.BlockInfo`
// and publish the removed BlockInfo in the `superwatcher.FilterResult.ReorgedBlocks`
func mapLogs(
	fromBlock uint64,
	toBlock uint64,
	logs []*types.Log, // Use pointers to avoid expensive copy of []types.Log
	tracker *blockInfoTracker,
) (
	map[uint64]bool, // Maps blockNumber to reorged status
	map[uint64]common.Hash, // Maps blockNumber to blockHash
	map[uint64][]*types.Log, // Maps blockNumber to logs
) {
	// We can actually just return |mapRemovedBlocks|, but we will have to iterate the logs anyway,
	// so we collect the data into maps here to save costs later.
	mapRemovedBlocks := make(map[uint64]bool)
	mapFreshHashes := make(map[uint64]common.Hash)
	mapFreshLogs := make(map[uint64][]*types.Log)

	// Collect mapFreshLogs and mapFreshHashes. All new logs will be collected,
	// and each log's blockHash will be compare with its neighbors in the same block
	// to see if the hashes differ. If it differs, the code panics, since it should not happen.
	for _, log := range logs {
		number := log.BlockNumber

		// Add new blockNumber hash
		if h, ok := mapFreshHashes[number]; !ok {
			mapFreshHashes[number] = log.BlockHash
		} else {
			if h != log.BlockHash {
				logger.Panic(
					"fresh logs on same block has different blockHash",
					zap.Uint64("blockNumber", number),
					zap.String("blockHash0", h.String()),
					zap.String("blockHash1", log.BlockHash.String()),
				)
			}
		}

		// Add all new logs to mapFreshLogs
		mapFreshLogs[number] = append(mapFreshLogs[number], log)
	}

	// i.e. poller.DoReorg == false
	if tracker == nil {
		return mapRemovedBlocks, mapFreshHashes, mapFreshLogs
	}

	// Compare all known (tracker) block hashes to new ones
	for number := toBlock; number >= fromBlock; number-- {
		// Continue if log.BlockNumber was never saved to tracker
		trackerBlock, found := tracker.getTrackerBlockInfo(number)
		if !found {
			continue
		}

		// If we have same blockHash with same logs length,
		// we can assume that the block was not reorged.
		if trackerBlock.Hash == mapFreshHashes[number] {
			if len(trackerBlock.Logs) == len(mapFreshLogs[number]) {
				continue
			}
		}

		mapRemovedBlocks[number] = true
		for _, trackerLog := range trackerBlock.Logs {
			trackerLog.Removed = true
		}
	}

	return mapRemovedBlocks, mapFreshHashes, mapFreshLogs
}
