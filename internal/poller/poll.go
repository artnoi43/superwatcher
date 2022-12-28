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
			return p.client.FilterLogs(ctx, ethereum.FilterQuery{ //nolint:wrapcheck
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

	removedBlocks, mapFreshHashes, mapFreshLogs, err := mapLogs(
		ctx,
		fromBlock,
		toBlock,
		gslutils.CollectPointers(eventLogs), // Use pointers here, to avoid expensive copy
		p.tracker,
		p.client.HeaderByNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error in mapLogs")
	}

	// Fills |result| and saves current data back to tracker first.
	result := new(superwatcher.FilterResult)
	for blockNumber := fromBlock; blockNumber <= toBlock; blockNumber++ {
		// Reorged blocks (the ones that were removed) will be published with data from tracker
		if removedBlocks[blockNumber] && p.doReorg {
			trackerBlock, foundInTracker := p.tracker.getTrackerBlockInfo(blockNumber)
			if !foundInTracker {
				p.debugger.Debug(
					1, "blockInfo marked as reorged but was not found in tracker",
					zap.Uint64("blockNumber", blockNumber),
					zap.String("freshHash", trackerBlock.String()),
				)

				return nil, errors.Wrapf(
					superwatcher.ErrProcessReorg, "reorgedBlock %d not found in tracker", blockNumber,
				)
			}

			// Logs may be moved from blockNumber, hence there's no value in map
			freshHash, ok := mapFreshHashes[blockNumber]
			if !ok {
				return nil, errors.Wrap(superwatcher.ErrProcessReorg, err.Error())
			}

			p.debugger.Debug(
				1, "chain reorg detected",
				zap.Uint64("blockNumber", blockNumber),
				zap.String("freshHash", freshHash.String()),
				zap.String("trackerHash", trackerBlock.String()),
			)

			// Copy to allow us to update blocks with LogsMigrated inside the tracker
			// without mutating values in |result.ReorgedBlocks|
			copiedFromTracker := *trackerBlock
			result.ReorgedBlocks = append(result.ReorgedBlocks, &copiedFromTracker)

			// Save updated recent block info back to tracker (there won't be this block in mapFreshLogs below)
			if trackerBlock.LogsMigrated {
				p.debugger.Debug(
					1, "logs missing from block, updating trackerBlock with nil logs and freshHash",
					zap.Uint64("blockNumber", blockNumber),
					zap.String("freshHash", freshHash.String()),
					zap.String("trackerHash", trackerBlock.String()),
				)

				// Update values in tracker
				trackerBlock.Hash = freshHash
				trackerBlock.Logs = nil
			}
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
	result.LastGoodBlock = superwatcher.LastGoodBlock(result)

	p.lastRecordedBlock = result.LastGoodBlock

	return result, nil
}
