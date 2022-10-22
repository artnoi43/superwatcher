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

type demoFSM struct {
	sync.RWMutex
	poolFactoryStates engine.ServiceFSM[subengines.DemoKey]
}

func NewDemoFSM(
	poolFactoryStates engine.ServiceFSM[subengines.DemoKey],
) engine.ServiceFSM[subengines.DemoKey] {
	return &demoFSM{
		poolFactoryStates: poolFactoryStates,
	}
}

func (fsm *demoFSM) SetServiceState(key subengines.DemoKey, state engine.ServiceItemState) {
	fsm.Lock()
	defer fsm.Unlock()

	stateUseCase := key.ForSubEngine()
	switch stateUseCase {
	case subengines.SubEngineUniswapv3Factory:
		poolFactoryKey, ok := key.(entity.Uniswapv3FactoryWatcherKey)
		if !ok {
			logger.Panic("key not Uniswapv3FactoryWatcherKey", zap.String("actual type", reflect.TypeOf(key).String()))
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

func (fsm *demoFSM) GetServiceState(key subengines.DemoKey) engine.ServiceItemState {
	fsm.RLock()
	defer fsm.RUnlock()

	stateUseCase := key.ForSubEngine()
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
