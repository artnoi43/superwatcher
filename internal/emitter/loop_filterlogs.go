package emitter

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

type filterLogStatus struct {
	IsReorging bool   `json:"isReorging"`
	FromBlock  uint64 `json:"fromBlock"`
	ToBlock    uint64 `json:"toBlock"`
}

func (e *emitter[H]) loopFilterLogs(ctx context.Context) error {
	sleep := func() {
		time.Sleep(time.Second * time.Duration(e.config.LoopInterval))
	}
	loopCtx := context.Background()
	// Assume that this is a normal first start (watcher restarted).
	lookBackFirstStart := true

	// This will keep track of fromBlock/toBlock, as well as reorg status
	var status filterLogStatus

	for {
		// Don't sleep or log status on first loop
		if !lookBackFirstStart {
			e.debugMsg("new loopFilterLogs loop", zap.Any("current_status", status))
			sleep()
		}

		select {
		// Graceful shutdown in main
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Get chain's tallest block number and compare it with lastRecordedBlock
			currentBlockNumber, err := e.client.BlockNumber(loopCtx)
			if err != nil {
				return errors.Wrap(err, "failed to get current block number from node")
			}

			// lastRecordedBlock was saved by engine.
			// The saved value is block number of the last element of FilterResult.GoodBlocks
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

			e.debugMsg(
				"recent blocks",
				zap.Uint64("currentBlock", currentBlockNumber),
				zap.Uint64("lastRecordedBlock", lastRecordedBlock),
			)

			// Continue if there's no new block yet
			if lastRecordedBlock == currentBlockNumber {
				e.debugMsg(
					"no new block, skipping",
					zap.Uint64("currentBlock", currentBlockNumber),
				)
				continue
			}

			// If first run or status.IsReorging, then we go back to filter previous blocks.
			var fromBlock, toBlock uint64
			if lookBackFirstStart || status.IsReorging {
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

			toggleStatusIsReorging := func(isReorging bool) {
				status.IsReorging = isReorging
				status.FromBlock = fromBlock
				status.ToBlock = toBlock
			}

			e.debugMsg(
				"calling filterLogs",
				zap.Uint64("fromBlock", fromBlock),
				zap.Uint64("toBlock", toBlock),
				zap.Any("current_status", status),
			)

			if err := e.filterLogs(
				loopCtx,
				fromBlock,
				toBlock,
			); err != nil {
				if errors.Is(err, errFromBlockReorged) {
					// Continue to filter from fromBlock
					toggleStatusIsReorging(true)

					logger.Warn("fromBlock reorged", zap.Any("filterLogStatus", status))
					continue
				} else if errors.Is(err, errFetchError) {
					// TODO: decide this
					// Continue if client failed to get headers/logs
					logger.Warn("client failed to fetch", zap.Error(err))
					continue
				}

				return errors.Wrap(err, "unexpected filterLogs error")
			}

			toggleStatusIsReorging(false)

			e.debugMsg(
				"filterLogs returned successfully",
				zap.Any("current_status", status),
			)
		}
	}
}
