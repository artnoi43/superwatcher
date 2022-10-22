package demoengine

import (
	"errors"
	"reflect"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase"
)

func (e *demoEngine) itemToService(item engine.ServiceItem[usecase.DemoKey]) engine.ServiceEngine[usecase.DemoKey, engine.ServiceItem[usecase.DemoKey]] {
	itemUseCase := item.ItemKey().GetUseCase()

	serviceEngine, ok := e.services[itemUseCase]
	if !ok {
		logger.Panic(
			"usecase has no service",
			zap.String("usecase", itemUseCase.String()),
			zap.String("type of item", reflect.TypeOf(item).String()),
		)
	}

	return serviceEngine
}

func (e *demoEngine) ServiceStateTracker() (
	engine.ServiceFSM[usecase.DemoKey],
	error,
) {
	if e.fsm == nil {
		return nil, errors.New("nil *demoEngine.fsm")
	}

	return e.fsm, nil
}

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *demoEngine) MapLogToItem(
	log *types.Log,
) (
	engine.ServiceItem[usecase.DemoKey],
	error,
) {
	logUseCase, ok := e.usecases[log.Address]
	if !ok {
		logger.Panic("usecase not found", zap.String("usecase", logUseCase.String()))
	}

	serviceEngine, ok := e.services[logUseCase]
	if !ok {
		logger.Panic("")
	}

	return serviceEngine.MapLogToItem(log)
}

// Unused by this service
func (e *demoEngine) ActionOptions(
	item engine.ServiceItem[usecase.DemoKey],
	engineState engine.EngineLogState,
	serviceState engine.ServiceItemState,
) (
	[]interface{},
	error,
) {
	serviceEngine := e.itemToService(item)
	return serviceEngine.ActionOptions(item, engineState, serviceState)
}

// ItemAction just logs incoming pool
func (e *demoEngine) ItemAction(
	item engine.ServiceItem[usecase.DemoKey],
	engineState engine.EngineLogState,
	serviceState engine.ServiceItemState,
	options ...interface{},
) (
	engine.ServiceItemState,
	error,
) {
	if serviceState == nil {
		logger.Panic("nil serviceState")
	}
	if item == nil {
		logger.Panic("nil item")
	}

	serviceEngine := e.itemToService(item)
	return serviceEngine.ItemAction(item, engineState, serviceState, options)
}

// Unused by this service
func (e *demoEngine) ReorgOptions(
	item engine.ServiceItem[usecase.DemoKey],
	engineState engine.EngineLogState,
	serviceState engine.ServiceItemState,
) (
	[]interface{},
	error,
) {
	return nil, nil
}

// HandleReorg handles reorged event
// In uniswapv3poolfactory case, we only revert PoolCreated in the db.
// Other service may need more elaborate HandleReorg.
func (e *demoEngine) HandleReorg(
	item engine.ServiceItem[usecase.DemoKey],
	engineState engine.EngineLogState,
	serviceState engine.ServiceItemState,
	options ...interface{},
) (
	engine.ServiceItemState,
	error,
) {
	// TODO: remove the nil check if-blocks
	if serviceState == nil {
		logger.Panic("nil serviceState")
	}
	if item == nil {
		logger.Panic("nil item")
	}

	serviceEngine := e.itemToService(item)
	return serviceEngine.HandleReorg(item, engineState, serviceState, options...)
}

// Unused by this service
func (e *demoEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))
	return nil
}
