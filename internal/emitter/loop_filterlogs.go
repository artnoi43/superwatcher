package emitter

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

type filterLogStatus struct {
	FromBlock         uint64 `json:"fromBlock"`
	ToBlock           uint64 `json:"toBlock"`
	CurrentBlock      uint64 `json:"currentBlock"`
	LastRecordedBlock uint64 `json:"lastRecordedBlock"`
	IsReorging        bool   `json:"isReorging"`
}

func (e *emitter) sleep() {
	time.Sleep(time.Second * time.Duration(e.conf.LoopInterval))
}

// loopFilterLogs is the emitter's main loop. It dynamically computes fromBlock and toBlock for `*emitter.filterLogs`,
// and will only returns if (some) errors happened.
func (e *emitter) loopFilterLogs(ctx context.Context, status *filterLogStatus) error {
	// Assume that this is a normal first start (watcher restarted).
	goBackFirstStart := true

	// This will keep track of fromBlock/toBlock, as well as reorg status
	loopCtx := context.Background()
	for {
		// Don't sleep or log status on first loop
		if !goBackFirstStart {
			e.debugger.Debug(1, "new loopFilterLogs loop")
			e.sleep()
		}

		select {
		// Graceful shutdown in main
		case <-ctx.Done():
			return ctx.Err()

		default:
			newStatus, err := e.computeFromBlockToBlock(
				loopCtx,
				&goBackFirstStart,
				status,
			)
			if goBackFirstStart && errors.Is(err, errNoNewBlock) {
				err = nil
			}

			if err != nil {
				if errors.Is(err, errNoNewBlock) {
					// Use zap.String because this is not actually an error
					e.debugger.Debug(1, "skipping", zap.String("reason", err.Error()))
					continue
				}
				if errors.Is(err, errFetchError) {
					e.debugger.Debug(1, "fetch error", zap.Error(err))
					continue
				}

				return errors.Wrap(err, "emitter failed to compute fromBlock and toBlock")
			}

			toggleStatusIsReorging := func(isReorging bool) {
				status = newStatus
				status.IsReorging = isReorging
			}

			e.debugger.Debug(
				2,
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

				}

				return errors.Wrap(err, "unexpected filterLogs error")
			}

			toggleStatusIsReorging(false)

			e.debugger.Debug(
				1,
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
	goBackFirstStart *bool,
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

		// If no lastRecordedBlock, then it means the emitter has never been run on the host: there's no need to look back.
		*goBackFirstStart = false
		// If no lastRecordedBlock, use startBlock (contract genesis block)
		lastRecordedBlock = e.conf.StartBlock
	}
	status.LastRecordedBlock = lastRecordedBlock

	e.debugger.Debug(
		1,
		"recent blocks",
		zap.Uint64("currentChainBlock", currentBlock),
		zap.Uint64("lastRecordedBlock", lastRecordedBlock),
	)

	// Continue if there's no new block yet
	if !*goBackFirstStart {
		if lastRecordedBlock == currentBlock {
			return status, errors.Wrapf(errNoNewBlock, "block %d", currentBlock)
		}
	}

	fromBlock, toBlock := computeFromBlockToBlock(
		currentBlock,
		lastRecordedBlock,
		e.conf.FilterRange,
		e.conf.GoBackRetries,
		goBackFirstStart,
		status,
		e.debugger,
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
	filterRange uint64,
	goBackRetries uint64,
	goBackFirstStart *bool,
	status *filterLogStatus,
	debugger *debugger.Debugger,
) (
	uint64,
	uint64,
) {
	var fromBlock, toBlock uint64

	// Special case
	if *goBackFirstStart || status.IsReorging {
		// Toggle
		*goBackFirstStart = false

		// TODO: Implement maxGoBack
		// Start with going back filterRange * goBackRetries times if watcher was restarted
		base := lastRecordedBlock + 1
		goBack := filterRange * goBackRetries

		// Prevent overflow
		if goBack > base {
			fromBlock = 1
		} else {
			fromBlock = base - goBack
		}
		toBlock = fromBlock + filterRange

		if debugger != nil {
			debugger.Debug(
				1,
				"emitter: first run, going back",
				zap.Uint64("lastRecordedBlock", lastRecordedBlock),
				zap.Uint64("goBack", goBack),
				zap.Uint64("fromBlock", fromBlock),
			)
		}

	} else {
		fromBlock, toBlock = fromBlockToBlock(currentBlock, lastRecordedBlock, filterRange)
	}

	// Update status here too
	status.FromBlock = fromBlock
	status.ToBlock = toBlock

	return fromBlock, toBlock
}
