package watcher

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/watcher/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
)

func (w *watcher) filterLogs(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) error {
	headersByBlockNumber := make(map[uint64]*types.Header)
	errChan := make(chan error)
	var wg sync.WaitGroup
	var mut sync.Mutex
	var eventLogs []types.Log
	var err error

	// wait waits for wg and consume errors from errChanHeaders
	wait := func(errChan chan error) error {
		go func() {
			wg.Wait()
			close(errChan)
		}()
		if errChan != nil {
			var errs []string
			for err := range errChan {
				errs = append(errs, err.Error())
			}
			if len(errs) > 0 {
				return errors.New(strings.Join(errs, ","))
			}
		}
		return nil
	}

	// getLogs calls FilterLogs from fromBlock to toBlock
	getLogs := func() {
		eventLogs, err = w.client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(toBlock)),
			Addresses: w.addresses,
			Topics:    w.topics, // TODO: Topics not working yet
		})
		if err != nil {
			errChan <- errors.Wrap(err, "error filtering event logs")
		}
	}

	// getHeader gets block header for a blockNumber
	getHeader := func(blockNumber uint64) {
		err := retry.Do(func() error {
			header, err := w.client.HeaderByNumber(ctx, big.NewInt(int64(blockNumber)))
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
			w.errChan <- errors.Wrapf(err, "failed to get header for block %d", blockNumber)
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
	if err := wait(errChan); err != nil {
		logger.Error("get headers failed", zap.String("error", err.Error()))
		return errors.Wrap(err, "get headers")
	}

	lenLogs := len(eventLogs)
	logger.Info("got event logs", zap.Int("number of filtered logs", lenLogs))
	logger.Info("got headers and logs", zap.Uint64("fromBlock", fromBlock), zap.Uint64("toBlock", toBlock))

	// Clear all tracker's blocks before fromBlock - lookBackBlocks
	logger.Info("clearing tracker", zap.Uint64("clearUntil", fromBlock-w.config.LookBackBlocks))
	w.tracker.ClearUntil(fromBlock - w.config.LookBackBlocks)

	/*
		Detect chain reorg:
		[1st for-loop]: Populate processLogs from fresh event logs,
		and populate freshHashesByBlockNumber with fresh block hashes.

		[2nd for-loop]: Check if block hash saved in tracker matches the fresh block hash.
		If they are different, old logs from w.tracker will be tagged as Removed and
		PREPENDED in processLogs[blockNumber]

		Let's say we have these logs in the tracker:

		{block:68, hash:"0x68"}, {block: 69, hash:"0x69"}, {block:70, hash:"0x70"}

		And then we have these fresh logs:

		{block:68, hash:"0x68"}, {block: 69, hash:"0x112"}, {block:70, hash:"0x70"}

		The result processLogs will look like this map:
		{
			68: [{block:68, hash:"0x68"}]
			69: [{block: 69, hash:"0x69", removed: true}, {block: 69, hash:"0x112"}]
			70: [{block:70, hash:"0x70"}]
		}
	*/

	// Collect fresh per-block blockHashes into freshHashesByBlockNumber
	freshHashesByBlockNumber := make(map[uint64]common.Hash)
	// Fresh logs by BlockNumber (from var eventLogs)
	freshLogsByBlockNumber := make(map[uint64][]*types.Log)
	// Tracker + fresh logs by blockNumber - Watcher will process logs in this map
	processLogsByBlockNumber := make(map[uint64][]*types.Log)
	for i := range eventLogs {
		freshBlockNumber := eventLogs[i].BlockNumber
		freshBlockHash := eventLogs[i].BlockHash

		// Check if we saw the block hash
		if _, ok := freshHashesByBlockNumber[freshBlockNumber]; !ok {
			// We've never seen this block hash before
			freshHashesByBlockNumber[freshBlockNumber] = freshBlockHash
		} else {
			if freshHashesByBlockNumber[freshBlockNumber] != freshBlockHash {
				// If we saw this log's block hash before from the fresh eventLogs
				// from this loop's previous loop iteration(s), but the hashes are different:
				// Fatal, should not happen
				logger.Fatal("fresh blockHash differs from tracker blockHash",
					zap.String("tag", "filterLogs bug"),
					zap.Uint64("freshBlockNumber", freshBlockNumber),
					zap.Any("known tracker blockHash", freshHashesByBlockNumber[freshBlockNumber]),
					zap.Any("fresh blockHash", freshBlockHash))
			}
		}
		// Collect this log into freshLogs and processLogs
		thisLog := &eventLogs[i]
		freshLogsByBlockNumber[freshBlockNumber] = append(freshLogsByBlockNumber[freshBlockNumber], thisLog)
		processLogsByBlockNumber[freshBlockNumber] = append(processLogsByBlockNumber[freshBlockNumber], thisLog)
	}

	// noReorgs tracks block numbers whose fresh hash and tracker hash is identical
	noReorgs := make(map[uint64]bool)
	// Detect and stamp removed/reverted event logs using tracker
	for blockNumber := fromBlock; blockNumber <= toBlock; blockNumber++ {
		// If the block had not been saved into w.tracker (new blocks), it's probably fresh blocks,
		// which are not yet 'reorged' at the execution time.
		trackerBlock, foundInTracker := w.tracker.GetTrackerBlockInfo(blockNumber)
		if !foundInTracker {
			continue
		}

		// If tracker's is the same from recently filtered hash, i.e. no reorg
		// logger.Info("found block in tracker, comparing hashes in tracker", zap.Uint64("blockNumber", blockNumber))
		if h := freshHashesByBlockNumber[blockNumber]; h == trackerBlock.Hash {
			// Mark blockNumber with identical hash (no reorg)
			if len(freshLogsByBlockNumber[blockNumber]) == len(trackerBlock.Logs) {
				noReorgs[blockNumber] = true
				continue
			}
		}

		// Mark this block as reorged and call publishReorg
		logger.Info("chain reorg: detected", zap.Uint64("blockNumber", blockNumber), zap.String("freshHash", freshHashesByBlockNumber[blockNumber].String()), zap.String("trackerHash", trackerBlock.Hash.String()))
		w.publishReorg(trackerBlock)
		// Mark every log in this block as removed
		for _, oldLog := range trackerBlock.Logs {
			oldLog.Removed = true
		}
		// Concat logs from the same block, old logs first, into freshLogs
		processLogsByBlockNumber[blockNumber] = append(trackerBlock.Logs, processLogsByBlockNumber[blockNumber]...)
	}

	isFirstIdentical := false

	for blockNumber := fromBlock; blockNumber <= toBlock; blockNumber++ {
		// Skip blockNumber marked in identicals (no reorg)
		if noReorgs[blockNumber] && !isFirstIdentical {
			continue
		}
		isFirstIdentical = true

		// Populate blockInfo with fresh info
		b := reorg.NewBlockInfo(blockNumber, freshHashesByBlockNumber[blockNumber])
		b.Logs = freshLogsByBlockNumber[blockNumber]

		// Process every log for this block
		for _, l := range processLogsByBlockNumber[blockNumber] {
			// watcher.processLog will populate the log with decoded data.
			// After sucessful return, orderStates values will be the decoded block values
			w.publishLog(l)
		}

		// Add to tracker here
		w.tracker.AddTrackerBlock(b)
	}

	// End loop
	logger.Info(
		"number of logs processed by filterLogs",
		zap.Int("eventLogs (filtered)", lenLogs),
		zap.Int("processLogs (all logs processed)", len(processLogsByBlockNumber)),
	)
	if err := w.stateDataGateway.SetLastRecordedBlock(ctx, toBlock); err != nil {
		return errors.Wrap(err, "failed to save lastRecordedBlock to redis")
	}
	logger.Info("set lastRecordedBlock", zap.Uint64("blockNumber", toBlock))

	return nil
}
