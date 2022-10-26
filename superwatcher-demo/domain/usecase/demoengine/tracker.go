package demoengine

import (
	"reflect"
	"sync"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
)

type demoStateTracker struct {
	sync.RWMutex
	poolFactoryStates engine.ServiceStateTracker
}

func NewDemoFSM(
	poolFactoryStates engine.ServiceStateTracker,
) engine.ServiceStateTracker {
	return &demoStateTracker{
		poolFactoryStates: poolFactoryStates,
	}
}

func (fsm *demoStateTracker) SetServiceState(key engine.ItemKey, state engine.ServiceItemState) {
	fsm.Lock()
	defer fsm.Unlock()

	demoKey := subengines.AssertDemoKey(key)
	stateUseCase := demoKey.ForSubEngine()

	switch stateUseCase {
	case subengines.SubEngineUniswapv3Factory:
		poolFactoryKey, ok := key.(entity.Uniswapv3FactoryWatcherKey)
		if !ok {
			logger.Panic("type assertion failed: poolFactorykey is not Uniswapv3FactoryWatcherKey", zap.String("actual type", reflect.TypeOf(key).String()))
		}

		fsm.poolFactoryStates.SetServiceState(poolFactoryKey, state)

	default:
		logger.Panic(
			"unhandled usecase for *demoFSM.SetServiceState",
			zap.String("usecase", stateUseCase.String()),
			zap.Any("usecase", stateUseCase),
		)
	}

}

func (fsm *demoStateTracker) GetServiceState(key engine.ItemKey) engine.ServiceItemState {
	fsm.RLock()
	defer fsm.RUnlock()

	demoKey := subengines.AssertDemoKey(key)
	stateUseCase := demoKey.ForSubEngine()

	switch stateUseCase {
	case subengines.SubEngineUniswapv3Factory:
		poolFactoryKey, ok := key.(entity.Uniswapv3FactoryWatcherKey)
		if !ok {
			logger.Panic("key not Uniswapv3FactoryWatcherKey", zap.String("actual type", reflect.TypeOf(key).String()))
		}
		return fsm.poolFactoryStates.GetServiceState(poolFactoryKey)

	default:
		logger.Panic(
			"unhandled usecase for *demoFSM.GetServiceState",
			zap.String("usecase", stateUseCase.String()),
			zap.Any("usecase", stateUseCase),
		)
	}

	return nil
}
