package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

// handleLog is a pure function for handling a single Ethereum event log
func handleLog[K itemKey, T ServiceItem[K]](
	log *types.Log,
	serviceEngine ServiceEngine[K, T],
	serviceFSM ServiceFSM[K],
	engineFSM EngineFSM[K],
	debugMode bool,
) error {
	if log.Removed {
		// TODO: Now what??
		logger.Info(
			"got removed log",
			zap.String("address", log.Address.String()),
			zap.String("txHash", log.TxHash.String()),
		)
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
	engineState := engineFSM.GetEngineState(key)

	switch engineState {
	case EngineStateNull:
		engineState.Fire(EngineEventGotLog)
		engineFSM.SetEngineState(key, engineState)
	case
		EngineStateSeen,
		EngineStateReorged,
		EngineStateProcessed,
		EngineStateReorgHandled:

		// If we saw/processed this log/item, skip it
		debug.DebugMsg(debugMode, "handleLog skip due to engineState", zap.String("state", engineState.String()))
		return nil
	}

	itemServiceState := serviceFSM.GetServiceState(key)
	processedState, err := serviceEngine.ItemAction(
		item,
		engineState,
		itemServiceState,
		// Get options for ItemAction from serviceEngine code
		serviceEngine.ActionOptions(item, engineState, itemServiceState),
	)
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to perform actions an %s",
			item.DebugString(),
		)
	}

	engineState.Fire(EngineLogEvent(EngineStateProcessed))
	if !engineState.IsValid() {
		return fmt.Errorf("invalid state %s", engineState.String())
	}

	// Update states
	engineFSM.SetEngineState(key, engineState)
	serviceFSM.SetServiceState(key, processedState)

	return nil
}
