package poller

import (
	"context"
	"math/big"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

// mapLogsResult represents information of fresh blocks mapped by mapLogs.
// It contains fresh data, i.e. not from tracker.
type mapLogsResult struct {
	reorged bool                     // true if the tracker block hash differs from fresh block hash
	hash    common.Hash              // fresh block hash
	header  superwatcher.BlockHeader // fresh block header
	logs    []*types.Log             // fresh block logs
}

// mapLogs compare |logs| and their blockHashes to known blockHashes in |tracker.
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
	client superwatcher.EthClient,
) (
	map[uint64]*mapLogsResult,
	error,
) {
	mapResults := make(map[uint64]*mapLogsResult)

	// Collect mapFreshLogs and mapFreshHashes. All new logs will be collected,
	// and each log's blockHash will be compare with its neighbors in the same block
	// to see if the hashes differ. If it differs, the code panics, since it should not happen.
	for _, log := range logs {
		number := log.BlockNumber

		// Add new hash to mapResults
		if mapResult, ok := mapResults[number]; !ok {
			mapResults[number] = &mapLogsResult{
				hash: log.BlockHash,
			}
		} else if ok {
			if mapResult.hash != log.BlockHash {
				return mapResults, errors.Wrapf(
					errHashesDiffer, "logs in the same block %d has different hash %s vs %s",
					number, gslutils.StringerToLowerString(mapResult.hash), gslutils.StringerToLowerString(log.BlockHash),
				)
			}
		}

		// Add all new logs to mapFreshLogs
		mapResults[number].logs = append(mapResults[number].logs, log)
	}

	if doHeader {
		// Get block headers using BatchCallContext
		var batchCalls []superwatcher.BatchCallable
		for i := toBlock; i >= fromBlock; i-- {
			if _, ok := mapResults[i]; !ok {
				continue
			}

			batchCalls = append(batchCalls, &headerByNumberBatch{
				Number: i,
			})
		}

		if err := superwatcher.BatchCall(ctx, client, batchCalls); err != nil {
			return nil, errors.Wrap(superwatcher.ErrFetchError, err.Error())
		}

		for _, getHeaderCall := range batchCalls {
			call, ok := getHeaderCall.(*headerByNumberBatch)
			if !ok {
				return nil, errors.New("type assertion failed: getHeaderCall is not *superwatcher.HeaderByNumberBatchCallable")
			}
			mapResult, ok := mapResults[call.Number]
			if !ok {
				return nil, errors.Wrap(superwatcher.ErrProcessReorg, "a header did not have result")
			}

			if headerHash := call.Header.Hash(); headerHash != mapResult.hash {
				return nil, errors.Wrapf(superwatcher.ErrFetchError, "header block hash %s on block %d is different than log block hash %s", headerHash.String(), call.Number, mapResult.hash.String())
			}

			mapResult.header = call.Header
		}
	}

	// i.e. if poller.doReorg == false
	if tracker == nil {
		return mapResults, nil
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
		mapResult, ok := mapResults[number]
		if !ok {
			// Prevent nil pointer dereference
			mapResult = new(mapLogsResult)
			mapResults[number] = mapResult
		}

		if mapResult == nil {
			return nil, errors.Wrap(superwatcher.ErrProcessReorg, "nil result after check")
		}

		// If hash is empty, we know for a fact that the logs were entirely moved from this block
		if mapResult.hash == emptyHash {
			// Fetch header again for this block (with missing log) in case doHeader is false
			freshHeader, err := client.HeaderByNumber(ctx, big.NewInt(int64(number)))
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

	return mapResults, nil
}
