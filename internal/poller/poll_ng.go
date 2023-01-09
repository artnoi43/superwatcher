package poller

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/artnoi43/gsl/concurrent"
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

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
	pollResults, err := poll(ctx, p.addresses, p.topics, fromBlock, toBlock, p.policy, p.client)
	if err != nil {
		return nil, err
	}

	mapResults, err := mapLogsNg(ctx, fromBlock, toBlock, p.policy, pollResults, p.client, p.tracker)
	if err != nil {
		return nil, err
	}

	result, err := findReorg(fromBlock, toBlock, p.policy, mapResults, p.tracker, p.debugger)
	if err != nil {
		return result, err
	}

	p.lastRecordedBlock = result.LastGoodBlock

	fromBlockResult, ok := mapResults[fromBlock]
	if ok {
		if fromBlockResult.forked && p.doReorg {
			return result, errors.Wrapf(
				superwatcher.ErrFromBlockReorged, "fromBlock %d was removed/reorged", fromBlock,
			)
		}
	}

	return result, nil
}

func poll( // nolint:unused
	ctx context.Context,
	addresses []common.Address,
	topics [][]common.Hash,
	fromBlock uint64,
	toBlock uint64,
	policy superwatcher.Policy,
	client superwatcher.EthClient,
) (
	map[uint64]superwatcher.Block,
	error,
) {
	pollResults := make(map[uint64]superwatcher.Block)

	q := ethereum.FilterQuery{ //nolint:wrapcheck
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: addresses,
		Topics:    topics,
	}

	switch {
	case policy >= superwatcher.PolicyExpensiveBlock: // Get blocks and event logs concurrently
		return nil, errors.New("PolicyExpensiveBlock not implemented")

	case policy == superwatcher.PolicyExpensive:
		// PolicyExpensive will get

		var headers map[uint64]superwatcher.BlockHeader
		var logs []types.Log

		var blockNumbers []uint64
		for n := fromBlock; n <= toBlock; n++ {
			blockNumbers = append(blockNumbers, n)
		}

		// Get block headers and event logs concurrently
		errChan := make(chan error)
		var wg sync.WaitGroup
		wg.Add(2) // Get headers and logs concurrently in 2 Goroutines

		go func() {
			defer wg.Done()
			var err error
			logs, err = client.FilterLogs(ctx, q)
			if err != nil {
				errChan <- errors.Wrapf(superwatcher.ErrFetchError, "filterLogs returned error")
			}
		}()

		go func() {
			defer wg.Done()
			var err error
			headers, err = getHeadersByNumbers(ctx, client, blockNumbers)
			if err != nil {
				errChan <- errors.Wrapf(err, "getHeadersByNumbers returned error")
			}
		}()

		if err := concurrent.WaitAndCollectErrors(&wg, errChan); err != nil {
			return nil, errors.Wrap(err, "concurrent fetch error")
		}

		if len(blockNumbers) != len(headers) {
			return nil, errors.Wrap(superwatcher.ErrFetchError, "headers and blockNumbers length not matched")
		}

		// Collect logs into map
		mappedLogs := make(map[uint64][]*types.Log)
		for i, log := range logs {
			blockLogs, ok := mappedLogs[log.BlockNumber]
			if !ok {
				mappedLogs[log.BlockNumber] = []*types.Log{&logs[i]}
				continue
			}

			blockLogs = append(blockLogs, &logs[i])
			mappedLogs[log.BlockNumber] = blockLogs
		}

		// Collect pollResults
		for n := fromBlock; n <= toBlock; n++ {
			pollResult, ok := pollResults[n]
			if ok {
				panic(fmt.Sprintf("duplicate pollResults on block %d", n))
			}

			header, ok := headers[n]
			if !ok {
				return nil, errors.Wrapf(superwatcher.ErrSuperwatcherBug, "missing block header %d", n)
			}

			headerHash := header.Hash()
			logs, ok := mappedLogs[n]
			if ok {
				// If this block has logs, see if their blockHash matches headerHash
				for _, log := range logs {
					if log.BlockHash != headerHash {
						deleteMapResults(pollResults, n)
						return nil, errors.Wrapf(
							errHashesDiffer, "block %d header has different hash than log blockHash: %s vs %s",
							log.BlockNumber, hashStr(headerHash), hashStr(log.BlockHash),
						)
					}
				}

				pollResult.Logs = logs
			}

			pollResult.Hash = headerHash
			pollResult.Header = header
			pollResults[n] = pollResult
		}

	case policy <= superwatcher.PolicyNormal:
		// Just get event logs
		logs, err := client.FilterLogs(ctx, q)
		if err != nil {
			return nil, errors.Wrap(superwatcher.ErrFetchError, err.Error())
		}

		// Collect pollResults
		for i, log := range logs {
			b, ok := pollResults[log.BlockNumber]
			if !ok {
				b = superwatcher.Block{
					Number: log.BlockNumber,
					Hash:   log.BlockHash,
					Header: nil,
					Logs:   []*types.Log{&logs[i]},
				}

				pollResults[log.BlockNumber] = b
				continue
			} else if ok {
				if b.Hash != log.BlockHash {
					return nil, errors.Wrapf(
						errHashesDiffer, "logs block hashes on block %d differ: %s vs %s",
						log.BlockNumber, b.Hash.String(), log.BlockHash.String(),
					)
				}

				b.Logs = append(b.Logs, &logs[i])
			}
		}
	}

	return pollResults, nil
}

