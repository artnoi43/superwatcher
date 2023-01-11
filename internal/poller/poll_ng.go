package poller

// TODO: remove fmt debug prints

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/gsl"
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

func (p *poller) PollNg(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) (
	*superwatcher.PollerResult,
	error,
) {
	p.Lock()
	defer p.Unlock()

	pollResults := make(map[uint64]*mapLogsResult)
	pollResults, err := poll(ctx, fromBlock, toBlock, p.addresses, p.topics, p.policy, p.client, pollResults)
	if err != nil {
		return nil, err
	}

	pollResults, err = pollMissing(ctx, fromBlock, toBlock, p.policy, p.client, p.tracker, pollResults)
	if err != nil {
		return nil, err
	}

	result, err := pollerResult(fromBlock, toBlock, p.policy, p.tracker, p.debugger, pollResults)
	if err != nil {
		return result, err
	}

	result.FromBlock, result.ToBlock = fromBlock, toBlock
	result.LastGoodBlock = superwatcher.LastGoodBlock(result)

	p.lastRecordedBlock = result.LastGoodBlock

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

func poll(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
	addresses []common.Address,
	topics [][]common.Hash,
	policy superwatcher.Policy,
	client superwatcher.EthClient,
	pollResults map[uint64]*mapLogsResult,
) (
	map[uint64]*mapLogsResult,
	error,
) {
	switch {
	case policy >= superwatcher.PolicyExpensiveBlock: // Get blocks and event logs concurrently
		return nil, errors.New("PolicyExpensiveBlock not implemented")

	case policy == superwatcher.PolicyExpensive:
		return pollExpensive(ctx, fromBlock, toBlock, addresses, topics, client, pollResults)

	case policy <= superwatcher.PolicyNormal:
		return pollCheap(ctx, fromBlock, toBlock, addresses, topics, client, pollResults)
	}

	panic("invalid policy " + policy.String())
}

func pollMissing(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
	policy superwatcher.Policy,
	client superwatcher.EthClient,
	tracker *blockTracker, // tracker is used as read-only in here. Don't write.
	pollResults map[uint64]*mapLogsResult,
) (
	map[uint64]*mapLogsResult,
	error,
) {
	// Find missing blocks (blocks in tracker that are not in pollResults)
	var blocksMissing []uint64
	if tracker != nil {
		for n := toBlock; n >= fromBlock; n-- {
			if _, ok := tracker.getTrackerBlock(n); !ok {
				continue
			}
			if _, ok := pollResults[n]; ok {
				continue
			}

			blocksMissing = append(blocksMissing, n)
		}
	}

	fmt.Println()
	fmt.Println("blocksMissing", blocksMissing)

	headers, err := getHeadersByNumbers(ctx, client, blocksMissing)
	if err != nil {
		return nil, errors.Wrap(superwatcher.ErrFetchError, "failed to get block headers in mapLogsNg")
	}

	lenHeads, lenBlocks := len(headers), len(blocksMissing)
	if lenHeads != lenBlocks {
		return nil, errors.Wrapf(superwatcher.ErrFetchError, "expecting %d headers, got %d", lenBlocks, lenHeads)
	}

	// Collect headers for blocksMissing
	_, err = collectHeaders(pollResults, fromBlock, toBlock, headers)
	if err != nil {
		if errors.Is(err, errHashesDiffer) {
			// deleteMapResults(pollResults, lastGood)
			return pollResults, err
		}

		return nil, errors.Wrap(err, "collectHeaders error")
	}

	if tracker == nil {
		return pollResults, nil
	}

	// Detect chain reorg using tracker
	for n := fromBlock; n <= toBlock; n++ {
		trackerBlock, ok := tracker.getTrackerBlock(n)
		if !ok {
			continue
		}

		pollResult, ok := pollResults[n]
		if !ok {
			return nil, errors.Wrapf(superwatcher.ErrProcessReorg, "pollResult missing for trackerBlock %d", n)
		}

		if trackerBlock.Hash == pollResult.Hash && len(trackerBlock.Logs) == len(pollResult.Logs) {
			continue
		}

		if gsl.Contains(blocksMissing, n) {
			pollResult.LogsMigrated = true
		}

		pollResult.forked = true
	}

	return pollResults, nil
}

func pollerResult(
	fromBlock uint64,
	toBlock uint64,
	policy superwatcher.Policy,
	tracker *blockTracker,
	debugger *debugger.Debugger,
	pollResults map[uint64]*mapLogsResult,
) (*superwatcher.PollerResult, error) {
	// Fills |result| and saves current data back to tracker first.
	result := new(superwatcher.PollerResult)
	for number := fromBlock; number <= toBlock; number++ {
		// Only blocks in pollResults are worth processing. There are 3 reasons a block is in pollResults:
		// (1) block has >=1 interesting log
		// (2) block _did_ have >= logs from the last call, but was reorged and no longer has any interesting logs
		// If (2), then it will removed from tracker, and will no longer appear in pollResults after this call.
		pollResult, ok := pollResults[number]
		if !ok {
			continue
		}

		// Reorged blocks (the ones that were removed) will be published with data from tracker
		if pollResult.forked && tracker != nil {
			trackerBlock, ok := tracker.getTrackerBlock(number)
			if !ok {
				debugger.Debug(
					1, "block marked as reorged but was not found in tracker",
					zap.Uint64("blockNumber", number),
					zap.String("freshHash", trackerBlock.String()),
				)

				return nil, errors.Wrapf(
					superwatcher.ErrProcessReorg, "reorgedBlock %d not found in trackfromBlocker", number,
				)
			}

			// Logs may be moved from blockNumber, hence there's no value in map
			freshHash := pollResult.Hash
			debugger.Debug(
				1, "chain reorg detected",
				zap.Uint64("blockNumber", number),
				zap.String("freshHash", freshHash.String()),
				zap.String("trackerHash", trackerBlock.String()),
				zap.Int("freshLogs", len(pollResult.Logs)),
				zap.Int("trackerLogs", len(trackerBlock.Logs)),
			)

			// Copy to avoid mutated trackerBlock which might break poller logic.
			// After the copy, result.ReorgedBlocks consumer may freely mutate their *Block.
			copiedFromTracker := *trackerBlock
			result.ReorgedBlocks = append(result.ReorgedBlocks, &copiedFromTracker)

			// Block used to have interesting logs, but chain reorg occurred
			// and its logs were moved to somewhere else, or just removed altogether.
			if pollResult.LogsMigrated {
				debugger.Debug(
					1, "logs missing from block",
					zap.Uint64("blockNumber", number),
					zap.String("freshHash", freshHash.String()),
					zap.String("trackerHash", trackerBlock.String()),
					zap.Int("old logs", len(trackerBlock.Logs)),
				)

				err := handleBlocksMissingPolicy(number, tracker, trackerBlock, freshHash, policy)
				if err != nil {
					return nil, errors.Wrap(superwatcher.ErrProcessReorg, err.Error())
				}
			}
		}

		freshBlock := pollResult.Block
		freshBlockTxHashes := gsl.Map(freshBlock.Logs, func(l *types.Log) (string, bool) {
			return gsl.StringerToLowerString(l.TxHash), true
		})

		fmt.Println("pollerResult: freshBlock", freshBlock.Number, freshBlockTxHashes)
		addTrackerBlockPolicy(tracker, &freshBlock, policy)

		// Copy goodBlock to avoid poller users mutating goodBlock values inside of tracker.
		goodBlock := freshBlock
		result.GoodBlocks = append(result.GoodBlocks, &goodBlock)
	}

	return result, nil
}

func handleBlocksMissingPolicy(
	number uint64,
	tracker *blockTracker,
	trackerBlock *superwatcher.Block,
	freshHash common.Hash,
	policy superwatcher.Policy,
) error {
	switch {
	case policy == superwatcher.PolicyFast:
		// Remove from tracker if block has 0 logs, and poller will cease to
		// get block header for this empty block after this call.
		fmt.Println("removing block", number, trackerBlock.String(), len(trackerBlock.Logs))
		if err := tracker.removeBlock(number); err != nil {
			return errors.Wrap(superwatcher.ErrProcessReorg, err.Error())
		}

	default:
		// Save new empty block information back to tracker. This will make poller
		// continues to get header for this block until it goes out of filter (poll) scope.
		trackerBlock.Hash = freshHash
		trackerBlock.Logs = nil
	}

	return nil
}

func addTrackerBlockPolicy(tracker *blockTracker, block *superwatcher.Block, policy superwatcher.Policy) {
	if tracker == nil {
		return
	}
	if policy == superwatcher.PolicyFast {
		if len(block.Logs) == 0 {
			return
		}
	}

	tracker.addTrackerBlock(block)
}
