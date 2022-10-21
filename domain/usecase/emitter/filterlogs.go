package emitter

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/artnoi43/gsl/concurrent"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
)

// filterLogs filters Ethereum event logs from fromBlock to toBlock,
// and sends *types.Log and *reorg.BlockInfo through w.logChan and w.reorgChan respectively.
// If an error is encountered, filterLogs returns with error.
// filterLogs should not be the one sending the error through w.errChan.
func (e *emitter) filterLogs(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) error {
	var wg sync.WaitGroup
	var mut sync.Mutex
	var eventLogs []types.Log
	var err error

	headersByBlockNumber := make(map[uint64]*types.Header)
	getErrChan := make(chan error)

	// getLogs calls FilterLogs from fromBlock to toBlock
	getLogs := func() {
		eventLogs, err = e.client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(toBlock)),
			Addresses: e.addresses,
			Topics:    e.topics,
		})
		if err != nil {
			getErrChan <- errors.Wrap(err, "error filtering event logs")
		}
	}

	// getHeader gets block header for a blockNumber
	getHeader := func(blockNumber uint64) {
		err := retry.Do(func() error {
			header, err := e.client.HeaderByNumber(ctx, big.NewInt(int64(blockNumber)))
			if err != nil {
				return err
			}
			mut.Lock()
			headersByBlockNumber[blockNumber] = header
			mut.Unlock()
			return nil
		},
			retry.Attempts(10),
			retry.Delay(300*time.Millisecond),
			retry.DelayType(retry.FixedDelay),
		)
		if err != nil {
			getErrChan <- errors.Wrapf(err, "failed to get header for block %d", blockNumber)
		}
	}

	// Get fresh logs, and block headers (fromBlock-toBlock)
	// to compare the headers with that of w.tracker's to detect chain reorg

	wg.Add(1)
	go func() {
		defer wg.Done()
		getLogs()
	}()

	for i := fromBlock; i <= toBlock; i++ {
		wg.Add(1)
		go func(blockNumber uint64) {
			defer wg.Done()
			getHeader(blockNumber)
		}(i)
	}

	// Wait here for logs and headers
	if err := concurrent.WaitAndCollectErrors(&wg, getErrChan); err != nil {
		logger.Error("get fresh data from blockchain failed", zap.String("error", err.Error()))
		return errors.Wrap(err, "get blockchain data")
	}

	lenLogs := len(eventLogs)
	logger.Info("got event logs", zap.Int("number of filtered logs", lenLogs))
	logger.Info("got headers and logs", zap.Uint64("fromBlock", fromBlock), zap.Uint64("toBlock", toBlock))

	// Clear all tracker's blocks before fromBlock - lookBackBlocks
	logger.Info("clearing tracker", zap.Uint64("clearUntil", fromBlock-e.config.LookBackBlocks))
	e.tracker.ClearUntil(fromBlock - e.config.LookBackBlocks)

	/* Use code from reorg package to manage/handle chain reorg */
	// Use fresh hashes and fresh logs to populate these 3 maps
	freshHashesByBlockNumber, freshLogsByBlockNumber, processLogsByBlockNumber := reorg.PopulateInitialMaps(eventLogs, headersByBlockNumber)
	// wasReorged maps block numbers whose fresh hash and tracker hash differ
	wasReorged := reorg.ProcessReorged(
		e.tracker,
		fromBlock,
		toBlock,
		freshHashesByBlockNumber,
		freshLogsByBlockNumber,
		processLogsByBlockNumber,
	)

	e.debugMsg("wasReorged", zap.Any("wasReorged", wasReorged))

	// If fromBlock was reorged, then return to loopFilterLogs
	if wasReorged[fromBlock] {
		return errors.Wrapf(errFromBlockReorged, "fromBlock %d was removed (chain reorganization)", fromBlock)
	}

	// Publish log(s) and reorged block, and add canon block to tracker
	for blockNumber := fromBlock; blockNumber <= toBlock; blockNumber++ {
		if wasReorged[blockNumber] {
			logger.Info(
				"chain reorg detected",
				zap.Uint64("blockNumber", blockNumber),
				zap.String("freshHash", freshHashesByBlockNumber[blockNumber].String()),
			)

			// For debugging
			reorgBlock, foundInTracker := e.tracker.GetTrackerBlockInfo(blockNumber)
			if !foundInTracker {
				logger.Panic(
					"blockInfo marked as reorged but was not found in tracker",
					zap.Uint64("blockNumber", blockNumber),
					zap.String("freshHash", reorgBlock.String()),
				)
			}

			e.publishReorg(reorgBlock)
			continue
		}

		// Populate blockInfo with fresh info
		b := reorg.NewBlockInfo(blockNumber, freshHashesByBlockNumber[blockNumber])
		b.Logs = freshLogsByBlockNumber[blockNumber]

		// Publish block with > 0 block
		if len(b.Logs) > 0 {
			e.publishBlock(b)
		}
		// Add ONLY CANONICAL block into tracker
		e.tracker.AddTrackerBlock(b)
	}

	// End loop
	logger.Info(
		"number of logs processed by filterLogs",
		zap.Int("eventLogs (filtered)", lenLogs),
		zap.Int("processLogs (all logs processed)", len(processLogsByBlockNumber)),
	)
	if err := e.stateDataGateway.SetLastRecordedBlock(ctx, toBlock); err != nil {
		return errors.Wrap(err, "failed to save lastRecordedBlock to redis")
	}
	logger.Info("set lastRecordedBlock", zap.Uint64("blockNumber", toBlock))

	return nil
}
