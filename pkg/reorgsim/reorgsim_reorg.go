package reorgsim

import (
	"go.uber.org/zap"
)

// triggerForkChain needs 2 calls with event.ReorgTrigger within range.
// The 1st call updates `r.triggered` to r.currentReorgEvent
// if the current ReorgEvent.ReorgTrigger is within range [rangeStart, rangeEnd].
// The 2nd call will call forkChain if the current ReorgEvent.ReorgTrigger is within range.
func (r *ReorgSim) triggerForkChain(rangeStart, rangeEnd uint64) {
	// No need to trigger
	if l := len(r.events); l == 0 || r.currentReorgEvent >= l {
		return
	}

	event := r.events[r.currentReorgEvent]

	// If event.ReorgBlock is within range, then trigger
	if rangeStart <= event.ReorgTrigger && rangeEnd >= event.ReorgTrigger {
		r.debugger.Debug(
			1, "triggering event",
			zap.Uint64("rangeStart", rangeStart),
			zap.Uint64("rangeEnd", rangeEnd),
			zap.Uint64("triggerBlock", event.ReorgTrigger),
			zap.Uint64("reorgBlock", event.ReorgBlock),
			zap.Uint64("currentBlock", r.currentBlock),
			zap.Int("currentReorgEvent", r.currentReorgEvent),
			zap.Int("triggered", r.triggered),
			zap.Int("forked", r.forked),
		)

		// First trigger on ReorgTrigger will not call forkChain
		if r.triggered < r.currentReorgEvent {
			r.triggered = r.currentReorgEvent

			r.debugger.Debug(
				1, "triggered, will NOT call forkChain now",
				zap.Uint64("rangeStart", rangeStart),
				zap.Uint64("rangeEnd", rangeEnd),
				zap.Uint64("triggerBlock", event.ReorgTrigger),
				zap.Uint64("reorgBlock", event.ReorgBlock),
				zap.Uint64("currentBlock", r.currentBlock),
				zap.Int("currentReorgEvent", r.currentReorgEvent),
				zap.Int("triggered", r.triggered),
				zap.Int("forked", r.forked),
			)

			return
		}

		// See if there's unforked r.reorgedChains before forking
		if r.triggered > r.forked {
			r.debugger.Debug(
				1, "triggered, will call forkChain now",
				zap.Uint64("rangeStart", rangeStart),
				zap.Uint64("rangeEnd", rangeEnd),
				zap.Uint64("triggerBlock", event.ReorgTrigger),
				zap.Uint64("reorgBlock", event.ReorgBlock),
				zap.Uint64("currentBlock", r.currentBlock),
				zap.Int("currentReorgEvent", r.currentReorgEvent),
				zap.Int("triggered", r.triggered),
				zap.Int("forked", r.forked),
			)

			r.forkChain()
		}
	}
}

// forkChain performs chain reorg logic if the current ReorgEvent.ReorgBlock is within range [from, to]
// and if r.seen[ReorgEvent.ReorgBlock] is >= 1. The latter check allows for the poller/emitter to see
// pre-fork block hash once, so that we can test the poller/emitter logic.
func (r *ReorgSim) forkChain() {
	if r.triggered == r.forked {
		r.debugger.Debug(
			1, "trigger already forked, returning without forking",
			zap.Uint64("currentBlock", r.currentBlock),
			zap.Int("currentReorgEvent", r.currentReorgEvent),
			zap.Int("triggered", r.triggered),
			zap.Int("forked", r.forked),
		)

		return
	}

	event := r.events[r.currentReorgEvent]

	var currentChain BlockChain
	var lastReorg bool

	if l := len(r.events); r.currentReorgEvent < l {
		currentChain = r.reorgedChains[r.currentReorgEvent]
	} else {
		lastReorg = true
		currentChain = r.reorgedChains[l-1]
	}

	if !lastReorg {
		if r.forked < r.currentReorgEvent {
			r.chain = currentChain
			r.currentBlock = event.ReorgBlock
			r.forked = r.triggered

			r.currentReorgEvent++
		}
	}

	r.debugger.Debug(
		1, "REORGED! fork done",
		zap.Uint64("currentBlock", r.currentBlock),
		zap.Uint64("reorgTrigger", event.ReorgTrigger),
		zap.Uint64("reorgBlock", event.ReorgBlock),
		zap.Int("currentReorgEvent", r.currentReorgEvent),
		zap.Int("triggered", r.triggered),
		zap.Int("forked", r.forked),
	)
}