func mapLogsNg(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
	policy superwatcher.Policy,
	pollResults map[uint64]superwatcher.Block,
	client superwatcher.EthClient,
	tracker *blockTracker,
) (
	map[uint64]*mapLogsResult,
	error,
) {
	var missingLogs []uint64
	if tracker != nil {
		for n := fromBlock; n <= toBlock; n++ {
			_, ok := tracker.getTrackerBlock(n)
			if !ok {
				continue
			}

			_, ok = pollResults[n]
			if !ok {
				missingLogs = append(missingLogs, n)
			}
		}
	}

	mapResults := make(map[uint64]*mapLogsResult)

	var blocksWithResults []uint64
	if policy < superwatcher.PolicyExpensive {
		// PolicyExpensive already fetched block headers fpr pollResults in poller.poll
		for n := fromBlock; n <= toBlock; n++ {
			pollResult, ok := pollResults[n]
			if !ok {
				continue
			}

			// Collect pollResult into mapResults
			mapResults[n] = &mapLogsResult{
				Block: pollResult,
			}

			// We will get headers for missingLogs anyway, so don't append if n's also there
			if gslutils.Contains(missingLogs, n) {
				continue
			}

			blocksWithResults = append(blocksWithResults, n)
		}
	}

	headers, err := getHeadersByNumbers(ctx, client, append(blocksWithResults, missingLogs...))
	if err != nil {
		return nil, errors.Wrap(superwatcher.ErrFetchError, "failed to get block headers in mapLogsNg")
	}

	lenHeads, lenBlocks := len(headers), len(blocksWithResults)+len(missingLogs)
	if lenHeads != lenBlocks {
		return nil, errors.Wrapf(superwatcher.ErrFetchError, "expecting %d headers, got %d", lenBlocks, lenHeads)
	}

	if tracker == nil {
		return mapResults, nil
	}

	// Detect chain reorg using tracker
	if tracker != nil {
		for n := fromBlock; n <= toBlock; n++ {
			trackerBlock, ok := tracker.getTrackerBlock(n)
			if !ok {
				continue
			}

			pollResult, ok := pollResults[n]
			if !ok {
				// If !ok and PolicyNormal, then n is one of the missing logs
				if policy <= superwatcher.PolicyNormal {
					continue
				}

				return nil, errors.Wrapf(superwatcher.ErrSuperwatcherBug, "missing pollResult for block %d with PolicyExpensive", n)
			}

			var header superwatcher.BlockHeader
			switch {
			case policy == superwatcher.PolicyExpensive:

				if gslutils.Contains(missingLogs, n) {
					header = headers[n]
					if pollResult.Header.Hash() != header.Hash() {
						// Our 2 most recent hashes for block |n| are unusable because they differ.
						deleteMapResults(mapResults, n)
						return mapResults, nil
					}
				} else {
					header = pollResult.Header
				}

			case policy < superwatcher.PolicyExpensive:
				header = headers[n]
			}

			if header == nil {
				return nil, errors.Wrapf(superwatcher.ErrSuperwatcherBug, "missing header for block %d with policy %s", n, policy.String())
			}

			headerHash := header.Hash()
			if header.Hash() != pollResult.Hash {
				deleteMapResults(mapResults, n)

				return mapResults, errors.Wrapf(
					errHashesDiffer, "block %d header hash differs from log block hash: %s vs %s",
					n, hashStr(headerHash), hashStr(pollResult.Hash),
				)
			}

			mapResult, ok := mapResults[n]
			if !ok || mapResult == nil {
				mapResult = new(mapLogsResult)
				mapResults[n] = mapResult
			}

			mapResult.Header = header
			mapResult.Hash = headerHash

			if trackerBlock.Hash == pollResult.Hash && len(trackerBlock.Logs) == len(pollResult.Logs) {
				continue
			}

			mapResult.forked = true
		}
	}

	return mapResults, nil
}

