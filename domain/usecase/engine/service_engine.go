package engine

import "github.com/ethereum/go-ethereum/core/types"

// ServiceEngine[T] defines what service should implement and inject into engine.
type ServiceEngine[K ItemKey, T ServiceItem[K]] interface {
	// ServiceStateTracker returns service-specific finite state machine.
	ServiceStateTracker() (ServiceFSM[K], error)

	// MapLogToItem maps Ethereum event log into service-specific type T.
	MapLogToItem(l *types.Log) (T, error)

	// ActionOptions can be implemented to define arbitary options to be passed to ItemAction.
	ProcessOptions(T, EngineLogState, ServiceItemState) ([]interface{}, error)

	// ProcessItem is called every time a new, fresh log is converted into ServiceItem,
	// The state returned represents the service state that will be assigned to the ServiceItem.
	ProcessItem(T, EngineLogState, ServiceItemState, ...interface{}) (ServiceItemState, error)

	// ReorgOption can be implemented to define arbitary options to be passed to HandleReorg.
	ReorgOptions(T, EngineLogState, ServiceItemState) ([]interface{}, error)

	// HandleReorg is called in *engine.handleReorg.
	HandleReorg(T, EngineLogState, ServiceItemState, ...interface{}) (ServiceItemState, error)

	// TODO: work this out
	HandleEmitterError(error) error
}
