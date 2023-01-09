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

	pollResults, err = findMissing(ctx, fromBlock, toBlock, p.policy, pollResults, p.client, p.tracker)
	if err != nil {
		return nil, err
	}

	result, err := collect(fromBlock, toBlock, p.policy, pollResults, p.tracker, p.debugger)
	if err != nil {
		return result, err
	}

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
	addresses []common.Address,
	topics [][]common.Hash,
	fromBlock uint64,
	toBlock uint64,
	policy superwatcher.Policy,
	client superwatcher.EthClient,
) (
	map[uint64]*mapLogsResult,
	error,
) {
	pollResults := make(map[uint64]*mapLogsResult)

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
		resultLogs := make(map[uint64][]*types.Log)
		for i, log := range logs {
			blockLogs, ok := resultLogs[log.BlockNumber]
			if !ok {
				resultLogs[log.BlockNumber] = []*types.Log{&logs[i]}
				continue
			}

			blockLogs = append(blockLogs, &logs[i])
			resultLogs[log.BlockNumber] = blockLogs
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
			logs, ok := resultLogs[n]
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
			result, ok := pollResults[log.BlockNumber]
			if !ok {
				b := superwatcher.Block{
					Number: log.BlockNumber,
					Hash:   log.BlockHash,
					Header: nil,
					Logs:   []*types.Log{&logs[i]},
				}

				pollResults[log.BlockNumber] = &mapLogsResult{Block: b}
				continue
			} else if ok {
				if result.Hash != log.BlockHash {
					return nil, errors.Wrapf(
						errHashesDiffer, "logs block hashes on block %d differ: %s vs %s",
						log.BlockNumber, result.Hash.String(), log.BlockHash.String(),
					)
				}

				result.Logs = append(result.Logs, &logs[i])
			}
		}
	}

	return pollResults, nil
}

func findMissing(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
	policy superwatcher.Policy,
	pollResults map[uint64]*mapLogsResult,
	client superwatcher.EthClient,
	tracker *blockTracker, // tracker is used as read-only in here. Don't write.
) (
	map[uint64]*mapLogsResult,
	error,
) {
	var blocksMissing []uint64
	if tracker != nil {
		for n := fromBlock; n <= toBlock; n++ {
			_, ok := tracker.getTrackerBlock(n)
			if !ok {
				continue
			}

			_, ok = pollResults[n]
			if !ok {
				blocksMissing = append(blocksMissing, n)
			}
		}
	}

	var blocksWithResults []uint64
	if policy < superwatcher.PolicyExpensive {
		// PolicyExpensive already fetched block headers for pollResults in poller.poll
		for n := fromBlock; n <= toBlock; n++ {
			_, ok := pollResults[n]
			if !ok {
				continue
			}

			// We will get headers for missingLogs anyway, so don't append if n's also there
			if gslutils.Contains(blocksMissing, n) {
				continue
			}

			blocksWithResults = append(blocksWithResults, n)
		}
	}

	headers, err := getHeadersByNumbers(ctx, client, append(blocksWithResults, blocksMissing...))
	if err != nil {
		return nil, errors.Wrap(superwatcher.ErrFetchError, "failed to get block headers in mapLogsNg")
	}

	lenHeads, lenBlocks := len(headers), len(blocksWithResults)+len(blocksMissing)
	if lenHeads != lenBlocks {
		return nil, errors.Wrapf(superwatcher.ErrFetchError, "expecting %d headers, got %d", lenBlocks, lenHeads)
	}

	if tracker == nil {
		return pollResults, nil
	}

	// Detect chain reorg using tracker
	if tracker != nil {
		for n := fromBlock; n <= toBlock; n++ {
			trackerBlock, ok := tracker.getTrackerBlock(n)
			if !ok {
				continue
			}

			pollResult, ok := pollResults[n]

			// !ok is when we had this block in tracker but not in pollResults
			var header superwatcher.BlockHeader
			var headerHash common.Hash
			if !ok {
				// It should have been tagged as missing
				if !gslutils.Contains(blocksMissing, n) {
					return nil, errors.Wrapf(superwatcher.ErrSuperwatcherBug, "missing pollResult for trackerBlock %d", n)
				}

				// We should ALWAYS have headers for blocksMissing
				header, ok = headers[n]
				if !ok {
					return nil, errors.Wrapf(superwatcher.ErrSuperwatcherBug, "missing headers for blocksMissing in tracker block %d", n)
				}

				headerHash = header.Hash()
				pollResult = &mapLogsResult{
					Block: superwatcher.Block{
						Number:       n,
						Header:       header,
						Hash:         headerHash,
						LogsMigrated: true,
					},

					forked: true,
				}

				pollResults[n] = pollResult
				continue
			}

			// ok means we have this block in tracker, and in pollResult. So here we check if
			// pollResult.Hash differs from header.Hash
			header, ok = headers[n]
			if !ok {
				return nil, errors.Wrapf(superwatcher.ErrSuperwatcherBug, "missing header for pollResult")
			}

			headerHash = header.Hash()
			if pollHash := pollResult.Hash; headerHash != pollHash {
				deleteMapResults(pollResults, n)
				return pollResults, errors.Wrapf(
					errHashesDiffer, "block %d headerHash %s differs from pollResultHash %s",
					n, headerHash.String(), pollHash.String(),
				)
			}

			pollResult.Header = header
			pollResults[n] = pollResult

			if header.Hash() != pollResult.Hash {
				deleteMapResults(pollResults, n)
				return pollResults, errors.Wrapf(
					errHashesDiffer, "block %d header hash differs from log block hash: %s vs %s",
					n, hashStr(headerHash), hashStr(pollResult.Hash),
				)
			}

			if trackerBlock.Hash == pollResult.Hash && len(trackerBlock.Logs) == len(pollResult.Logs) {
				continue
			}

			pollResult.forked = true
		}
	}

	return pollResults, nil
}

func collect(
	fromBlock uint64,
	toBlock uint64,
	policy superwatcher.Policy,
	pollResults map[uint64]*mapLogsResult,
	tracker *blockTracker,
	debugger *debugger.Debugger,
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

		freshBlock := pollResult.Block

		if tracker != nil {
			tracker.addTrackerBlock(&freshBlock)
		}

		// Copy goodBlock to avoid poller users mutating goodBlock values inside of tracker.
		goodBlock := freshBlock
		result.GoodBlocks = append(result.GoodBlocks, &goodBlock)
	}

	result.FromBlock, result.ToBlock = fromBlock, toBlock
	result.LastGoodBlock = superwatcher.LastGoodBlock(result)

	// If fromBlock reorged, return the result but with non-nil error.
	// This way, the result still gets emitted, leaving no unseen Block for the managed engine.
	fromBlockResult, ok := pollResults[fromBlock]
	if ok {
		if fromBlockResult.forked && tracker != nil {
			return result, errors.Wrapf(
				superwatcher.ErrFromBlockReorged, "fromBlock %d was removed/reorged", fromBlock,
			)
		}
	}

	return result, nil
}
