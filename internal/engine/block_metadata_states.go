package engine

import (
	"fmt"
)

// Package `engine` uses a simple state machine to track each block's state.
// See STATES.md for explanation of the design decision.

// A block's initial (default) state `blockState` is always `stateNull`,
// and we can mutate the block's state by firing blockEvent event on the state.
// Note that each blockState corresponds the a block's hash (`see metadataTracker`).

type (
	blockState uint8 // blockState is used by Engine to determine if it should pass a block's logs to ServiceEngine.
	blockEvent uint8 // blockEvent is used by Engine to mutate blockState according to the state machine.
)

// All states and events are defined in the same const block to avoid collision.
const (
	stateNull         blockState = iota // Block was never seen before by Engine (default blockState)
	stateSeen                           // Block was seen by Engine
	stateHandled                        // Block was processed by ServiceEngine
	stateReorged                        // Block was present in a PollResult.ReorgedBlocks
	stateHandledReorg                   // Block's reorg was handled by ServiceEngine
	stateInvalid                        // Invalid blockState - program will panic

	eventInvalid     blockEvent = iota // Invalid blockEvent - program will panic (default blockEvent)
	eventSeeBlock                      // When Engine sees a block
	eventHandle                        // When ServiceEngine has processed the block's logs
	eventSeeReorg                      // When Engine sees the block in PollResult.ReorgedBlocks
	eventHandleReorg                   // When ServiceEngine has handled the reorg event
)

type stateEvent = struct {
	state blockState
	event blockEvent
}

// TODO: Update this table should we implement new feature "soft-errors".
// The soft errors features would allow ServiceEngine to inject "soft errors",
// which would NOT cause the engine to exit on seeing these soft errors.
// To implement this feature, we must allow certain states to be handled again,
// since some blocks with state stateSeen might fail and never progressed to stateHandled.
// i.e. as in stateSeen + eventSeeBlock = stateSeen, to allow the engine to re-handle the blocks.
// The transitions that need to be updated are tagged with TODO comments.
var watcherEngineStateMachine = map[stateEvent]blockState{
	{state: stateNull, event: eventSeeBlock}:    stateSeen,
	{state: stateNull, event: eventSeeReorg}:    stateInvalid,
	{state: stateNull, event: eventHandle}:      stateInvalid,
	{state: stateNull, event: eventHandleReorg}: stateInvalid,

	{state: stateSeen, event: eventSeeBlock}:    stateInvalid, // TODO: Change to stateSeen if implementing soft-errors
	{state: stateSeen, event: eventSeeReorg}:    stateInvalid, // TODO: Change to stateReorged if implementing soft-errors
	{state: stateSeen, event: eventHandle}:      stateHandled,
	{state: stateSeen, event: eventHandleReorg}: stateInvalid,

	{state: stateHandled, event: eventSeeBlock}:    stateHandled,
	{state: stateHandled, event: eventSeeReorg}:    stateReorged,
	{state: stateHandled, event: eventHandle}:      stateInvalid,
	{state: stateHandled, event: eventHandleReorg}: stateInvalid,

	{state: stateReorged, event: eventSeeBlock}:    stateInvalid,
	{state: stateReorged, event: eventSeeReorg}:    stateInvalid, // TODO: Change to stateReorged if implementing soft-errors
	{state: stateReorged, event: eventHandle}:      stateInvalid,
	{state: stateReorged, event: eventHandleReorg}: stateHandledReorg,

	{state: stateHandledReorg, event: eventSeeBlock}:    stateInvalid,
	{state: stateHandledReorg, event: eventSeeReorg}:    stateInvalid,
	{state: stateHandledReorg, event: eventHandle}:      stateInvalid,
	{state: stateHandledReorg, event: eventHandleReorg}: stateInvalid,
}

func (state *blockState) Fire(event blockEvent) {
	if !event.IsValid() {
		panic(fmt.Sprintf("invalid Engine event: %d", event))
	}

	this := stateEvent{state: *state, event: event}

	newState, ok := watcherEngineStateMachine[this]
	if !ok {
		*state = stateInvalid
		return
	}

	*state = newState
}

func (state blockState) String() string {
	switch state {
	case stateNull:
		return "NULL"
	case stateSeen:
		return "SEEN"
	case stateHandled:
		return "HANDLED"
	case stateReorged:
		return "REORGED"
	case stateHandledReorg:
		return "HANDLED_REORG"
	case stateInvalid:
		return "INVALID_BLOCK_STATE"
	}

	panic(fmt.Sprintf("invalid Engine state: %d", state))
}

func (state blockState) IsValid() bool {
	switch state {
	case stateInvalid:
		return false
	case
		stateNull,
		stateSeen,
		stateHandled,
		stateReorged,
		stateHandledReorg:
		return true
	}

	panic(fmt.Sprintf("invalid Engine state: %d", state))
}

func (event blockEvent) String() string {
	switch event {
	case eventSeeBlock:
		return "See Block"
	case eventHandle:
		return "Handle Block"
	case eventSeeReorg:
		return "See Reorg"
	case eventHandleReorg:
		return "Handle Reorg"
	}

	panic(fmt.Sprintf("invalid Engine event: %d", event))
}

func (event blockEvent) IsValid() bool {
	switch event {
	case
		eventSeeBlock,
		eventHandle,
		eventSeeReorg,
		eventHandleReorg:
		return true
	}

	return false
}
