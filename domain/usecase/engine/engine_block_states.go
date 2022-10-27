package engine

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
)

type (
	EngineBlockState uint8
	EngineBlockEvent uint8
)

const (
	EngineBlockStateNull EngineBlockState = iota
	EngineBlockStateSeen
	EngineBlockStateProcessed
	EngineBlockStateReorged
	EngineBlockStateReorgHandled
	EngineBlockStateInvalid

	EngineBlockEventInvalid EngineBlockEvent = iota
	EngineBlockEventGotLog
	EngineBlockEventProcess
	EngineBlockEventReorg
	EngineBlockEventHandleReorg
)

type stateEvent = struct {
	state EngineBlockState
	event EngineBlockEvent
}

var engineStateTransitionTable = map[stateEvent]EngineBlockState{
	{state: EngineBlockStateNull, event: EngineBlockEventGotLog}:      EngineBlockStateSeen,
	{state: EngineBlockStateNull, event: EngineBlockEventProcess}:     EngineBlockStateInvalid,
	{state: EngineBlockStateNull, event: EngineBlockEventReorg}:       EngineBlockStateReorged,
	{state: EngineBlockStateNull, event: EngineBlockEventHandleReorg}: EngineBlockStateInvalid,

	{state: EngineBlockStateSeen, event: EngineBlockEventGotLog}:      EngineBlockStateSeen,
	{state: EngineBlockStateSeen, event: EngineBlockEventProcess}:     EngineBlockStateProcessed,
	{state: EngineBlockStateSeen, event: EngineBlockEventReorg}:       EngineBlockStateReorged,
	{state: EngineBlockStateSeen, event: EngineBlockEventHandleReorg}: EngineBlockStateReorgHandled,

	{state: EngineBlockStateProcessed, event: EngineBlockEventGotLog}:      EngineBlockStateProcessed,
	{state: EngineBlockStateProcessed, event: EngineBlockEventProcess}:     EngineBlockStateProcessed,
	{state: EngineBlockStateProcessed, event: EngineBlockEventReorg}:       EngineBlockStateReorged,
	{state: EngineBlockStateProcessed, event: EngineBlockEventHandleReorg}: EngineBlockStateReorgHandled,

	{state: EngineBlockStateReorged, event: EngineBlockEventGotLog}:      EngineBlockStateReorged,
	{state: EngineBlockStateReorged, event: EngineBlockEventProcess}:     EngineBlockStateInvalid,
	{state: EngineBlockStateReorged, event: EngineBlockEventReorg}:       EngineBlockStateReorged,
	{state: EngineBlockStateReorged, event: EngineBlockEventHandleReorg}: EngineBlockStateReorgHandled,

	{state: EngineBlockStateReorgHandled, event: EngineBlockEventGotLog}:      EngineBlockStateInvalid,
	{state: EngineBlockStateReorgHandled, event: EngineBlockEventProcess}:     EngineBlockStateInvalid,
	{state: EngineBlockStateReorgHandled, event: EngineBlockEventReorg}:       EngineBlockStateReorged,
	{state: EngineBlockStateReorgHandled, event: EngineBlockEventHandleReorg}: EngineBlockStateInvalid,

	{state: EngineBlockStateInvalid, event: EngineBlockEventGotLog}:      EngineBlockStateInvalid,
	{state: EngineBlockStateInvalid, event: EngineBlockEventProcess}:     EngineBlockStateInvalid,
	{state: EngineBlockStateInvalid, event: EngineBlockEventReorg}:       EngineBlockStateInvalid,
	{state: EngineBlockStateInvalid, event: EngineBlockEventHandleReorg}: EngineBlockStateInvalid,
}

func (state *EngineBlockState) Fire(event EngineBlockEvent) {
	if !event.IsValid() {
		logger.Panic("invalid event", zap.String("event", event.String()))
	}

	self := stateEvent{state: *state, event: event}
	newState := engineStateTransitionTable[self]
	*state = newState
}

func (state EngineBlockState) String() string {
	switch state {
	case EngineBlockStateNull:
		return "NULL"
	case EngineBlockStateSeen:
		return "SEEN"
	case EngineBlockStateProcessed:
		return "PROCESSED"
	case EngineBlockStateReorged:
		return "REORGED"
	case EngineBlockStateReorgHandled:
		return "REORG_HANDLED"
	case EngineBlockStateInvalid:
		return "INVALID_ENGINE_STATE"
	}

	panic(fmt.Sprintf("unexpected invalid state: %d", state))
}

func (state EngineBlockState) IsValid() bool {
	switch state {
	case EngineBlockStateInvalid:
		return false
	case
		EngineBlockStateNull,
		EngineBlockStateSeen,
		EngineBlockStateProcessed,
		EngineBlockStateReorged,
		EngineBlockStateReorgHandled:
		return true
	}

	panic(fmt.Sprintf("unexpected invalid state: %d", state))
}

func (event EngineBlockEvent) String() string {
	switch event {
	case EngineBlockEventGotLog:
		return "Got Log"
	case EngineBlockEventProcess:
		return "Process"
	case EngineBlockEventReorg:
		return "Got Reorg"
	case EngineBlockEventHandleReorg:
		return "Handle Reorg"
	}

	panic(fmt.Sprintf("unexpected invalid event: %d", event))
}

func (event EngineBlockEvent) IsValid() bool {
	switch event {
	case EngineBlockEventInvalid:
		return false
	case
		EngineBlockEventGotLog,
		EngineBlockEventProcess,
		EngineBlockEventReorg,
		EngineBlockEventHandleReorg:
		return true
	}

	panic(fmt.Sprintf("unexpected invalid event: %d", event))
}
