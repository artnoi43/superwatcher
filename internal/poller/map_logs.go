package poller

import (
	"context"
	"math/big"
	"sync"

	"github.com/artnoi43/gsl/concurrent"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

// mapLogs compare |logs| and their blockHashes to known blockHashes in |tracker.
// If a block has different tracker blockHash, the blockNumber will be marked true in |mapRemovedBlocks|.
// The emitter can then use the information in |mapRemovedBlocks| to get the removed `superwatcher.BlockInfo`
// and publish the removed BlockInfo in the `superwatcher.PollResult.ReorgedBlocks`.
// We avoid getting block headers for all the range fromBlock-toBlock, so we use block hashes from the logs.
// In case the known logs were removed from a particular block, then we won't have the fresh hash for that block.
// In that case, we'll use |getHeaderFunc| to get the header for that particular block.
func mapLogs(
	// For calling |getHeaderFunc|
	ctx context.Context,
	// fromBlock from poller.Poll
	fromBlock uint64,
	// toBlock from poller.Poll
	toBlock uint64,
	// logs are slice of pointers to avoid expensive copies.
	logs []*types.Log,
	// tracker is used to store known BlockInfo from previous calls
	tracker *blockInfoTracker,
	// getHeaderFunc is used to get block hash for known block with missing logs.
	getHeaderFunc func(context.Context, *big.Int) (superwatcher.BlockHeader, error),
) (
	map[uint64]bool, // Maps blockNumber to reorged status
	map[uint64]superwatcher.BlockHeader, // Maps blockNumber to blockHeader
	map[uint64]common.Hash, // Maps blockNumber to blockHash
	map[uint64][]*types.Log, // Maps blockNumber to logs
	error,
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
		} else if ok {
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

	// Get block headers from blocks in mapFreshLogs concurrently
	var safeMapFreshHeaders struct {
		sync.RWMutex
		m map[uint64]superwatcher.BlockHeader
	}
	if getHeaderFunc != nil {
		safeMapFreshHeaders.m = make(map[uint64]superwatcher.BlockHeader)
		errChan := make(chan error)

		var wg sync.WaitGroup
		for number := range mapFreshLogs {
			wg.Add(1)
			go func(n uint64) {
				defer wg.Done()

				h, err := getHeaderFunc(ctx, big.NewInt(int64(n)))
				if err != nil {
					errChan <- err
				}

				safeMapFreshHeaders.Lock()
				defer safeMapFreshHeaders.Unlock()

				safeMapFreshHeaders.m[n] = h
			}(number)
		}

		// wg.Wait() is called in WaitAndCollectErrors
		err := concurrent.WaitAndCollectErrors(&wg, errChan)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "failed to get headers")
		}
	}

	var mapFreshHeaders map[uint64]superwatcher.BlockHeader
	if safeMapFreshHeaders.m != nil {
		mapFreshHeaders = safeMapFreshHeaders.m
	}

	// i.e. poller.DoReorg == false
	if tracker == nil {
		return mapRemovedBlocks, mapFreshHeaders, mapFreshHashes, mapFreshLogs, nil
	}

	// Compare all known (tracker) block hashes to new ones
	for number := toBlock; number >= fromBlock; number-- {
		// Continue if log.BlockNumber was never saved to tracker
		trackerBlock, found := tracker.getTrackerBlockInfo(number)
		if !found {
			continue
		}

		freshHash, ok := mapFreshHashes[number]
		// Logs may be moved from blockNumber, hence there's no value in mapFreshHashes
		if !ok {
			var freshHeader superwatcher.BlockHeader
			if mapFreshHeaders != nil {
				freshHeader = mapFreshHeaders[number]
			}

			// Fetch header again for this block (with missing log) in case doHeader is false
			if freshHeader == nil {
				var err error
				freshHeader, err = getHeaderFunc(ctx, big.NewInt(int64(number)))
				if err != nil {
					return nil, nil, nil, nil, errors.Wrap(superwatcher.ErrFetchError, err.Error())
				}
			}

			freshHash = freshHeader.Hash()
			mapFreshHashes[number] = freshHash
			mapFreshLogs[number] = nil
			trackerBlock.LogsMigrated = true
		}

		// If we have same blockHash with same logs length we can assume that the block was not reorged.
		if trackerBlock.Hash == freshHash && len(trackerBlock.Logs) == len(mapFreshLogs[number]) {
			continue
		}

		for _, trackerLog := range trackerBlock.Logs {
			trackerLog.Removed = true
		}

		mapRemovedBlocks[number] = true
	}

	return mapRemovedBlocks, mapFreshHeaders, mapFreshHashes, mapFreshLogs, nil
}
