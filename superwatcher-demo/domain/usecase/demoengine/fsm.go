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
	poolFactoryStates engine.ServiceFSM[entity.Uniswapv3FactoryWatcherKey]
}

func (fsm *demoFSM) SetServiceState(key DemoKey, state engine.ServiceItemState) {
	fsm.Lock()
	defer fsm.Unlock()

	switch key.GetUseCase() {
	case usecase.UseCaseUniswapv3Factory:
		poolFactoryKey, ok := key.(entity.Uniswapv3FactoryWatcherKey)
		if !ok {
			logger.Panic("key not Uniswapv3FactoryWatcherKey", zap.String("actual type", reflect.TypeOf(key).String()))
		}

		fsm.poolFactoryStates.SetServiceState(poolFactoryKey, state)
	}

	logger.Panic("unhandled usecase for *demoFSM.SetServiceState")
}

func (fsm *demoFSM) GetServiceState(key DemoKey) engine.ServiceItemState {
	fsm.RLock()
	defer fsm.RUnlock()

	switch key.GetUseCase() {
	case usecase.UseCaseUniswapv3Factory:
		poolFactoryKey, ok := key.(entity.Uniswapv3FactoryWatcherKey)
		if !ok {
			logger.Panic("key not Uniswapv3FactoryWatcherKey", zap.String("actual type", reflect.TypeOf(key).String()))
		}
		return fsm.poolFactoryStates.GetServiceState(poolFactoryKey)
	}

	logger.Panic("unhandled usecase for *demoFSM.GetServiceState")
	return nil
}
