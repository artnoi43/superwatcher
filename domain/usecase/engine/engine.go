package engine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

type itemKey constraints.Ordered

// ServiceItem is The service "domain"-type representation of the log
type ServiceItem[K itemKey] interface {
	ItemKey() K
	DebugString() string
}

// ServiceFSM[T] is the service's implementation of chain reorganization state machine
// that operates on T ServiceItem
type ServiceFSM[K itemKey] interface {
	SetServiceState(K, ServiceItemState)                            // Overwrites state blindly
	GetServiceState(K) ServiceItemState                             // Gets current item state
	FireServiceEvent(K, ServiceItemEvent) (ServiceItemState, error) // Traverses FSM
}

// ServiceEngine[T] defines what service should implement and inject into engine.
type ServiceEngine[K itemKey, T ServiceItem[K]] interface {
	// ServiceStateTracker returns service-specific finite state machine.
	ServiceStateTracker() (ServiceFSM[K], error)

	// MapLogToItem maps Ethereum event log into service-specific type T.
	MapLogToItem(l *types.Log) (T, error)

	// ActionOptions can be implemented to define arbitary options to be passed to ItemAction.
	ActionOptions(T, EngineLogState, ServiceItemState) []interface{}

	// ItemAction is called every time a new, fresh log is converted into ServiceItem,
	// The state returned represents the service state that will be assigned to the ServiceItem.
	ItemAction(T, EngineLogState, ServiceItemState, ...interface{}) (State, error)

	// ReorgOption can be implemented to define arbitary options to be passed to HandleReorg.
	ReorgOptions(T) []interface{}

	// HandleReorg is called in *engine.handleReorg.
	HandleReorg(T, EngineLogState, ServiceItemState, ...interface{}) (State, error)

	// TODO: work this out
	HandleEmitterError(error) error
}

type engine[K itemKey, T ServiceItem[K]] struct {
	client        watcherClient[T]
	serviceEngine ServiceEngine[K, T]
	engineFSM     EngineFSM[K]
	debug         bool
}

func (e *engine[K, T]) handleLog() error {
	serviceEngine, serviceFSM, engineFSM, err := e.initStuff("handleLog")
	if err != nil {
		return err
	}

	for {
		newLog := e.client.WatcherCurrentLog()
		if err := handleLog(newLog, serviceEngine, serviceFSM, engineFSM, e.debug); err != nil {
			return errors.Wrap(err, "handleLog failed in engine.HandleLog")
		}
	}
}

func (e *engine[K, T]) handleBlock() error {
	serviceEngine, serviceFSM, engineFSM, err := e.initStuff("handleBlock")
	if err != nil {
		return err
	}

	for {
		newBlock := e.client.WatcherCurrentBlock()

		for _, log := range newBlock.Logs {
			if err := handleLog(log, serviceEngine, serviceFSM, engineFSM, e.debug); err != nil {
				return errors.Wrap(err, "handleLog failed in handleBlock")
			}
		}
	}
}

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
			currentEngineState := engineFSM.GetEngineState(key)

			// TODO: How to handle this?
			switch currentEngineState {
			case
				EngineStateSeen,
				EngineStateProcessed:
				currentEngineState.Fire(EngineEventGotReorg)
				if !currentEngineState.IsValid() {
					return errors.Wrap(err, "failed to update engine state to EngineEventGotReorg")
				}
			}

			handleReorgOptions := serviceEngine.ReorgOptions(reorgedItem)
			handledState, err := serviceEngine.HandleReorg(
				reorgedItem,
				currentEngineState,
				serviceFSM.GetServiceState(key),
				handleReorgOptions,
			)
			if err != nil {
				return errors.Wrapf(err, "failed to handle reorg for item %s", reorgedItem.DebugString())
			}

			serviceFSM.SetServiceState(key, handledState)
			engineFSM.SetEngineState(key, EngineStateProcessedReorged)
		}
	}
}

func (e *engine[K, T]) handleError() error {
	return errors.New("not implemented")
}

func (e *engine[K, T]) initStuff(method string) (ServiceEngine[K, T], ServiceFSM[K], EngineFSM[K], error) {
	serviceFSM, err := e.serviceEngine.ServiceStateTracker()
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to init serviceFSM for %s", method)
	}

	return e.serviceEngine, serviceFSM, e.engineFSM, nil
}
