package engine

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

// DEPRECATED - use handleLog instead
func handleBlock(
	block *reorg.BlockInfo,
	serviceEngine ServiceEngine,
	serviceStateTracker ServiceStateTracker,
	engineStateTracker EngineStateTracker,
	debugMode bool,
) error {
	for _, log := range block.Logs {
		if err := handleLog(
			log,
			serviceEngine,
			serviceStateTracker,
			engineStateTracker,
			debugMode,
		); err != nil {
			return errors.Wrapf(err, "error handling block number %d txHash %s", log.BlockNumber, log.TxHash.String())
		}
	}

	return nil
}

// handleLog is a pure function for handling a single Ethereum event log
func handleLog(
	log *types.Log,
	serviceEngine ServiceEngine,
	serviceStateTracker ServiceStateTracker,
	engineStateTracker EngineStateTracker,
	debugMode bool,
) error {
	debug.DebugMsg(debugMode, "*engine.handleLog: gotLog", zap.Any("log", log))
	if log.Removed {
		// TODO: Now what??
		logger.Info(
			"got removed log",
			zap.String("address", log.Address.String()),
			zap.String("txHash", log.TxHash.String()),
		)
	}

	engineKey := engineLogStateKeyFromLog(log)
	debug.DebugMsg(debugMode, "*engine.handleLog: got engineKey", zap.Any("key", engineKey), zap.String("keyType", reflect.TypeOf(engineKey).String()))
	engineState := engineStateTracker.GetEngineState(engineKey)

	switch engineState {
	case EngineLogStateNull:
		engineState.Fire(EngineLogEventGotLog)
		engineStateTracker.SetEngineState(engineKey, engineState)
	case
		EngineLogStateSeen,
		EngineLogStateReorged,
		EngineLogStateProcessed,
		EngineLogStateReorgHandled:

		// If we saw/processed this log/item, skip it
		debug.DebugMsg(
			debugMode,
			"handleLog skip due to engineState",
			zap.String("engineState", engineState.String()),
		)
		return nil
	}

	item, err := serviceEngine.MapLogToItem(log)
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to map log (txHash %s) to service item",
			log.TxHash.String(),
		)
	}

	key := item.ItemKey()
	itemServiceState := serviceStateTracker.GetServiceState(key)
	if itemServiceState == nil {
		logger.Panic("nil service state", zap.Any("itemKey", key))
	}

	actionOptions, err := serviceEngine.ProcessOptions(item, engineState, itemServiceState)
	if err != nil {
		return errors.Wrap(err, "failed to get itemAction options from service")
	}
	// TODO: Or the returned type from ProcessItem should be event?
	processedState, err := serviceEngine.ProcessItem(
		item,
		engineState,
		itemServiceState,
		actionOptions...,
	// Get options for ProcessItem from serviceEngine code
	)
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to perform actions an %s",
			item.DebugString(),
		)
	}

	engineState.Fire(EngineLogEventProcess)
	if !engineState.IsValid() {
		return fmt.Errorf("invalid state %s", engineState.String())
	}

	// Update states
	engineStateTracker.SetEngineState(engineKey, engineState)
	serviceStateTracker.SetServiceState(key, processedState)

	return nil
}