func findReorg(
	fromBlock uint64,
	toBlock uint64,
	policy superwatcher.Policy,
	mapResults map[uint64]*mapLogsResult,
	tracker *blockTracker,
	debugger *debugger.Debugger,
) (*superwatcher.PollerResult, error) {
	// Fills |result| and saves current data back to tracker first.
	result := new(superwatcher.PollerResult)
	for number := fromBlock; number <= toBlock; number++ {
		// Only blocks in mapResults are worth processing. There are 3 reasons a block is in mapResults:
		// (1) block has >=1 interesting log
		// (2) block _did_ have >= logs from the last call, but was reorged and no longer has any interesting logs
		// If (2), then it will removed from tracker, and will no longer appear in mapResults after this call.
		mapResult, ok := mapResults[number]
		if !ok {
			continue
		}

		// Reorged blocks (the ones that were removed) will be published with data from tracker
		if mapResult.forked && tracker != nil {
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
			freshHash := mapResult.Hash

			debugger.Debug(
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
				debugger.Debug(
					1, "logs missing from block, removing from tracker",
					zap.Uint64("blockNumber", number),
					zap.String("freshHash", freshHash.String()),
					zap.String("trackerHash", trackerBlock.String()),
				)

				switch {
				case policy == superwatcher.PolicyFast:
					// Remove from tracker if block has 0 logs, and poller will cease to
					// get block header for this empty block after this call.
					if err := tracker.removeBlock(number); err != nil {
						return nil, errors.Wrap(superwatcher.ErrProcessReorg, err.Error())
					}

				default:
					// Save new empty block information back to tracker. This will make poller
					// continues to get header for this block until it goes out of filter (poll) scope.
					trackerBlock.Hash = freshHash
					trackerBlock.Logs = nil
				}
			}
		}

		goodBlock := mapResult.Block

		if tracker != nil {
			tracker.addTrackerBlock(&goodBlock)
		}

		// Copy goodBlock to avoid poller users mutating goodBlock values inside of tracker.
		resultBlock := goodBlock
		result.GoodBlocks = append(result.GoodBlocks, &resultBlock)
	}

	result.FromBlock, result.ToBlock = fromBlock, toBlock
	result.LastGoodBlock = superwatcher.LastGoodBlock(result)

	// If fromBlock reorged, return the result but with non-nil error.
	// This way, the result still gets emitted, leaving no unseen Block for the managed engine.
	fromBlockResult, ok := mapResults[fromBlock]
	if ok {
		if fromBlockResult.forked && tracker != nil {
			return result, errors.Wrapf(
				superwatcher.ErrFromBlockReorged, "fromBlock %d was removed/reorged", fromBlock,
			)
		}
	}

	return result, nil
}
