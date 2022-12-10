package engine

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
)

type (
	BlockState uint8
	BlockEvent uint8
)

const (
	StateNull BlockState = iota
	StateSeen
	StateProcessed
	StateReorged
	StateReorgHandled
	StateInvalid

	EventInvalid BlockEvent = iota
	EventGotLog
	EventProcess
	EventReorg
	EventHandleReorg
)

type stateEvent = struct {
	state BlockState
	event BlockEvent
}

var engineStateTransitionTable = map[stateEvent]BlockState{
	{state: StateNull, event: EventGotLog}:      StateSeen,
	{state: StateNull, event: EventProcess}:     StateInvalid,
	{state: StateNull, event: EventReorg}:       StateReorged,
	{state: StateNull, event: EventHandleReorg}: StateInvalid,

	{state: StateSeen, event: EventGotLog}:      StateSeen,
	{state: StateSeen, event: EventProcess}:     StateProcessed,
	{state: StateSeen, event: EventReorg}:       StateReorged,
	{state: StateSeen, event: EventHandleReorg}: StateReorgHandled,

	{state: StateProcessed, event: EventGotLog}:      StateProcessed,
	{state: StateProcessed, event: EventProcess}:     StateProcessed,
	{state: StateProcessed, event: EventReorg}:       StateReorged,
	{state: StateProcessed, event: EventHandleReorg}: StateReorgHandled,

	{state: StateReorged, event: EventGotLog}:      StateReorged,
	{state: StateReorged, event: EventProcess}:     StateInvalid,
	{state: StateReorged, event: EventReorg}:       StateReorged,
	{state: StateReorged, event: EventHandleReorg}: StateReorgHandled,

	{state: StateReorgHandled, event: EventGotLog}:      StateInvalid,
	{state: StateReorgHandled, event: EventProcess}:     StateInvalid,
	{state: StateReorgHandled, event: EventReorg}:       StateReorged,
	{state: StateReorgHandled, event: EventHandleReorg}: StateInvalid,

	{state: StateInvalid, event: EventGotLog}:      StateInvalid,
	{state: StateInvalid, event: EventProcess}:     StateInvalid,
	{state: StateInvalid, event: EventReorg}:       StateInvalid,
	{state: StateInvalid, event: EventHandleReorg}: StateInvalid,
}

func (state *BlockState) Fire(event BlockEvent) {
	if !event.IsValid() {
		logger.Panic("invalid event", zap.String("event", event.String()))
	}

	self := stateEvent{state: *state, event: event}
	newState := engineStateTransitionTable[self]
	*state = newState
}

func (state BlockState) String() string {
	switch state {
	case StateNull:
		return "NULL"
	case StateSeen:
		return "SEEN"
	case StateProcessed:
		return "PROCESSED"
	case StateReorged:
		return "REORGED"
	case StateReorgHandled:
		return "REORG_HANDLED"
	case StateInvalid:
		return "INVALID_ENGINE_STATE"
	}

	panic(fmt.Sprintf("unexpected invalid state: %d", state))
}

func (state BlockState) IsValid() bool {
	switch state {
	case StateInvalid:
		return false
	case
		StateNull,
		StateSeen,
		StateProcessed,
		StateReorged,
		StateReorgHandled:
		return true
	}

	panic(fmt.Sprintf("unexpected invalid state: %d", state))
}

func (event BlockEvent) String() string {
	switch event {
	case EventGotLog:
		return "Got Log"
	case EventProcess:
		return "Process"
	case EventReorg:
		return "Got Reorg"
	case EventHandleReorg:
		return "Handle Reorg"
	}

	panic(fmt.Sprintf("unexpected invalid event: %d", event))
}

func (event BlockEvent) IsValid() bool {
	switch event {
	case EventInvalid:
		return false
	case
		EventGotLog,
		EventProcess,
		EventReorg,
		EventHandleReorg:
		return true
	}

	panic(fmt.Sprintf("unexpected invalid event: %d", event))
}
