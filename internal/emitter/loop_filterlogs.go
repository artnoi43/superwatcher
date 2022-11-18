package emitter

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/logger/debug"
)

type filterLogStatus struct {
	IsReorging bool   `json:"isReorging"`
	FromBlock  uint64 `json:"fromBlock"`
	ToBlock    uint64 `json:"toBlock"`
}

// loopFilterLogs is the emitter's main loop. It dynamically computes fromBlock and toBlock for `*emitter.filterLogs`,
// and will only returns if (some) errors happened.
func (e *emitter) loopFilterLogs(ctx context.Context) error {
	sleep := func() {
		time.Sleep(time.Second * time.Duration(e.config.LoopInterval))
	}
	loopCtx := context.Background()
	// Assume that this is a normal first start (watcher restarted).
	lookBackFirstStart := true

	// This will keep track of fromBlock/toBlock, as well as reorg status
	status := new(filterLogStatus)
	for {
		// Don't sleep or log status on first loop
		if !lookBackFirstStart {
			e.debugMsg("new loopFilterLogs loop", zap.Any("emitterStatus", status))
			sleep()
		}

		select {
		// Graceful shutdown in main
		case <-ctx.Done():
			return ctx.Err()
		default:
			fromBlock, toBlock, err := e.computeFromBlockToBlock(loopCtx, &lookBackFirstStart, status)
			if err != nil {
				if errors.Is(errNoNewBlock, err) {
					// Use zap.String because this is not actually an error
					e.debugMsg("skipping", zap.String("reason", err.Error()))
					continue
				}

				return errors.Wrap(err, "emitter failed to compute fromBlock and toBlock")
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

					logger.Warn("fromBlock reorged", zap.Any("emitterStatus", status))
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
				zap.Any("emitterStatus", status),
			)
		}
	}
}

// computeFromBlockToBlock gets relevant block numbers using emitter's resources
// such as |e.client| and |e.stateDataGateway| before calling the other (non-method) `computeFromBlockToBlock`,
// which does not use emitter's resources.
func (e *emitter) computeFromBlockToBlock(
	ctx context.Context,
	lookBackFirstStart *bool,
	status *filterLogStatus,
) (
	fromBlock uint64,
	toBlock uint64,
	err error,
) {
	fromBlock, toBlock = status.FromBlock, status.ToBlock
	// Get chain's tallest block number and compare it with lastRecordedBlock
	currentBlock, err := e.client.BlockNumber(ctx)
	if err != nil {
		return status.FromBlock, status.ToBlock, errors.Wrap(err, "failed to get current block number from node")
	}

	// lastRecordedBlock was saved by engine.
	// The value to be saved should be superwatcher.FilterResult.LastGoodBlock
	lastRecordedBlock, err := e.stateDataGateway.GetLastRecordedBlock(ctx)
	if err != nil {
		// Return error if not datagateway.ErrRecordNotFound
		if !errors.Is(err, datagateway.ErrRecordNotFound) {
			return fromBlock, toBlock, errors.Wrap(err, "failed to get last recorded block from Redis")
		}

		// If no lastRecordedBlock => watcher has never been run on the host: there's no need to look back.
		*lookBackFirstStart = false
		// If no lastRecordedBlock, use startBlock (contract genesis block)
		lastRecordedBlock = e.startBlock
	}

	e.debugMsg(
		"recent blocks",
		zap.Uint64("currentBlock", currentBlock),
		zap.Uint64("lastRecordedBlock", lastRecordedBlock),
	)

	// Continue if there's no new block yet
	if lastRecordedBlock == currentBlock {
		e.debugMsg(
			"no new block, skipping",
			zap.Uint64("currentBlock", currentBlock),
		)
		return fromBlock, toBlock, errors.Wrapf(errNoNewBlock, "block %d", currentBlock)
	}

	fromBlock, toBlock = computeFromBlockToBlock(e.config, lastRecordedBlock, currentBlock, lookBackFirstStart, status, e.debug)
	return fromBlock, toBlock, nil
}

// computeFromBlockToBlock defines a more higher-level computation logic for getting fromBlock and toBlock.
// It takes into account the emitter's config, the emitter status, chain status, etc.
// TODO: Test this func
func computeFromBlockToBlock(
	conf *config.Config,
	lastRecordedBlock, currentBlock uint64,
	lookBackFirstStart *bool,
	status *filterLogStatus,
	isDebug bool,
) (
	fromBlock uint64,
	toBlock uint64,
) {
	// Special case
	if *lookBackFirstStart || status.IsReorging {
		// Toggle
		*lookBackFirstStart = false

		// Start with going back lookBack * maxLookBack times if watcher was restarted
		goBack := conf.LookBackBlocks * conf.LookBackRetries
		fromBlock = lastRecordedBlock + 1 - goBack
		toBlock = fromBlock + conf.LookBackBlocks

		debug.DebugMsg(
			isDebug,
			"first watcher run, going back",
			zap.Uint64("lastRecordedBlock", lastRecordedBlock),
			zap.Uint64("goBack", goBack),
			zap.Uint64("fromBlock", fromBlock),
		)

	} else {
		fromBlock, toBlock = fromBlockToBlock(currentBlock, lastRecordedBlock, conf.LookBackBlocks)
	}

	return fromBlock, toBlock
}
