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
	forked bool // true if the tracker block hash differs from fresh block hash
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
	// policy changes poller behavior when getting headers.
	policy superwatcher.Policy,
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
				// Delete all results after block |number|
				deleteUnusableResult(mapResults, number)

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
	// so we have them in tracker but not in mapResults. We can know if a block has missing logs
	// only if we have the tracker.
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
	var targetBlocks []uint64

	// Get block headers for targetBlocks using BatchCallContext via getHeadersByNumbers
	if doHeader {
		for n := toBlock; n >= fromBlock; n-- {
			// Only add blocks that are not in blocksMissingLogs
			if gslutils.Contains(blocksMissingLogs, n) {
				continue
			}

			switch {
			case policy >= superwatcher.PolicyExpensive:
				// If PolicyExpensive, then we will get all headers in range [fromBlock, toBlock]
				// so here we append targetBlocks and continue.
				targetBlocks = append(targetBlocks, n)
				continue

			default:
				// Otherwise, only get headers for blocks with interesting logs (i.e. blocks present in mapResults)
				if _, ok := mapResults[n]; !ok {
					continue
				}

				targetBlocks = append(targetBlocks, n)
			}
		}

		var err error
		headers, err = getHeadersByNumbers(ctx, client, append(targetBlocks, blocksMissingLogs...))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get block headers for new logs filtered")
		}

		// Check if we got all headers we need
		var lenTargets int
		switch {
		case doHeader:
			lenTargets = len(targetBlocks) + len(blocksMissingLogs)
		case !doHeader:
			lenTargets = len(blocksMissingLogs)
		}
		if lenHeaders := len(headers); lenTargets != lenHeaders {
			return nil, errors.Wrapf(superwatcher.ErrFetchError, "expecting %d headers, got %d", lenTargets, lenHeaders)
		}

		// Collect headers into mapResults
		// (1) If Expensive, then collect all headers.
		// (2) If !Expensive, then only collect headers for blocks already in mapResults.
		// If header hashes with known hashes in mapResults differs, return errHashesDiffer.
		for n := fromBlock; n <= toBlock; n++ { // Go from fromBlock -> toBlock to detect the first bad blocks
			header, ok := headers[n]
			if !ok {
				// Expensive policy should have headers for all blocks
				if policy >= superwatcher.PolicyExpensive {
					return nil, errors.Wrapf(
						superwatcher.ErrFetchError,
						"policy is EXPENSIVE, but missing header for %d", n,
					)
				}

				continue
			}

			mapResult, ok := mapResults[n]
			// If no block n in mapResults (i.e. block n has to interesting logs this time)
			if !ok {
				switch {
				case policy >= superwatcher.PolicyExpensive:
					// Create new mapResult with data from header
					mapResult = &mapLogsResult{
						Block: superwatcher.Block{
							Number: n,
							Header: header,        // Assign header
							Hash:   header.Hash(), // Assign hash
						},
					}

					mapResults[n] = mapResult

				case policy <= superwatcher.PolicyNormal:
					// Check for superwatcher bug and just continue,
					// as we won't process headers for blocks outside of mapResults
					if !gslutils.Contains(blocksMissingLogs, n) {
						return nil, errors.Wrap(superwatcher.ErrProcessReorg, "got headers but no result and not missing logs")
					}

					continue
				}
			}

			// If code reaches here, then we have the header's block in mapResults, and so we'll compare their hashes.
			// If headerHash differs from log blockHash, then the chain is reorging as we run this code.
			if headerHash := header.Hash(); headerHash != mapResult.Block.Hash {
				// Delete all results after block mapResult.Number
				deleteUnusableResult(mapResults, mapResult.Number)

				return mapResults, errors.Wrapf(
					errHashesDiffer, "header block hash %s on block %d is different than log block hash %s",
					headerHash.String(), n, mapResult.Block.Hash.String(),
				)
			}

			// Assign header to good block in mapResult
			mapResult.Block.Header = header
		}
	} else {
		// We get headers for blocksMissingLogs anyway regardless of deHeader or policy values
		var err error
		headers, err = getHeadersByNumbers(ctx, client, blocksMissingLogs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get headers for blocks with missing logs")
		}
	}

	// i.e. if poller.doReorg == false,
	// return now after collecting headers and don't compare hashes
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
				forked: true,
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

		// Stamp Removed field
		for _, trackerLog := range trackerBlock.Logs {
			trackerLog.Removed = true
		}

		mapResult.forked = true
	}

	return mapResults, nil
}

// deleteUnusableResult removes all map keys that are >= lastGood.
// It used when poller.mapLogs encounter a situation where fresh chain data
// on a same block has different block hashes.
func deleteUnusableResult(mapResults map[uint64]*mapLogsResult, lastGood uint64) {
	for k := range mapResults {
		if k >= lastGood {
			delete(mapResults, k)
		}
	}
}
