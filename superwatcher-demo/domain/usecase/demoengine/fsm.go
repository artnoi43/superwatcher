package demoengine

import (
	"reflect"
	"sync"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase"
)

type demoFSM struct {
	sync.RWMutex
	poolFactoryStates engine.ServiceFSM[usecase.DemoKey]
}

func NewDemoFSM(
	poolFactoryStates engine.ServiceFSM[usecase.DemoKey],
) engine.ServiceFSM[usecase.DemoKey] {
	return &demoFSM{
		poolFactoryStates: poolFactoryStates,
	}
}

func (fsm *demoFSM) SetServiceState(key usecase.DemoKey, state engine.ServiceItemState) {
	fsm.Lock()
	defer fsm.Unlock()

	stateUseCase := key.GetUseCase()
	switch stateUseCase {
	case usecase.UseCaseUniswapv3Factory:
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

func (fsm *demoFSM) GetServiceState(key usecase.DemoKey) engine.ServiceItemState {
	fsm.RLock()
	defer fsm.RUnlock()

	stateUseCase := key.GetUseCase()
	switch stateUseCase {
	case usecase.UseCaseUniswapv3Factory:
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
