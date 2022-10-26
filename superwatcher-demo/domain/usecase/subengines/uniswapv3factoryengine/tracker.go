package uniswapv3factoryengine

import (
	"reflect"
	"sync"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

type poolFactoryStateTracker struct {
	sync.RWMutex
	states map[entity.Uniswapv3FactoryWatcherKey]engine.ServiceItemState
}

func (fsm *poolFactoryStateTracker) SetServiceState(key engine.ItemKey, state engine.ServiceItemState) {
	fsm.Lock()
	defer fsm.Unlock()

	poolKey, ok := key.(entity.Uniswapv3FactoryWatcherKey)
	if !ok {
		logger.Panic(
			"type assetion failed: key is not of type entity.Uniswapv3FactoryWatcherKey",
			zap.String("actual type", reflect.TypeOf(key).String()),
		)
	}

	fsm.states[poolKey] = state
}

func (fsm *poolFactoryStateTracker) GetServiceState(key engine.ItemKey) engine.ServiceItemState {
	fsm.RLock()
	defer fsm.RUnlock()

	poolKey, ok := key.(entity.Uniswapv3FactoryWatcherKey)
	if !ok {
		logger.Panic(
			"type assetion failed: key is not of type entity.Uniswapv3FactoryWatcherKey",
			zap.String("actual type", reflect.TypeOf(key).String()),
		)
	}

	state := fsm.states[poolKey]
	if state == nil {
		return PoolFactoryStateNull
	} else {
		return state
	}
}
