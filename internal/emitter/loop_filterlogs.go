package emitter

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/logger/debug"
)

type filterLogStatus struct {
	IsReorging        bool   `json:"isReorging"`
	CurrentBlock      uint64 `json:"currentBlock"`
	LastRecordedBlock uint64 `json:"lastRecordedBlock"`
	FromBlock         uint64 `json:"fromBlock"`
	ToBlock           uint64 `json:"toBlock"`
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
			newStatus, err := e.computeFromBlockToBlock(
				loopCtx,
				&lookBackFirstStart,
				status,
			)
			if err != nil {
				if errors.Is(errNoNewBlock, err) {
					// Use zap.String because this is not actually an error
					e.debugMsg("skipping", zap.String("reason", err.Error()))
					continue
				}

				return errors.Wrap(err, "emitter failed to compute fromBlock and toBlock")
			}

			toggleStatusIsReorging := func(isReorging bool) {
				status = newStatus
				status.IsReorging = isReorging
			}

			e.debugMsg(
				"calling filterLogs",
				zap.Any("current_status", newStatus),
			)

			if err := e.filterLogs(
				loopCtx,
				newStatus.FromBlock,
				newStatus.ToBlock,
			); err != nil {
				if errors.Is(err, errFromBlockReorged) {
					// Continue to filter from fromBlock
					toggleStatusIsReorging(true)

					logger.Warn("fromBlock reorged", zap.Any("emitterStatus", newStatus))
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
				zap.Any("emitterStatus", newStatus),
			)
		}
	}
}

// computeFromBlockToBlock gets relevant block numbers using emitter's resources
// such as |e.client| and |e.stateDataGateway| before calling the other (non-method) `computeFromBlockToBlock`,
// which does not use emitter's resources.
// It returns old status with updated values if there was an error, or new status if successful.
func (e *emitter) computeFromBlockToBlock(
	ctx context.Context,
	lookBackFirstStart *bool,
	status *filterLogStatus,
) (
	*filterLogStatus,
	error,
) {
	// Get chain's tallest block number and compare it with lastRecordedBlock
	var err error
	currentBlock, err := e.client.BlockNumber(ctx)
	if err != nil {
		return status, errors.Wrap(err, "failed to get current block number from node")
	}
	status.CurrentBlock = currentBlock

	// lastRecordedBlock was saved by engine.
	// The value to be saved should be superwatcher.FilterResult.LastGoodBlock
	lastRecordedBlock, err := e.stateDataGateway.GetLastRecordedBlock(ctx)
	if err != nil {
		// Return error if not datagateway.ErrRecordNotFound
		if !errors.Is(err, datagateway.ErrRecordNotFound) {
			return status, errors.Wrap(err, "failed to get last recorded block from Redis")
		}

		// If no lastRecordedBlock => watcher has never been run on the host: there's no need to look back.
		*lookBackFirstStart = false
		// If no lastRecordedBlock, use startBlock (contract genesis block)
		lastRecordedBlock = e.startBlock
	}
	status.LastRecordedBlock = lastRecordedBlock

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
		return status, errors.Wrapf(errNoNewBlock, "block %d", currentBlock)
	}

	fromBlock, toBlock := computeFromBlockToBlock(
		currentBlock,
		lastRecordedBlock,
		e.config.LookBackBlocks,
		e.config.LookBackRetries,
		lookBackFirstStart,
		status,
		e.debug,
	)

	return &filterLogStatus{
		FromBlock:         fromBlock,
		ToBlock:           toBlock,
		CurrentBlock:      currentBlock,
		LastRecordedBlock: lastRecordedBlock,
		IsReorging:        status.IsReorging,
	}, nil
}

// computeFromBlockToBlock defines a more higher-level computation logic for getting fromBlock and toBlock.
// It takes into account the emitter's config, the emitter status, chain status, etc.
// TODO: Test this func
func computeFromBlockToBlock(
	currentBlock uint64,
	lastRecordedBlock uint64,
	lookBackBlocks uint64,
	lookBackRetries uint64,
	lookBackFirstStart *bool,
	status *filterLogStatus,
	isDebug bool,
) (
	uint64,
	uint64,
) {
	var fromBlock, toBlock uint64

	// Special case
	if *lookBackFirstStart || status.IsReorging {
		// Toggle
		*lookBackFirstStart = false

		// Start with going back lookBack * maxLookBack times if watcher was restarted
		goBack := lookBackBlocks * lookBackRetries
		base := lastRecordedBlock + 1

		// Prevent overflow
		if goBack > base {
			fromBlock = 1
		} else {
			fromBlock = base - goBack
		}
		toBlock = fromBlock + lookBackBlocks

		debug.DebugMsg(
			isDebug,
			"first watcher run, going back",
			zap.Uint64("lastRecordedBlock", lastRecordedBlock),
			zap.Uint64("goBack", goBack),
			zap.Uint64("fromBlock", fromBlock),
		)

	} else {
		fromBlock, toBlock = fromBlockToBlock(currentBlock, lastRecordedBlock, lookBackBlocks)
	}

	// Update status here too
	status.FromBlock = fromBlock
	status.ToBlock = toBlock

	return fromBlock, toBlock
}
