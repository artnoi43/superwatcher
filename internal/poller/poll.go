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

// Poll filters Ethereum event logs from |fromBlock| to |toBlock|,
// then gathers the result as `superwatcher.PollResult`, and returns the result.
func (p *poller) Poll(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) (
	*superwatcher.PollResult,
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

	mapResults, err := mapLogs(
		ctx,
		fromBlock,
		toBlock,
		gslutils.CollectPointers(eventLogs), // Use pointers here, to avoid expensive copy
		p.doHeader,
		p.tracker,
		p.client,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error in mapLogs")
	}

	wasReorged := make(map[uint64]bool)
	for number, mapResult := range mapResults {
		wasReorged[number] = mapResult.reorged
	}

	// Fills |result| and saves current data back to tracker first.
	result := new(superwatcher.PollResult)
	for number := fromBlock; number <= toBlock; number++ {
		// Reorged blocks (the ones that were removed) will be published with data from tracker
		if wasReorged[number] && p.doReorg {
			trackerBlock, ok := p.tracker.getTrackerBlock(number)
			if !ok {
				p.debugger.Debug(
					1, "block marked as reorged but was not found in tracker",
					zap.Uint64("blockNumber", number),
					zap.String("freshHash", trackerBlock.String()),
				)

				return nil, errors.Wrapf(
					superwatcher.ErrProcessReorg, "reorgedBlock %d not found in tracker", number,
				)
			}

			// Logs may be moved from blockNumber, hence there's no value in map
			freshHash := mapResults[number].Block.Hash

			p.debugger.Debug(
				1, "chain reorg detected",
				zap.Uint64("blockNumber", number),
				zap.String("freshHash", freshHash.String()),
				zap.String("trackerHash", trackerBlock.String()),
			)

			// Copy to avoid mutated trackerBlock which might break poller logic.
			// After the copy, result.ReorgedBlocks consumer may freely mutate their *Block.
			copiedFromTracker := *trackerBlock
			result.ReorgedBlocks = append(result.ReorgedBlocks, &copiedFromTracker)

			// Save updated recent block info back to tracker (there won't be this block in mapFreshLogs below)
			if trackerBlock.LogsMigrated {
				p.debugger.Debug(
					1, "logs missing from block, updating trackerBlock with nil logs and freshHash",
					zap.Uint64("blockNumber", number),
					zap.String("freshHash", freshHash.String()),
					zap.String("trackerHash", trackerBlock.String()),
				)

				// Remove from tracker if there's no logs, to avoid mapLogs getting headers
				// for empty blocks more than once.
				if err := p.tracker.removeBlock(number); err != nil {
					return nil, errors.Wrap(superwatcher.ErrProcessReorg, err.Error())
				}
			}
		}

		// New blocks will use fresh information. This includes new block after a reorg.
		// Only block with logs will be appended to |result.GoodBlocks|.
		mapResult, ok := mapResults[number]
		if !ok {
			continue
		}

		goodBlock := mapResult.Block

		if p.doReorg {
			p.tracker.addTrackerBlock(&goodBlock)
		}

		// Copy goodBlock to avoid poller users mutating goodBlock values inside of tracker.
		resultBlock := goodBlock
		result.GoodBlocks = append(result.GoodBlocks, &resultBlock)
	}

	result.FromBlock, result.ToBlock = fromBlock, toBlock
	result.LastGoodBlock = superwatcher.LastGoodBlock(result)

	p.lastRecordedBlock = result.LastGoodBlock

	// If fromBlock reorged, return the result but with non-nil error.
	// This way, the result still gets emitted, leaving no unseen Block for the managed engine.
	if wasReorged[fromBlock] && p.doReorg {
		return result, errors.Wrapf(
			superwatcher.ErrFromBlockReorged, "fromBlock %d was removed/reorged", fromBlock,
		)
	}

	return result, nil
}
