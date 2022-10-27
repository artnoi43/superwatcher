package uniswapv3factoryengine

import (
	"fmt"
)

type (
	uniswapv3PoolFactoryState uint8
	uniswapv3PoolFactoryEvent uint8

	uniswapv3StateTableKey struct {
		state uniswapv3PoolFactoryState
		event uniswapv3PoolFactoryEvent
	}
)

const (
	PoolFactoryStateNull uniswapv3PoolFactoryState = iota
	PoolFactoryStateCreated

	PoolFactoryEventPoolCreated uniswapv3PoolFactoryEvent = iota
)

var uniswapv3PoolFactoryStateTransitionTable = map[uniswapv3StateTableKey]uniswapv3PoolFactoryState{
	{state: PoolFactoryStateNull, event: PoolFactoryEventPoolCreated}: PoolFactoryStateCreated,
}

func (state uniswapv3PoolFactoryState) String() string {
	switch state {
	case PoolFactoryStateNull:
		return "NULL"
	case PoolFactoryStateCreated:
		return "POOLCREATED"
	}

	panic(fmt.Sprintf("invalid state: %d", state))
}

func (state uniswapv3PoolFactoryState) IsValid() bool {
	switch state {
	case
		PoolFactoryStateNull,
		PoolFactoryStateCreated:
		return true
	}

	return false
}

func (state uniswapv3PoolFactoryState) Fire(event uniswapv3PoolFactoryEvent) uniswapv3PoolFactoryState {
	newState, found := uniswapv3PoolFactoryStateTransitionTable[uniswapv3StateTableKey{
		state: state,
		event: event,
	}]

	if !found {
		panic("unknown path in state transition table")
	}

	state = newState
	return state
}

func (event uniswapv3PoolFactoryEvent) String() string {
	switch event {
	case PoolFactoryEventPoolCreated:
		return "PoolCreated"
	}

	panic(fmt.Sprintf("invalid event: %d", event))
}

func (event uniswapv3PoolFactoryEvent) IsValid() bool {
	switch event {
	case PoolFactoryEventPoolCreated:
		return true
	}

	return false
}
