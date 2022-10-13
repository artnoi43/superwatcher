package engine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func handleLog[K itemKey, T ServiceItem[K]](
	log *types.Log,
	serviceEngine ServiceEngine[K, T],
	serviceFSM ServiceFSM[K],
	engineFSM EngineFSM[K],
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
	engineFSM.SetEngineState(key, EngineStateSeen)

	processedState, err := serviceEngine.ItemAction(
		item,
		serviceEngine.ActionOptions(item), // Get options for ItemAction
	)
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to perform actions an %s",
			item.DebugString(),
		)
	}

	// If processed just fine, update state
	serviceFSM.SetServiceState(key, processedState)
	engineFSM.SetEngineState(key, EngineStateProcessed)

	return nil
}
