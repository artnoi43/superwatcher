package engine

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
)

type (
	EngineLogState uint8
	EngineLogEvent uint8
)

const (
	EngineLogStateNull EngineLogState = iota
	EngineLogStateSeen
	EngineLogStateProcessed
	EngineLogStateReorged
	EngineLogStateReorgHandled
	EngineLogStateInvalid

	EngineLogEventInvalid EngineLogEvent = iota
	EngineLogEventGotLog
	EngineLogEventProcess
	EngineLogEventReorg
	EngineLogEventHandleReorg
)

type stateEvent = struct {
	state EngineLogState
	event EngineLogEvent
}

var engineStateTransitionTable = map[stateEvent]EngineLogState{
	{state: EngineLogStateNull, event: EngineLogEventGotLog}:      EngineLogStateSeen,
	{state: EngineLogStateNull, event: EngineLogEventProcess}:     EngineLogStateInvalid,
	{state: EngineLogStateNull, event: EngineLogEventReorg}:       EngineLogStateReorged,
	{state: EngineLogStateNull, event: EngineLogEventHandleReorg}: EngineLogStateInvalid,

	{state: EngineLogStateSeen, event: EngineLogEventGotLog}:      EngineLogStateSeen,
	{state: EngineLogStateSeen, event: EngineLogEventProcess}:     EngineLogStateProcessed,
	{state: EngineLogStateSeen, event: EngineLogEventReorg}:       EngineLogStateReorged,
	{state: EngineLogStateSeen, event: EngineLogEventHandleReorg}: EngineLogStateReorgHandled,

	{state: EngineLogStateProcessed, event: EngineLogEventGotLog}:      EngineLogStateProcessed,
	{state: EngineLogStateProcessed, event: EngineLogEventProcess}:     EngineLogStateProcessed,
	{state: EngineLogStateProcessed, event: EngineLogEventReorg}:       EngineLogStateReorged,
	{state: EngineLogStateProcessed, event: EngineLogEventHandleReorg}: EngineLogStateReorgHandled,

	{state: EngineLogStateReorged, event: EngineLogEventGotLog}:      EngineLogStateReorged,
	{state: EngineLogStateReorged, event: EngineLogEventProcess}:     EngineLogStateInvalid,
	{state: EngineLogStateReorged, event: EngineLogEventReorg}:       EngineLogStateReorged,
	{state: EngineLogStateReorged, event: EngineLogEventHandleReorg}: EngineLogStateReorgHandled,

	{state: EngineLogStateReorgHandled, event: EngineLogEventGotLog}:      EngineLogStateInvalid,
	{state: EngineLogStateReorgHandled, event: EngineLogEventProcess}:     EngineLogStateInvalid,
	{state: EngineLogStateReorgHandled, event: EngineLogEventReorg}:       EngineLogStateReorged,
	{state: EngineLogStateReorgHandled, event: EngineLogEventHandleReorg}: EngineLogStateInvalid,
}

func (state *EngineLogState) Fire(event EngineLogEvent) {
	if !event.IsValid() {
		logger.Panic("invalid event", zap.String("event", event.String()))
	}
	self := stateEvent{state: *state, event: event}
	newState := engineStateTransitionTable[self]
	*state = newState
}

func (state EngineLogState) String() string {
	switch state {
	case EngineLogStateNull:
		return "NULL"
	case EngineLogStateSeen:
		return "SEEN"
	case EngineLogStateProcessed:
		return "PROCESSED"
	case EngineLogStateReorged:
		return "REORGED"
	case EngineLogStateReorgHandled:
		return "REORG_HANDLED"
	case EngineLogStateInvalid:
		return "INVALID_ENGINE_STATE"
	}

	panic(fmt.Sprintf("unexpected invalid state: %d", state))
}

func (state EngineLogState) IsValid() bool {
	switch state {
	case EngineLogStateInvalid:
		return false
	case
		EngineLogStateNull,
		EngineLogStateSeen,
		EngineLogStateProcessed,
		EngineLogStateReorged,
		EngineLogStateReorgHandled:
		return true
	}

	panic(fmt.Sprintf("unexpected invalid state: %d", state))
}

func (event EngineLogEvent) String() string {
	switch event {
	case EngineLogEventGotLog:
		return "Got Log"
	case EngineLogEventProcess:
		return "Process"
	case EngineLogEventReorg:
		return "Got Reorg"
	case EngineLogEventHandleReorg:
		return "Handle Reorg"
	}

	panic(fmt.Sprintf("unexpected invalid event: %d", event))
}

func (event EngineLogEvent) IsValid() bool {
	switch event {
	case EngineLogEventInvalid:
		return false
	case
		EngineLogEventGotLog,
		EngineLogEventProcess,
		EngineLogEventReorg,
		EngineLogEventHandleReorg:
		return true
	}

	panic(fmt.Sprintf("unexpected invalid event: %d", event))
}
