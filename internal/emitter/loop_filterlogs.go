package emitter

import (
	"context"
	"time"

	"github.com/artnoi43/gsl/gslutils"
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

	GoBackFirstStart bool   `json:"goBackFirstStart"`
	IsReorging       bool   `json:"isReorging"`
	RetriesCount     uint64 `json:"goBackRetries"` // Tracks how many times the emitter has to goBack because fromBlock was reorged
}

func (e *emitter) sleep() {
	time.Sleep(time.Second * time.Duration(e.conf.LoopInterval))
}

// loopFilterLogs is the emitter's main loop. It dynamically computes fromBlock and toBlock for `*emitter.filterLogs`,
// and will only returns if (some) errors happened.
func (e *emitter) loopFilterLogs(
	ctx context.Context,
	status *filterLogStatus,
) error {
	// Assume that this is a normal first start (watcher restarted).
	status.GoBackFirstStart = true

	// This will keep track of fromBlock/toBlock, as well as reorg status
	loopCtx := context.Background()
	for {
		// Don't sleep or log status on first loop
		if !status.GoBackFirstStart {
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
				status,
			)

			// If first run and there's no new block
			if status.GoBackFirstStart && errors.Is(err, errNoNewBlock) {
				err = nil
			}

			if err != nil {
				// Skip if there's no new block just yet
				if errors.Is(err, errNoNewBlock) {
					// Use zap.String because this is not actually an error
					e.debugger.Debug(
						1, "skipping", zap.String("reason", err.Error()),
						zap.Uint64("currentBlock", newStatus.CurrentBlock),
						zap.Uint64("lastRecordedBlock", newStatus.LastRecordedBlock),
					)
					continue
				}

				if errors.Is(err, errFetchError) {
					e.debugger.Debug(1, "fetch error", zap.Error(err))
					continue
				}

				return errors.Wrap(err, "emitter failed to compute fromBlock and toBlock")
			}

			// updateStatus is called after e.filterLogs returned. It updates status to newStatus,
			// and also increments the counter for tracking retries during reorg.
			// If |isReorging| is true, then the emitter *goes back* until the chain stops reorging.
			updateStatus := func(isReorging bool) {
				status = newStatus
				status.IsReorging = isReorging

				if isReorging {
					status.RetriesCount++
				} else {
					// Reset counter
					status.RetriesCount = 0
				}
			}

			e.debugger.Debug(
				2, "calling filterLogs",
				zap.Any("current_status", newStatus),
			)

			if err := e.filterLogs(
				loopCtx,
				newStatus.FromBlock,
				newStatus.ToBlock,
			); err != nil {
				if errors.Is(err, errFromBlockReorged) {
					// Continue to filter from fromBlock
					updateStatus(true)

					logger.Warn("fromBlock reorged", zap.Any("emitterStatus", newStatus))
					continue
				}

				if errors.Is(err, errProcessReorg) {
					e.debugger.Debug(
						1, "got errProcessReorg - contact prem@cleverse.com for reporting this bug",
						zap.Error(err),
					)
				}

				return errors.Wrap(err, "unexpected filterLogs error")
			}

			updateStatus(false)

			e.debugger.Debug(
				1, "filterLogs returned successfully",
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
	prevStatus *filterLogStatus,
) (
	*filterLogStatus,
	error,
) {
	// Get chain's tallest block number and compare it with lastRecordedBlock
	var err error
	currentBlock, err := e.client.BlockNumber(ctx)
	if err != nil {
		return prevStatus, errors.Wrap(err, "failed to get current block number from node")
	}

	// lastRecordedBlock was saved by engine.
	// The value to be saved should be superwatcher.FilterResult.LastGoodBlock
	lastRecordedBlock, err := e.stateDataGateway.GetLastRecordedBlock(ctx)
	if err != nil {
		// Return error if not datagateway.ErrRecordNotFound
		if !errors.Is(err, datagateway.ErrRecordNotFound) {
			return prevStatus, errors.Wrap(err, "failed to get last recorded block from Redis")
		}

		// If no lastRecordedBlock, then it means the emitter has never been run on the host: there's no need to look back.
		prevStatus.GoBackFirstStart = false
		// If no lastRecordedBlock, use startBlock (contract genesis block)
		lastRecordedBlock = e.conf.StartBlock
	}

	e.debugger.Debug(
		2,
		"recent blocks",
		zap.Uint64("currentChainBlock", currentBlock),
		zap.Uint64("lastRecordedBlock", lastRecordedBlock),
	)

	// Continue if there's no new block yet
	if lastRecordedBlock == currentBlock {
		if !prevStatus.GoBackFirstStart {
			return prevStatus, errors.Wrapf(errNoNewBlock, "block %d", currentBlock)
		}
	}

	// Update prevStatus with current states
	prevStatus.CurrentBlock = currentBlock
	prevStatus.LastRecordedBlock = lastRecordedBlock

	// And send the updated states to computeFromBlockToBlock
	fromBlock, toBlock, err := computeFromBlockToBlock(
		prevStatus,
		currentBlock,
		lastRecordedBlock,
		e.conf.FilterRange,
		e.conf.GoBackRetries,
		e.conf.StartBlock,
		e.debugger,
	)
	if err != nil {
		return prevStatus, errors.Wrap(err, "failed to compute fromBlockToBlock")
	}

	// If fromBlock was too far back, i.e. fromBlock = 0
	if fromBlock < e.conf.StartBlock && toBlock >= e.conf.StartBlock {
		fromBlock = e.conf.StartBlock
	}

	return &filterLogStatus{
		FromBlock:         fromBlock,
		ToBlock:           toBlock,
		CurrentBlock:      currentBlock,
		LastRecordedBlock: lastRecordedBlock,
		IsReorging:        prevStatus.IsReorging,
		RetriesCount:      prevStatus.RetriesCount,
	}, nil
}

// computeFromBlockToBlock defines a more higher-level computation logic for getting fromBlock and toBlock.
// It takes into account the emitter's config, the emitter status, chain status, etc.
// There are 3 possible cases when computing fromBlock and toBlock:
// (1) `goBackFirstStart` - this is when this function call is called when the emitter just started.
// The emitter will "go back" to `lastRecordedBlock - (filterRange * maxRetries)`
// (2) `status.IsReorging` - this is when emitter.filterLogs detected that the previous fromBlock was reorged.
// The emitter will "go back" to `previous fromBlock - filterRange`. This logic continues until fromBlock is no longer reorging.
// (3) Normal cases - this should be the base case for computing fromBlock and toBlock.
// The emitter will use `fromBlockToBlockNormal` function to compute the values.
func computeFromBlockToBlock(
	prevStatus *filterLogStatus,
	currentBlock uint64,
	lastRecordedBlock uint64,
	filterRange uint64,
	maxRetries uint64,
	startBlock uint64,
	debugger *debugger.Debugger,
) (
	uint64,
	uint64,
	error,
) {
	var fromBlock, toBlock uint64

	// Special case
	if prevStatus.GoBackFirstStart {
		// Start with going back for filterRange * goBackRetries blocks if watcher was restarted
		// lastRecordedBlock = 80, filterRange = 10, maxRetries = 5
		// 1st run: from(31), to(40)   [goBackFirstStart] -> lastRecordedBlock = 40 (goBack: range = 5)
		// 2nd run: from(31), to(50)   [normalCase]       -> lastRecordedBlock = 50 (range = 10)
		// 3rd run: from(41), to(60)   [normalCase]       -> lastRecordedBlock = 60 (range = 10)

		prevStatus.GoBackFirstStart = false
		goBack := filterRange * prevStatus.RetriesCount

		firstNewBlock := lastRecordedBlock + 1
		// Prevent overflow
		if goBack > firstNewBlock {
			debugger.Debug(
				1, "goBack > firstNewBlock",
				zap.Uint64("goBack", goBack),
				zap.Uint64("firstNewBlock", firstNewBlock),
			)

			fromBlock = startBlock
		} else {
			fromBlock = firstNewBlock - goBack
		}

		// The range is the same in firstStart
		toBlock = fromBlock + filterRange - 1

		debugger.Debug(
			1,
			"emitter: first run, going back",
			zap.Uint64("lastRecordedBlock", lastRecordedBlock),
			zap.Uint64("goBack", goBack),
			zap.Uint64("fromBlock", fromBlock),
		)
	} else if prevStatus.IsReorging {
		// The lookBack range will grow after each retries, but not the forward range
		// lastRecordedBlock = 80, filterRange = 10
		// 71 - 90   # none reorged in this loop       lastRecordedBlock = 90,  lookBack = 10, fwdRange = 90 - 80   = 10
		// 81 - 100  # none reoged in this loop        lastRecordedBlock = 100, lookBack = 10  fwdRange = 100 - 90  = 10
		// 91 - 110  # 91 reorged in this loop         lastRecordedBlock = 110, lookBack = 10  fwdRange = 110 - 100 = 10
		// 81 - 110  # 81 reorged in this loop         lastRecordedBlock = 110, lookBack = 15  fwdRange = 110 - 110 = 0
		// 71 - 110  # 71 reorged in this loop         lastRecordedBlock = 110, lookBack = 20  fwdRange = 110 - 110 = 0
		// 61 - 110  # none reorged in this loop       lastRecordedBlock = 110, lookBack = 25, fwdRange = 110 - 110 = 0
		// 101 - 120 # normal case continues           lastRecordedBlock = 120, lookBack = 10, fwdRange = 120 - 110 = 10

		if prevStatus.RetriesCount > maxRetries {
			return prevStatus.FromBlock, prevStatus.ToBlock, errors.Wrapf(ErrMaxRetriesReached, "%d goBackRetries", prevStatus.RetriesCount)
		}

		// goBack in this case (reorging) is fixed
		goBack := filterRange

		// goBack is 1500, but (prev) fromBlock is 1050
		if goBack > prevStatus.FromBlock {
			debugger.Debug(
				1, "goBack > fromBlock",
				zap.Uint64("goBack", goBack),
				zap.Uint64("fromBlock", fromBlock),
				zap.Uint64("currentRetries", prevStatus.RetriesCount),
			)

			fromBlock = startBlock
		} else {
			// goBack from the last fromBlock
			fromBlock = prevStatus.FromBlock - goBack
		}

		// toBlock does not go back, so we don't update it,
		// unless currentBlock was shrunk during reorg too.
		toBlock = gslutils.Min(prevStatus.ToBlock, currentBlock)
	} else {
		// Call fromBlockToBlock in normal cases
		// lastRecordedBlock = 80, filterRange = 10
		// 71 - 90   # none reorged in this loop       lastRecordedBlock = 90,  lookBack = 10, fwdRange = 90 - 80    = 10
		// 81 - 100  # none reoged in this loop        lastRecordedBlock = 100, lookBack = 10  fwdRange = 100 - 90   = 10
		// 91 - 110  # none reoged in this loop        lastRecordedBlock = 110, lookBack = 10  fwdRange = 110 - 100  = 10
		// 101 - 120 # none reoged in this loop        lastRecordedBlock = 120, lookBack = 10  fwdRange = 120 - 110  = 10

		fromBlock, toBlock = fromBlockToBlockNormal(startBlock, currentBlock, lastRecordedBlock, filterRange)
	}

	// Update status here too
	prevStatus.FromBlock = fromBlock
	prevStatus.ToBlock = toBlock

	return fromBlock, toBlock, nil
}
