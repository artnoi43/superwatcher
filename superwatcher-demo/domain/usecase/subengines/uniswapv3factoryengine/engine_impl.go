package uniswapv3factoryengine

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *uniswapv3PoolFactoryEngine) MapLogToItem(
	log *types.Log,
) (
	engine.ServiceItem,
	error,
) {
	logEventKey := log.Topics[0]
	for _, event := range e.contractEvents {
		// This engine is supposed to handle more than 1 event,
		// but it's not yet finished now.
		if logEventKey == event.ID || event.Name == "PoolCreated" {
			return mapLogToPoolCreated(e.contractABI, event.Name, log)
		}
	}

	return nil, fmt.Errorf("event topic %s not found", logEventKey)
}

// Unused by this service
func (e *uniswapv3PoolFactoryEngine) ProcessOptions(
	pool engine.ServiceItem,
	engineState engine.EngineLogState,
	serviceState engine.ServiceItemState,
) (
	[]interface{},
	error,
) {

	return nil, nil
}

// ProcessItem just logs incoming pool
func (e *uniswapv3PoolFactoryEngine) ProcessItem(
	pool engine.ServiceItem,
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

	poolState := serviceState.(uniswapv3PoolFactoryState)
	switch engineState {
	case engine.EngineLogStateNull:
		switch poolState {
		case PoolFactoryStateNull:
			// Pretend this is DB operations
			logger.Info("DEMO: got new poolCreated, writing to db..")
			return PoolFactoryStateCreated, nil
		}
	}

	return serviceState, nil
}

// Unused by this service
func (e *uniswapv3PoolFactoryEngine) ReorgOptions(
	pool engine.ServiceItem,
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
func (e *uniswapv3PoolFactoryEngine) HandleReorg(
	item engine.ServiceItem,
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

	poolState, ok := serviceState.(uniswapv3PoolFactoryState)
	if !ok {
		logger.Panic(
			"type assertion failed: serviceState is not of type uniswapv3PoolFactoryState",
			zap.String("actual type", reflect.TypeOf(serviceState).String()),
		)
	}
	pool, ok := item.(*entity.Uniswapv3PoolCreated)
	if !ok {
		logger.Panic(
			"type assertion failed: item is not of type *entity.Uniswapv3PoolCreated",
			zap.String("actual type", reflect.TypeOf(item).String()),
		)
	}

	switch engineState {
	case engine.EngineLogStateProcessed:
		switch poolState {
		case PoolFactoryStateCreated:
			if err := e.revertPoolCreated(pool); err != nil {
				return serviceState, errors.Wrapf(
					err, "failed to revert poolCreated for pool %s",
					pool.Address.String(),
				)
			}

			return PoolFactoryStateNull, nil
		}
	}

	logger.Panic(
		"unhandled scenario",
		zap.String("engineState", engineState.String()),
		zap.String("poolState", poolState.String()),
	)

	return serviceState, nil
}

// Unused by this service
func (e *uniswapv3PoolFactoryEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))
	return nil
}
