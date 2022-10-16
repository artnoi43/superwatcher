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

type EngineFSM[K itemKey] interface {
	SetEngineState(K, EngineLogState)
	GetEngineState(K) EngineLogState
}

const (
	EngineStateNull EngineLogState = iota
	EngineStateSeen
	EngineStateProcessed
	EngineStateReorged
	EngineStateReorgHandled
	EngineStateInvalid

	EngineEventInvalid EngineLogEvent = iota
	EngineEventGotLog
	EngineEventProcess
	EngineEventReorg
	EngineEventHandleReorg
)

type stateEvent = struct {
	state EngineLogState
	event EngineLogEvent
}

var engineStateTransitionTable = map[stateEvent]EngineLogState{
	{state: EngineStateNull, event: EngineEventGotLog}:      EngineStateSeen,
	{state: EngineStateNull, event: EngineEventProcess}:     EngineStateInvalid,
	{state: EngineStateNull, event: EngineEventReorg}:       EngineStateReorged,
	{state: EngineStateNull, event: EngineEventHandleReorg}: EngineStateInvalid,

	{state: EngineStateSeen, event: EngineEventGotLog}:      EngineStateSeen,
	{state: EngineStateSeen, event: EngineEventProcess}:     EngineStateProcessed,
	{state: EngineStateSeen, event: EngineEventReorg}:       EngineStateReorged,
	{state: EngineStateSeen, event: EngineEventHandleReorg}: EngineStateReorgHandled,

	{state: EngineStateProcessed, event: EngineEventGotLog}:      EngineStateProcessed,
	{state: EngineStateProcessed, event: EngineEventProcess}:     EngineStateProcessed,
	{state: EngineStateProcessed, event: EngineEventReorg}:       EngineStateReorged,
	{state: EngineStateProcessed, event: EngineEventHandleReorg}: EngineStateReorgHandled,

	{state: EngineStateReorged, event: EngineEventGotLog}:      EngineStateReorged,
	{state: EngineStateReorged, event: EngineEventProcess}:     EngineStateInvalid,
	{state: EngineStateReorged, event: EngineEventReorg}:       EngineStateReorged,
	{state: EngineStateReorged, event: EngineEventHandleReorg}: EngineStateReorgHandled,

	{state: EngineStateReorgHandled, event: EngineEventGotLog}:      EngineStateInvalid,
	{state: EngineStateReorgHandled, event: EngineEventProcess}:     EngineStateInvalid,
	{state: EngineStateReorgHandled, event: EngineEventReorg}:       EngineStateReorged,
	{state: EngineStateReorgHandled, event: EngineEventHandleReorg}: EngineStateInvalid,
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
	case EngineStateNull:
		return "NULL"
	case EngineStateSeen:
		return "SEEN"
	case EngineStateProcessed:
		return "PROCESSED"
	case EngineStateReorged:
		return "REORGED"
	case EngineStateReorgHandled:
		return "REORG_HANDLED"
	case EngineStateInvalid:
		return "INVALID_ENGINE_STATE"
	}

	panic(fmt.Sprintf("unexpected invalid state: %d", state))
}

func (state EngineLogState) IsValid() bool {
	switch state {
	case EngineStateInvalid:
		return false
	case
		EngineStateNull,
		EngineStateSeen,
		EngineStateProcessed,
		EngineStateReorged,
		EngineStateReorgHandled:
		return true
	}

	panic(fmt.Sprintf("unexpected invalid state: %d", state))
}

func (event EngineLogEvent) String() string {
	switch event {
	case EngineEventGotLog:
		return "Got Log"
	case EngineEventProcess:
		return "Process"
	case EngineEventReorg:
		return "Got Reorg"
	case EngineEventHandleReorg:
		return "Handle Reorg"
	}

	panic(fmt.Sprintf("unexpected invalid event: %d", event))
}

func (event EngineLogEvent) IsValid() bool {
	switch event {
	case EngineEventInvalid:
		return false
	case
		EngineEventGotLog,
		EngineEventProcess,
		EngineEventReorg,
		EngineEventHandleReorg:
		return true
	}

	panic(fmt.Sprintf("unexpected invalid event: %d", event))
}
