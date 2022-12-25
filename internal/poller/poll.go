package poller

import (
	"context"
	"fmt"
	"math/big"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

// Poll filters Ethereum event logs from fromBlock to toBlock,
// and sends *types.Log and *superwatcher.BlockInfo through w.logChan and w.reorgChan respectively.
// If an error is encountered, Poll returns with error. Poll should never be the one sending the error through w.errChan.
func (p *poller) Poll(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) (
	*superwatcher.FilterResult,
	error,
) {
	p.Lock()
	defer p.Unlock()

	// Filter event logs with retries
	eventLogs, err := gslutils.RetryWithReturn(
		fmt.Sprintf("getLogs from %d to %d", fromBlock, toBlock),

		func() ([]types.Log, error) {
			// No error wrap because in retry mode
			return p.filterFunc(ctx, ethereum.FilterQuery{ //nolint:wrapcheck
				FromBlock: big.NewInt(int64(fromBlock)),
				ToBlock:   big.NewInt(int64(toBlock)),
				Addresses: p.addresses,
				Topics:    p.topics,
			})
		},

		gslutils.Attempts(10),
		gslutils.Delay(4),
		gslutils.LastErrorOnly(true),
	)
	if err != nil {
		return nil, errors.Wrap(superwatcher.ErrFetchError, err.Error())
	}

	// Wait here for logs and headers
	lenLogs := len(eventLogs)
	p.debugger.Debug(2, "got headers and logs", zap.Int("logs", lenLogs))

	if p.tracker != nil {
		// Clear all tracker's blocks before fromBlock - filterRange
		until := fromBlock - p.filterRange
		p.debugger.Debug(2, "clearing tracker", zap.Uint64("untilBlock", until))
		p.tracker.clearUntil(until)
	}

	removedBlocks, mapFreshHashes, mapFreshLogs := mapLogs(
		fromBlock,
		toBlock,
		gslutils.CollectPointers(eventLogs), // Use pointers here, to avoid expensive copy
		p.tracker,
	)

	// Fills |result| and saves current data back to tracker first.
	result := new(superwatcher.FilterResult)
	for blockNumber := fromBlock; blockNumber <= toBlock; blockNumber++ {
		// Reorged blocks (the ones that were removed) will be published with data from tracker
		if removedBlocks[blockNumber] && p.doReorg {
			reorgedBlock, foundInTracker := p.tracker.getTrackerBlockInfo(blockNumber)
			if !foundInTracker {
				p.debugger.Warn(
					1, "blockInfo marked as reorged but was not found in tracker",
					zap.Uint64("blockNumber", blockNumber),
					zap.String("freshHash", reorgedBlock.String()),
				)

				return nil, errors.Wrapf(superwatcher.ErrProcessReorg, "reorgedBlock %d not found in tracker", blockNumber)
			}

			logger.Info(
				"chain reorg detected",
				zap.Uint64("blockNumber", blockNumber),
				zap.String("freshHash", mapFreshHashes[blockNumber].String()),
				zap.String("trackerHash", reorgedBlock.String()),
			)

			// ReorgedBlocks field contains the seen, removed blocks from tracker.
			// The new reroged blocks were added to |result| as new `result.GoodBlocks`.
			// This means that the engine should process ReorgedBlocks first to revert the TXs,
			// before processing the new, reorged logs.
			result.ReorgedBlocks = append(result.ReorgedBlocks, reorgedBlock)
		}

		// New blocks will use fresh information. This includes new block after a reorg.
		logs, ok := mapFreshLogs[blockNumber]
		if !ok || len(logs) == 0 {
			continue
		}
		hash, ok := mapFreshHashes[blockNumber]
		if !ok {
			return nil, errors.Wrapf(errNoHash, "blockNumber %d", blockNumber)
		}

		// Only add fresh, canonical blockInfo with interesting logs
		b := &superwatcher.BlockInfo{
			Number: blockNumber,
			Hash:   hash,
			Logs:   logs,
		}

		if p.doReorg {
			p.tracker.addTrackerBlockInfo(b)
		}

		// Publish only block with logs
		result.GoodBlocks = append(result.GoodBlocks, b)
	}

	// If fromBlock was reorged, then return to loopEmit.
	// Results will not be published, so the engine will never know that fromBlock is reorging.
	// **The reorged blocks different hashes have also been saved to tracker,
	// so if they come back in the next loop, with the same hash here, they'll be marked as non-reorg blocks**.
	if removedBlocks[fromBlock] && p.doReorg {
		return nil, errors.Wrapf(superwatcher.ErrFromBlockReorged, "fromBlock %d was removed/reorged", fromBlock)
	}

	// Publish filterResult via e.filterResultChan
	result.FromBlock = fromBlock
	result.ToBlock = toBlock
	result.LastGoodBlock = lastGoodBlock(result)

	p.lastRecordedBlock = result.LastGoodBlock

	return result, nil
}

// lastGoodBlock computes `superwatcher.FilterResult.lastGoodBlock` based on |result|.
// TODO: Finalize or just remove FilterResult.LastGoodBlock altogether.
func lastGoodBlock(
	result *superwatcher.FilterResult,
) uint64 {
	if len(result.ReorgedBlocks) != 0 {
		// If there's also goodBlocks during reorg
		if l := len(result.GoodBlocks); l != 0 {
			// Use last good block's number as LastGoodBlock
			lastGood := result.GoodBlocks[l-1].Number
			firstReorg := result.ReorgedBlocks[0].Number

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
