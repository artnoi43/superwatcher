package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	watcherengine "github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *uniswapv3FactoryEngine) MapLogToItem(
	log *types.Log,
) (
	*entity.Uniswapv3PoolCreated,
	error,
) {
	contractAddr := log.Address
	contractABI, ok := e.mapAddrToABI[contractAddr]
	if !ok {
		return nil, fmt.Errorf("abi not found for address %s", contractAddr.String())
	}
	contractInterestingEvents, ok := e.mapAddrToEvents[contractAddr]
	if !ok {
		return nil, fmt.Errorf("events not found for address %s", contractAddr.String())
	}

	logEventKey := log.Topics[0]
	for _, event := range contractInterestingEvents {
		if logEventKey == event.ID {
			return mapLogToItem(contractABI, event.Name, log)
		}
	}

	return nil, fmt.Errorf("event topic %s not found", logEventKey)
}

// Unused by this service
func (e *uniswapv3FactoryEngine) ActionOptions(
	pool *entity.Uniswapv3PoolCreated,
	engineState watcherengine.EngineLogState,
	serviceState watcherengine.ServiceItemState,
) (
	[]interface{},
	error,
) {

	return nil, nil
}

func (e *uniswapv3FactoryEngine) ItemAction(
	*entity.Uniswapv3PoolCreated,
	watcherengine.EngineLogState,
	watcherengine.ServiceItemState,
	...interface{},
) (
	watcherengine.ServiceItemState,
	error,
) {
	return nil, errors.New("not implemented")
}

// Unused by this service
func (e *uniswapv3FactoryEngine) ReorgOptions(
	pool *entity.Uniswapv3PoolCreated,
	engineState watcherengine.EngineLogState,
	serviceState watcherengine.ServiceItemState,
) (
	[]interface{},
	error,
) {
	return nil, nil
}

// HandleReorg handles reorged event
// In uniswapv3poolfactory case, we only revert PoolCreated in the db.
// Other service may need more elaborate HandleReorg.
func (e *uniswapv3FactoryEngine) HandleReorg(
	pool *entity.Uniswapv3PoolCreated,
	engineState watcherengine.EngineLogState,
	serviceState watcherengine.ServiceItemState,
) (
	watcherengine.ServiceItemState,
	error,
) {
	// TODO: remove the nil check if-block
	if serviceState == nil {
		logger.Panic("nil serviceState")
	}

	poolState := serviceState.(uniswapv3PoolFactoryState)
	switch engineState {
	case watcherengine.EngineStateProcessed:
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
func (e *uniswapv3FactoryEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))
	return nil
}
