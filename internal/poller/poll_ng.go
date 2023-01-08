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

	"github.com/artnoi43/superwatcher"
)

func poll( // nolint:unused
	ctx context.Context,
	addresses []common.Address,
	topics [][]common.Hash,
	fromBlock uint64,
	toBlock uint64,
	client superwatcher.EthClient,
	policy superwatcher.Policy,
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
	doReorg bool,
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

	if tracker != nil {
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

			mapResults[n].Header = header
			mapResults[n].Hash = headerHash

			if trackerBlock.Hash == pollResult.Hash && len(trackerBlock.Logs) == len(pollResult.Logs) {
				continue
			}

			mapResults[n].reorged = true
		}
	}

	return mapResults, nil
}
