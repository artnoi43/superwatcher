package emitter

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/lib/logger"
)

func (e *emitter) loopFilterLogs(ctx context.Context) error {
	sleep := func() {
		time.Sleep(time.Second * time.Duration(e.config.LoopInterval))
	}
	loopCtx := context.Background()
	// Assume that this is a normal first start (watcher restarted).
	lookBackFirstStart := true

filterLoop:
	for {
		logger.Info("starting filterLoop")
		if !lookBackFirstStart {
			sleep()
		}
		select {
		// Graceful shutdown in main
		case <-ctx.Done():
			return ctx.Err()
		default:
			currentBlockNumber, err := e.client.BlockNumber(loopCtx)
			if err != nil {
				return errors.Wrap(err, "failed to get current block number from node")
			}
			lastRecordedBlock, err := e.stateDataGateway.GetLastRecordedBlock(ctx)
			if err != nil {
				// Return error if not datagateway.ErrRecordNotFound
				if !errors.Is(err, datagateway.ErrRecordNotFound) {
					return errors.Wrap(err, "failed to get last recorded block from Redis")
				}
				// If no lastRecordedBlock => watcher has never been run on the host: there's no need to look back.
				lookBackFirstStart = false
				// If no lastRecordedBlock, use startBlock (contract genesis block)
				lastRecordedBlock = e.startBlock
			}

			logger.Info(
				"recent blocks",
				zap.Uint64("currentBlock", currentBlockNumber),
				zap.Uint64("lastRecordedBlock", lastRecordedBlock),
			)

			// Handle if there's no new block yet
			if lastRecordedBlock == currentBlockNumber {
				logger.Info(
					"no new block, skipping",
					zap.Uint64("currentBlock", currentBlockNumber),
					zap.Uint64("lastRecordedBlock", lastRecordedBlock),
				)
				continue filterLoop
			}

			var fromBlock, toBlock uint64
			if lookBackFirstStart {
				// Toggle
				lookBackFirstStart = false
				// Start with going back lookBack * maxLookBack times if watcher was restarted
				goBack := e.config.LookBackBlocks * e.config.LookBackRetries
				fromBlock = lastRecordedBlock + 1 - goBack
				toBlock = fromBlock + e.config.LookBackBlocks
				logger.Info(
					"first watcher run, going back",
					zap.Uint64("lastRecordedBlock", lastRecordedBlock),
					zap.Uint64("goBack", goBack),
					zap.Uint64("fromBlock", fromBlock),
				)
			} else {
				fromBlock, toBlock = fromBlockToBlock(currentBlockNumber, lastRecordedBlock, e.config.LookBackBlocks)
			}

			logger.Info(
				"calling filterLogs",
				zap.Uint64("fromBlock", fromBlock),
				zap.Uint64("toBlock", toBlock),
			)
			if err := e.filterLogs(
				loopCtx,
				fromBlock,
				toBlock,
			); err != nil {
				if errors.Is(err, errFromBlockReorged) {
					// Continue to filter from fromBlock
					logger.Info("fromBlock reorged", zap.String("error", err.Error()))
					continue
				}
				return errors.Wrap(err, "filterLogs error")
			}
		}
	}
}
