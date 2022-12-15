package emitter

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

// filterLogs filters Ethereum event logs from fromBlock to toBlock,
// and sends *types.Log and *lib.BlockInfo through w.logChan and w.reorgChan respectively.
// If an error is encountered, filterLogs returns with error.
// filterLogs should not be the one sending the error through w.errChan.
func (e *emitter) filterLogs(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) error {
	// Filter event logs with retries
	eventLogs, err := gslutils.RetryWithReturn(
		fmt.Sprintf("getLogs from %d to %d", fromBlock, toBlock),

		func() ([]types.Log, error) {
			// No error wrap because in retry mode
			return e.client.FilterLogs(ctx, ethereum.FilterQuery{ //nolint:wrapcheck
				FromBlock: big.NewInt(int64(fromBlock)),
				ToBlock:   big.NewInt(int64(toBlock)),
				Addresses: e.addresses,
				Topics:    e.topics,
			})
		},

		gslutils.Attempts(10),
		gslutils.Delay(4),
		gslutils.LastErrorOnly(true),
	)
	if err != nil {
		return errors.Wrap(errFetchError, err.Error())
	}

	// Wait here for logs and headers
	lenLogs := len(eventLogs)
	e.debugger.Debug(2, "got headers and logs", zap.Int("logs", lenLogs))

	// Clear all tracker's blocks before fromBlock - filterRange
	until := fromBlock - e.conf.FilterRange
	e.debugger.Debug(2, "clearing tracker", zap.Uint64("untilBlock", until))
	e.tracker.clearUntil(until)

	removedBlocks, mapFreshHashes, mapFreshLogs := mapLogs(
		fromBlock,
		toBlock,
		gslutils.CollectPointers(eventLogs), // Use pointers here, to avoid expensive copy
		e.tracker,
	)

	// Fills |filterResult| and saves current data back to tracker first.
	filterResult := new(superwatcher.FilterResult)
	for blockNumber := fromBlock; blockNumber <= toBlock; blockNumber++ {
		// Reorged blocks (the ones that were removed) will be published with data from tracker
		if removedBlocks[blockNumber] {
			reorgedBlock, foundInTracker := e.tracker.getTrackerBlockInfo(blockNumber)
			if !foundInTracker {
				logger.Panic(
					"blockInfo marked as reorged but was not found in tracker",
					zap.Uint64("blockNumber", blockNumber),
					zap.String("freshHash", reorgedBlock.String()),
				)
			}

			logger.Info(
				"chain reorg detected",
				zap.Uint64("blockNumber", blockNumber),
				zap.String("freshHash", mapFreshHashes[blockNumber].String()),
				zap.String("trackerHash", reorgedBlock.String()),
			)

			// ReorgedBlocks field contains the seen, removed blocks from tracker.
			// The new reroged blocks were added to |filterResult| as new GoodBlocks.
			// This means that the engine should process ReorgedBlocks first to revert the TXs,
			// before processing the new, reorged logs.
			filterResult.ReorgedBlocks = append(filterResult.ReorgedBlocks, reorgedBlock)
		}

		// New blocks will use fresh information. This includes new block after a reorg.
		logs, ok := mapFreshLogs[blockNumber]
		if !ok {
			continue
		}
		hash, ok := mapFreshHashes[blockNumber]
		if !ok {
			return errors.Wrapf(errNoHash, "blockNumber %d", blockNumber)
		}

		// Only add fresh, canonical blockInfo with interesting logs
		b := &superwatcher.BlockInfo{
			Number: blockNumber,
			Hash:   hash,
			Logs:   logs,
		}
		b.Logs = mapFreshLogs[blockNumber]
		e.tracker.addTrackerBlockInfo(b)

		// Publish only block with logs
		if b.Logs != nil {
			filterResult.GoodBlocks = append(filterResult.GoodBlocks, b)
		}
	}

	// If fromBlock was reorged, then return to loopFilterLogs.
	// Results will not be published, so the engine will never know that fromBlock is reorging.
	// **The reorged blocks different hashes have also been saved to tracker,
	// so if they come back in the next loop, with the same hash here, they'll be marked as non-reorg blocks**.
	if removedBlocks[fromBlock] {
		return errors.Wrapf(errFromBlockReorged, "fromBlock %d was removed/reorged", fromBlock)
	}

	// Publish filterResult via e.filterResultChan
	filterResult.FromBlock = fromBlock
	filterResult.ToBlock = toBlock
	filterResult.LastGoodBlock = lastGoodBlock(filterResult)
	e.emitFilterResult(filterResult)

	// End loop
	e.debugger.Debug(
		1, "number of logs published by filterLogs",
		zap.Int("eventLogs (filtered)", lenLogs),
		zap.Int("goodBlocks", len(filterResult.GoodBlocks)),
		zap.Int("reorgedBlocks", len(filterResult.ReorgedBlocks)),
		zap.Uint64("lastGoodBlock", filterResult.LastGoodBlock),
	)

	// Waits until engine syncs
	e.SyncsWithEngine()

	return nil
}

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
