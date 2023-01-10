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
// then gathers the result as `superwatcher.PollerResult`, and returns the result.
func (p *poller) Poll(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) (
	*superwatcher.PollerResult,
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
		until := p.lastRecordedBlock - p.filterRange
		p.debugger.Debug(2, "clearing tracker", zap.Uint64("untilBlock", until))
		p.tracker.clearUntil(until)
	}

	// mapLogs maps logs into a map of superwatcher.Block.
	// mapLogs also fetches other data, e.g. block headers, based on logs available and p.policy.
	// mapLogs will use tracker information as well as the block headers to detect chain reorg events.
	pollResults, err := mapLogs(
		ctx,
		fromBlock,
		toBlock,
		gslutils.CollectPointers(eventLogs), // Use pointers here, to avoid expensive copy
		p.doHeader || p.policy >= superwatcher.PolicyExpensive, // Force doHeader if policy >= Expensive
		p.tracker,
		p.client,
		p.policy,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error in mapLogs")
	}

	// Fills |result| and saves current data back to tracker first.
	result := new(superwatcher.PollerResult)
	for number := fromBlock; number <= toBlock; number++ {
		// Only blocks in mapResults are worth processing. There are 3 reasons a block is in mapResults:
		// (1) block has >=1 interesting log
		// (2) block _did_ have >= logs from the last call, but was reorged and no longer has any interesting logs
		// If (2), then it will removed from tracker, and will no longer appear in mapResults after this call.
		pollResult, ok := pollResults[number]
		if !ok {
			continue
		}

		// Reorged blocks (the ones that were removed) will be published with data from tracker
		if pollResult.forked && p.doReorg {
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
			freshHash := pollResult.Hash

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

			// Block used to have interesting logs, but chain reorg occurred
			// and its logs were moved to somewhere else, or just removed altogether.
			if trackerBlock.LogsMigrated {
				p.debugger.Debug(
					1, "logs missing from block, removing from tracker",
					zap.Uint64("blockNumber", number),
					zap.String("freshHash", freshHash.String()),
					zap.String("trackerHash", trackerBlock.String()),
				)

				switch {
				case p.policy == superwatcher.PolicyFast:
					// Remove from tracker if block has 0 logs, and poller will cease to
					// get block header for this empty block after this call.
					if err := p.tracker.removeBlock(number); err != nil {
						return nil, errors.Wrap(superwatcher.ErrProcessReorg, err.Error())
					}

				case p.policy >= superwatcher.PolicyNormal:
					// Save new empty block information back to tracker. This will make poller
					// continues to get header for this block until it goes out of filter (poll) scope.
					trackerBlock.Hash = freshHash
					trackerBlock.Logs = nil
				}
			}
		}

		freshBlock := pollResult.Block
		if p.doReorg {
			p.tracker.addTrackerBlock(&freshBlock)
		}

		// Copy goodBlock to avoid poller users mutating goodBlock values inside of tracker.
		goodBlock := freshBlock
		result.GoodBlocks = append(result.GoodBlocks, &goodBlock)
	}

	result.FromBlock, result.ToBlock = fromBlock, toBlock
	result.LastGoodBlock = superwatcher.LastGoodBlock(result)

	// Used for clearing tracker, and when doReorg policy changes (p.SetDoReorg(false))
	p.lastRecordedBlock = result.LastGoodBlock

	// If fromBlock reorged, return the result but with non-nil error.
	// This way, the result still gets emitted, leaving no unseen Block for the managed engine.
	fromBlockResult, ok := pollResults[fromBlock]
	if ok {
		if fromBlockResult.forked && p.doReorg {
			return result, errors.Wrapf(
				superwatcher.ErrFromBlockReorged, "fromBlock %d was removed/reorged", fromBlock,
			)
		}
	}

	return result, nil
}
