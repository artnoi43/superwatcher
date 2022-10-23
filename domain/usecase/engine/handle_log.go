package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

// DEPRECATED - use handleLog instead
func handleBlock[K ItemKey, T ServiceItem[K]](
	block *reorg.BlockInfo,
	serviceEngine ServiceEngine[K, T],
	serviceFSM ServiceFSM[K],
	engineFSM EngineFSM,
	debugMode bool,
) error {
	for _, log := range block.Logs {
		if err := handleLog(
			log,
			serviceEngine,
			serviceFSM,
			engineFSM,
			debugMode,
		); err != nil {
			return errors.Wrapf(err, "error handling block number %d txHash %s", log.BlockNumber, log.TxHash.String())
		}
	}

	return nil
}

// handleLog is a pure function for handling a single Ethereum event log
func handleLog[K ItemKey, T ServiceItem[K]](
	log *types.Log,
	serviceEngine ServiceEngine[K, T],
	serviceFSM ServiceFSM[K],
	engineFSM EngineFSM,
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
	engineState := engineFSM.GetEngineState(engineKey)
	switch engineState {
	case EngineStateNull:
		engineState.Fire(EngineEventGotLog)
		engineFSM.SetEngineState(engineKey, engineState)
	case
		EngineStateSeen,
		EngineStateReorged,
		EngineStateProcessed,
		EngineStateReorgHandled:

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
	itemServiceState := serviceFSM.GetServiceState(key)
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

	engineState.Fire(EngineEventProcess)
	if !engineState.IsValid() {
		return fmt.Errorf("invalid state %s", engineState.String())
	}

	// Update states
	engineFSM.SetEngineState(engineKey, engineState)
	serviceFSM.SetServiceState(key, processedState)

	return nil
}
