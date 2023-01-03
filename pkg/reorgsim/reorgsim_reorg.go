package reorgsim

import "go.uber.org/zap"

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

	debug := func(s string) {
		r.debugger.Debug(
			1, s,
			zap.Uint64("rangeStart", rangeStart),
			zap.Uint64("rangeEnd", rangeEnd),
			zap.Uint64("triggerBlock", event.ReorgTrigger),
			zap.Uint64("reorgBlock", event.ReorgBlock),
			zap.Uint64("currentBlock", r.currentBlock),
			zap.Int("currentReorgEvent", r.currentReorgEvent),
			zap.Int("triggered", r.triggered),
			zap.Int("forked", r.forked),
		)
	}

	// If event.ReorgBlock is within range, then trigger
	if rangeStart <= event.ReorgTrigger && rangeEnd >= event.ReorgTrigger {
		debug("triggering event")

		// First trigger on ReorgTrigger will not call forkChain
		if r.triggered < r.currentReorgEvent {
			r.triggered = r.currentReorgEvent

			debug("triggered, will NOT call forkChain now")
			return
		}

		// See if there's unforked r.reorgedChains before forking
		if r.triggered > r.forked && r.currentReorgEvent > r.forked {
			debug("triggered, will call forkChain now")
			r.forkChain()
		}
	}
}

// forkChain performs chain reorg logic after triggerForkChain validates all conditions
func (r *ReorgSim) forkChain() {
	event := r.events[r.currentReorgEvent]

	r.currentBlock = event.ReorgBlock
	r.forked = r.triggered
	r.chain = r.reorgedChains[r.currentReorgEvent]

	r.currentReorgEvent++

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
