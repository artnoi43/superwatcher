package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// handleReorgedLog is a pure function for handling a reorganized event log.
func handleReorgedLog(
	reorgedLog *types.Log,
	serviceEngine ServiceEngine,
	serviceStateTracker ServiceStateTracker,
	engineStateTracker EngineStateTracker,
	debug bool,
) error {
	engineKey := engineLogStateKeyFromLog(reorgedLog)
	engineState := engineStateTracker.GetEngineState(engineKey)

	// TODO: Work this out.
	// As of now, we will only handle reorg if it's 1st reorg.
	var err error
	switch engineState {
	case
		// First reorg of this log
		EngineLogStateSeen,
		EngineLogStateProcessed:
		engineState.Fire(EngineLogEventReorg)
		if !engineState.IsValid() {
			return errors.Wrap(err, "failed to update engine state to EngineEventGotReorg")
		}
		engineStateTracker.SetEngineState(engineKey, engineState)
	}

	reorgedItem, err := serviceEngine.MapLogToItem(reorgedLog)
	if err != nil {
		return errors.Wrapf(err, "failed to map reorged log (txHash %s) to item", reorgedLog.TxHash.String())
	}

	key := reorgedItem.ItemKey()
	handleReorgOptions, err := serviceEngine.ReorgOptions(
		reorgedItem,
		engineState,
		serviceStateTracker.GetServiceState(key),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to get reorgOptions from service")
	}

	stateAfterHandledReorged, err := serviceEngine.HandleReorg(
		reorgedItem,
		engineState,
		serviceStateTracker.GetServiceState(key),
		handleReorgOptions,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to handle reorg for item %s", reorgedItem.DebugString())
	}

	engineState.Fire(EngineLogEventHandleReorg)
	if !engineState.IsValid() {
		return fmt.Errorf("invalid engineState: %s (%d)", engineState.String(), engineState)
	}
	engineStateTracker.SetEngineState(engineKey, engineState)
	serviceStateTracker.SetServiceState(key, stateAfterHandledReorged)

	return nil
}
