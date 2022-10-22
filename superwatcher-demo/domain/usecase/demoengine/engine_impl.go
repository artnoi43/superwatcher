package demoengine

import (
	"errors"
	"reflect"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
)

func (e *demoEngine) itemToService(item engine.ServiceItem[subengines.DemoKey]) engine.ServiceEngine[subengines.DemoKey, engine.ServiceItem[subengines.DemoKey]] {
	itemUseCase := item.ItemKey().ForSubEngine()

	serviceEngine, ok := e.services[itemUseCase]
	if !ok {
		logger.Panic(
			"usecase has no service",
			zap.String("subengine usecase", itemUseCase.String()),
			zap.String("type of item", reflect.TypeOf(item).String()),
		)
	}

	return serviceEngine
}

func (e *demoEngine) ServiceStateTracker() (
	engine.ServiceFSM[subengines.DemoKey],
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
	engine.ServiceItem[subengines.DemoKey],
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
func (e *demoEngine) ProcessOptions(
	item engine.ServiceItem[subengines.DemoKey],
	engineState engine.EngineLogState,
	serviceState engine.ServiceItemState,
) (
	[]interface{},
	error,
) {
	serviceEngine := e.itemToService(item)
	return serviceEngine.ProcessOptions(item, engineState, serviceState)
}

// ProcessItem just logs incoming pool
func (e *demoEngine) ProcessItem(
	item engine.ServiceItem[subengines.DemoKey],
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
	return serviceEngine.ProcessItem(item, engineState, serviceState, options)
}

// Unused by this service
func (e *demoEngine) ReorgOptions(
	item engine.ServiceItem[subengines.DemoKey],
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
	item engine.ServiceItem[subengines.DemoKey],
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
