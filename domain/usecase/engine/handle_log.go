package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

// handleLog handles a single Ethereum event log
func handleLog[K itemKey, T ServiceItem[K]](
	log *types.Log,
	serviceEngine ServiceEngine[K, T],
	serviceFSM ServiceFSM[K],
	engineFSM EngineFSM[K],
	debugMode bool,
) error {
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
		EngineStateProcessedReorged:

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
