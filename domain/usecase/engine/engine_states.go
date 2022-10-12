package engine

import "fmt"

type EngineLogState uint8
type EngineLogEvent uint8

type EngineFSM[T ServiceItem] interface {
	SetEngineState(T, EngineLogState)
	GetEngineState(T) EngineLogState
	FireEngineEvent(T, EngineLogEvent) (EngineLogState, error)
}

const (
	EngineStateNull EngineLogState = iota
	EngineStateSeen
	EngineStateProcessed
	EngineStateReorged
	EngineStateProcessedReorged
	EngineStateError

	EngineEventNull EngineLogEvent = iota
	EngineEventGotLog
	EngineEventProcess
	EngineEventGotReorg
	EngineEventHandleReorg
)

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
	case EngineStateProcessedReorged:
		return "PROCESSED_REORGED"
	case EngineStateError:
		return "ERROR"
	}

	panic(fmt.Sprintf("invalid state: %d", state))
}

func (state EngineLogState) IsValid() bool {
	switch state {
	case
		EngineStateNull,
		EngineStateSeen,
		EngineStateProcessed,
		EngineStateReorged,
		EngineStateProcessedReorged:
		return true
	}

	return false
}

func (event EngineLogEvent) String() string {
	switch event {
	case EngineEventGotLog:
		return "Got Log"
	case EngineEventProcess:
		return "Process"
	case EngineEventGotReorg:
		return "Got Reorg"
	case EngineEventHandleReorg:
		return "Handle Reorg"
	}

	panic(fmt.Sprintf("invalid event: %d", event))
}

func (event EngineLogEvent) IsValid() bool {
	switch event {
	case
		EngineEventGotLog,
		EngineEventProcess,
		EngineEventGotReorg,
		EngineEventHandleReorg:
		return true
	}

	return false
}
