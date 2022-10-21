package engine

import "sync"

type EngineFSM interface {
	SetEngineState(engineLogStateKey, EngineLogState)
	GetEngineState(engineLogStateKey) EngineLogState
}

type engineFSM struct {
	sync.RWMutex
	states map[engineLogStateKey]EngineLogState
}

func NewEngineFSM() EngineFSM {
	return &engineFSM{
		states: make(map[engineLogStateKey]EngineLogState),
	}
}

func (fsm *engineFSM) SetEngineState(key engineLogStateKey, newState EngineLogState) {
	fsm.Lock()
	defer fsm.Unlock()

	fsm.states[key] = newState
}

func (fsm *engineFSM) GetEngineState(key engineLogStateKey) EngineLogState {
	fsm.RLock()
	defer fsm.RUnlock()

	state, ok := fsm.states[key]
	if !ok {
		return EngineStateNull
	}

	return state
}
