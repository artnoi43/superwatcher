package poller

import (
	"context"
	"math/big"
	"sync"

	"github.com/artnoi43/gsl/concurrent"
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

// safeMap wraps map types used in mapLogs for concurrent operations
type safeMap[V any] struct {
	sync.RWMutex
	m map[uint64]V
}

// mapLogsResult represents information of fresh blocks mapped by mapLogs.
// It contains fresh data, i.e. not from tracker.
type mapLogsResult struct {
	reorged bool
	hash    common.Hash
	header  superwatcher.BlockHeader
	logs    []*types.Log
}

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
	// doHeader specifies if mapLogs will call |getHeaderFunc| for blocks with logs
	doHeader bool,
	// tracker is used to store known BlockInfo from previous calls
	tracker *blockInfoTracker,
	// getHeaderFunc is used to get block hash for known block with missing logs.
	getHeaderFunc func(context.Context, *big.Int) (superwatcher.BlockHeader, error),
) (
	map[uint64]*mapLogsResult,
	error,
) {
	mapResults := safeMap[*mapLogsResult]{
		m: make(map[uint64]*mapLogsResult),
	}

	// Collect mapFreshLogs and mapFreshHashes. All new logs will be collected,
	// and each log's blockHash will be compare with its neighbors in the same block
	// to see if the hashes differ. If it differs, the code panics, since it should not happen.
	for _, log := range logs {
		number := log.BlockNumber

		// Add new hash to mapResults
		if mapResult, ok := mapResults.m[number]; !ok {
			mapResults.m[number] = &mapLogsResult{
				hash: log.BlockHash,
			}
		} else if ok {
			if mapResult.hash != log.BlockHash {
				return mapResults.m, errors.Wrapf(
					errHashesDiffer, "logs in the same block %d has different hash %s vs %s",
					number, gslutils.StringerToLowerString(mapResult.hash), gslutils.StringerToLowerString(log.BlockHash),
				)
			}
		}

		// Add all new logs to mapFreshLogs
		mapResults.m[number].logs = append(mapResults.m[number].logs, log)
	}

	if doHeader {
		errChan := make(chan error)

		var wg sync.WaitGroup
		// Get block headers concurrently for blocks with logs
		for number, mapResult := range mapResults.m {
			wg.Add(1)

			go func(n uint64, r *mapLogsResult) {
				defer wg.Done()

				// TODO: Make batch requests instead of Goroutines
				header, err := getHeaderFunc(ctx, big.NewInt(int64(n)))
				if err != nil {
					errChan <- err
				}

				func() {
					mapResults.Lock()
					defer mapResults.Unlock()

					r.header = header
				}()

				// If header hash differs from logHash, cause a return with error
				// TODO: Maybe we can better handle this?
				if l := r.logs[0]; l != nil {
					if headerHash := header.Hash(); headerHash != l.BlockHash {
						errChan <- errors.Wrapf(
							errHashesDiffer, "block %d header blockHash %s differ from log's blockHash %s",
							n, gslutils.StringerToLowerString(headerHash), gslutils.StringerToLowerString(l.BlockHash),
						)
					}
				}
			}(number, mapResult)
		}

		// wg.Wait() is called in WaitAndCollectErrors
		err := concurrent.WaitAndCollectErrors(&wg, errChan)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get headers")
		}
	}

	// i.e. if poller.doReorg == false
	if tracker == nil {
		return mapResults.m, nil
	}

	emptyHash := common.Hash{}
	// Compare all known (tracker) block hashes to new ones
	for number := toBlock; number >= fromBlock; number-- {
		// Continue if |number| was never saved to tracker
		trackerBlock, ok := tracker.getTrackerBlockInfo(number)
		if !ok {
			continue
		}

		// Logs may be moved from blockNumber, hence there's value in tracker but nil value in mapResults.m[number]
		mapResult, ok := mapResults.m[number]
		if !ok {
			// Prevent nil pointer dereference
			mapResult = new(mapLogsResult)
			mapResults.m[number] = mapResult
		}

		if mapResult == nil {
			return nil, errors.Wrap(superwatcher.ErrProcessReorg, "nil result after check")
		}

		// If hash is empty, we know for a fact that the logs were entirely moved from this block
		if mapResult.hash == emptyHash {
			// Fetch header again for this block (with missing log) in case doHeader is false
			freshHeader, err := getHeaderFunc(ctx, big.NewInt(int64(number)))
			if err != nil {
				return nil, errors.Wrap(superwatcher.ErrFetchError, err.Error())
			}

			freshHash := freshHeader.Hash()

			mapResult.hash = freshHash
			mapResult.header = freshHeader
			trackerBlock.LogsMigrated = true
		}

		// If we have same blockHash with same logs length we can assume that the block was not reorged.
		if trackerBlock.Hash == mapResult.hash && len(trackerBlock.Logs) == len(mapResult.logs) {
			continue
		}

		for _, trackerLog := range trackerBlock.Logs {
			trackerLog.Removed = true
		}

		mapResult.reorged = true
	}

	return mapResults.m, nil
}
