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
	blockState uint8 // blockState is used by WatcherEngine to determine if it should pass a block's logs to ServiceEngine.
	blockEvent uint8 // blockEvent is used by WatcherEngine to mutate blockState according to the state machine.
)

// All states and events are defined in the same const block to avoid collision.
const (
	stateNull         blockState = iota // Block was never seen before by WatcherEngine (default blockState)
	stateSeen                           // Block was seen by WatcherEngine
	stateHandled                        // Block was processed by ServiceEngine
	stateReorged                        // Block was present in a FilterResult.ReorgedBlocks
	stateHandledReorg                   // Block's reorg was handled by ServiceEngine
	stateInvalid                        // Invalid blockState - program will panic

	eventInvalid     blockEvent = iota // Invalid blockEvent - program will panic (default blockEvent)
	eventSeeBlock                      // When WatcherEngine sees a block
	eventHandle                        // When ServiceEngine has processed the block's logs
	eventSeeReorg                      // When WatcherEngine sees the block in FilterResult.ReorgedBlocks
	eventHandleReorg                   // When ServiceEngine has handled the reorg event
)

type stateEvent = struct {
	state blockState
	event blockEvent
}

var watcherEngineStateMachine = map[stateEvent]blockState{
	{state: stateNull, event: eventSeeBlock}:    stateSeen,
	{state: stateNull, event: eventSeeReorg}:    stateInvalid,
	{state: stateNull, event: eventHandle}:      stateInvalid,
	{state: stateNull, event: eventHandleReorg}: stateInvalid,

	{state: stateSeen, event: eventSeeBlock}:    stateSeen,    // Maybe stateInvalid is better?
	{state: stateSeen, event: eventSeeReorg}:    stateReorged, // Maybe stateInvalid is better?
	{state: stateSeen, event: eventHandle}:      stateHandled,
	{state: stateSeen, event: eventHandleReorg}: stateInvalid,

	{state: stateHandled, event: eventSeeBlock}:    stateHandled,
	{state: stateHandled, event: eventSeeReorg}:    stateReorged,
	{state: stateHandled, event: eventHandle}:      stateInvalid,
	{state: stateHandled, event: eventHandleReorg}: stateInvalid,

	{state: stateReorged, event: eventSeeBlock}:    stateInvalid,
	{state: stateReorged, event: eventSeeReorg}:    stateInvalid,
	{state: stateReorged, event: eventHandle}:      stateInvalid,
	{state: stateReorged, event: eventHandleReorg}: stateHandledReorg,

	{state: stateHandledReorg, event: eventSeeBlock}:    stateInvalid,
	{state: stateHandledReorg, event: eventSeeReorg}:    stateHandledReorg,
	{state: stateHandledReorg, event: eventHandle}:      stateInvalid,
	{state: stateHandledReorg, event: eventHandleReorg}: stateInvalid,

	{state: stateInvalid, event: eventSeeBlock}:    stateInvalid,
	{state: stateInvalid, event: eventSeeReorg}:    stateInvalid,
	{state: stateInvalid, event: eventHandle}:      stateInvalid,
	{state: stateInvalid, event: eventHandleReorg}: stateInvalid,
}

func (state *blockState) Fire(event blockEvent) {
	if !event.IsValid() {
		panic(fmt.Sprintf("invalid WatcherEngine event: %d", event))
	}

	self := stateEvent{state: *state, event: event}
	newState := watcherEngineStateMachine[self]
	*state = newState
}

func (state blockState) String() string {
	switch state {
	case stateNull:
		return "NULL"
	case stateSeen:
		return "SEEN"
	case stateHandled:
		return "PROCESSED"
	case stateReorged:
		return "REORGED"
	case stateHandledReorg:
		return "REORG_HANDLED"
	case stateInvalid:
		return "INVALID_ENGINE_STATE"
	}

	panic(fmt.Sprintf("invalid WatcherEngine state: %d", state))
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

	panic(fmt.Sprintf("invalid WatcherEngine state: %d", state))
}

func (event blockEvent) String() string {
	switch event {
	case eventSeeBlock:
		return "Got Log"
	case eventHandle:
		return "Process"
	case eventSeeReorg:
		return "Got Reorg"
	case eventHandleReorg:
		return "Handle Reorg"
	}

	panic(fmt.Sprintf("invalid WatcherEngine event: %d", event))
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
