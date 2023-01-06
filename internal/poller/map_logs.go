package poller

import (
	"context"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

// mapLogsResult represents information of fresh blocks mapped by mapLogs.
// It contains fresh data, i.e. not from tracker.
type mapLogsResult struct {
	reorged bool // true if the tracker block hash differs from fresh block hash
	superwatcher.Block
}

// mapLogs maps logs to mapLogsResult, and compares |logs| and their blockHashes to
// known blockHashes in |tracker|. If |tracker| is nil, mapLogs won't compare hashes.
// If |doHeader| is true, mapLogs also gets all block headers for blocks with logs.
// In case the known logs were removed from a particular block, then we won't have
// the fresh hash/header for that block to compare to tracker values.
// In that case, we'll use |client| to call getHeadersByNumbers to get the header for
// blocks with missing logs, regardless of |doHeader| value, so that we have the new block hash.
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
	// tracker is used to store known Block from previous calls
	tracker *blockTracker,
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
				Block: superwatcher.Block{
					Number: number,
					Hash:   log.BlockHash,
				},
			}
		} else if ok {
			if mapResult.Block.Hash != log.BlockHash {
				return mapResults, errors.Wrapf(
					errHashesDiffer, "logs in the same block %d has different hash %s vs %s",
					number, gslutils.StringerToLowerString(mapResult.Block.Hash), gslutils.StringerToLowerString(log.BlockHash),
				)
			}
		}

		// Add all new logs to mapFreshLogs
		mapResults[number].Block.Logs = append(mapResults[number].Block.Logs, log)
	}

	// Collect blocks with missing logs (we saw them with logs in previous calls to mapLogs)
	var blocksMissingLogs []uint64
	if tracker != nil {
		for n := toBlock; n >= fromBlock; n-- {
			if _, ok := tracker.getTrackerBlock(n); !ok {
				continue
			}
			if _, ok := mapResults[n]; ok {
				continue
			}

			blocksMissingLogs = append(blocksMissingLogs, n)
		}
	}

	var headers map[uint64]superwatcher.BlockHeader
	var blocksWithResults []uint64

	// Get block headers using BatchCallContext via getHeadersByNumbers
	if doHeader {
		for n := toBlock; n >= fromBlock; n-- {
			if _, ok := mapResults[n]; !ok {
				continue
			}

			blocksWithResults = append(blocksWithResults, n)
		}

		var err error
		headers, err = getHeadersByNumbers(ctx, client, append(blocksWithResults, blocksMissingLogs...))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get block headers for new logs filtered")
		}

		// Collect results of batchCalls into mapResults
		for n, header := range headers {
			mapResult, ok := mapResults[n]
			// Should be one from blocksMissingLogs
			if !ok {
				if !gslutils.Contains(blocksMissingLogs, n) {
					return nil, errors.Wrap(superwatcher.ErrProcessReorg, "got headers but no result and not missing logs")
				}

				continue
			}

			// If headerHash differs from log blockHash, then the chain is reorging as we run this code
			if headerHash := header.Hash(); headerHash != mapResult.Block.Hash {
				return nil, errors.Wrapf(
					superwatcher.ErrFetchError, "header block hash %s on block %d is different than log block hash %s",
					headerHash.String(), n, mapResult.Block.Hash.String(),
				)
			}

			mapResult.Block.Header = header
		}
	} else {
		// We get headers for blocksMissingLogs regardless of deHeader values
		var err error
		headers, err = getHeadersByNumbers(ctx, client, blocksMissingLogs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get headers for blocks with missing logs")
		}
	}

	// i.e. if poller.doReorg == false
	if tracker == nil {
		return mapResults, nil
	}

	// Compare all known block hashes in tracker with result blocks.
	// Blocks with missing logs will be dealt with later in the last for-loop.
	for n := toBlock; n >= fromBlock; n-- {
		// Continue if block number |n| was never saved to tracker
		trackerBlock, ok := tracker.getTrackerBlock(n)
		if !ok {
			continue
		}

		mapResult, ok := mapResults[n]

		// If !ok, then it's one of blocksMissingLogs
		if !ok {
			// Mark blocks with missing logs as reorged. Do not overwrite tracker data here,
			// as poller.Poll logic needs the old info in the tracker to return as reorged block.
			// If you update the tracker data now, downstream users won't get all the old data
			// they need to handle chain reorg events.
			newHeader, ok := headers[n]
			if !ok {
				return nil, errors.Wrapf(superwatcher.ErrProcessReorg, "block %d (missing logs) header not found", n)
			}

			newHash := newHeader.Hash()
			if newHash == trackerBlock.Hash {
				continue
			}

			mapResult = &mapLogsResult{
				reorged: true,
				Block: superwatcher.Block{
					Number: n,
					Header: newHeader,
					Hash:   newHash,
					Logs:   nil,
				},
			}

			mapResults[n] = mapResult
		}

		if trackerBlock.Hash == mapResult.Block.Hash && len(trackerBlock.Logs) == len(mapResult.Block.Logs) {
			// If we have same blockHash with same logs length we can assume that the block was not reorged.
			continue
		}

		for _, trackerLog := range trackerBlock.Logs {
			trackerLog.Removed = true
		}

		mapResult.reorged = true
	}

	return mapResults, nil
}
