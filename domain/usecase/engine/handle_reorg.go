package engine

import (
	"fmt"

	"github.com/pkg/errors"
)

func (e *engine[K, T]) handleReorg() error {
	serviceEngine, serviceFSM, engineFSM, err := e.initStuff("handleBlock")
	if err != nil {
		return err
	}

	for {
		reorgedBlock := e.client.WatcherReorg()
		for _, reorgedLog := range reorgedBlock.Logs {
			var err error
			reorgedItem, err := serviceEngine.MapLogToItem(reorgedLog)
			if err != nil {
				return errors.Wrapf(err, "failed to map reorged log (txHash %s) to item", reorgedLog.TxHash.String())
			}

			key := reorgedItem.ItemKey()
			engineState := engineFSM.GetEngineState(key)

			// TODO: Work this out.
			// As of now, we will only handle reorg if it's 1st reorg.
			switch engineState {
			case
				// First reorg of this log
				EngineStateSeen,
				EngineStateProcessed:
				engineState.Fire(EngineEventGotReorg)
				if !engineState.IsValid() {
					return errors.Wrap(err, "failed to update engine state to EngineEventGotReorg")
				}
				engineFSM.SetEngineState(key, engineState)
			}

			handleReorgOptions := serviceEngine.ReorgOptions(
				reorgedItem,
				engineState,
				serviceFSM.GetServiceState(key),
			)
			stateAfterHandledReorged, err := serviceEngine.HandleReorg(
				reorgedItem,
				engineState,
				serviceFSM.GetServiceState(key),
				handleReorgOptions,
			)
			if err != nil {
				return errors.Wrapf(err, "failed to handle reorg for item %s", reorgedItem.DebugString())
			}

			engineState.Fire(EngineEventHandleReorg)
			if !engineState.IsValid() {
				return fmt.Errorf("invalid engineState: %s (%d)", engineState.String(), engineState)
			}
			engineFSM.SetEngineState(key, engineState)
			serviceFSM.SetServiceState(key, stateAfterHandledReorged)
		}
	}
}
