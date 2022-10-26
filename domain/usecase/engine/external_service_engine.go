package engine

import "github.com/ethereum/go-ethereum/core/types"

// ServiceStateTracker is the service's implementation of chain reorganization state machine
type ServiceStateTracker interface {
	SetServiceState(ItemKey, ServiceItemState) // Overwrites state blindly
	GetServiceState(ItemKey) ServiceItemState  // Gets current item state
}

// ServiceEngine defines what service should implement and inject into engine.
type ServiceEngine interface {
	// ServiceStateTracker returns service-specific finite state machine.
	ServiceStateTracker() (ServiceStateTracker, error)

	// MapLogToItem maps Ethereum event log into service-specific type T.
	MapLogToItem(l *types.Log) (ServiceItem, error)

	// ActionOptions can be implemented to define arbitary options to be passed to ItemAction.
	ProcessOptions(ServiceItem, EngineLogState, ServiceItemState) ([]interface{}, error)

	// ProcessItem is called every time a new, fresh log is converted into ServiceItem,
	// The state returned represents the service state that will be assigned to the ServiceItem.
	ProcessItem(ServiceItem, EngineLogState, ServiceItemState, ...interface{}) (ServiceItemState, error)

	// ReorgOption can be implemented to define arbitary options to be passed to HandleReorg.
	ReorgOptions(ServiceItem, EngineLogState, ServiceItemState) ([]interface{}, error)

	// HandleReorg is called in *engine.handleReorg.
	HandleReorg(ServiceItem, EngineLogState, ServiceItemState, ...interface{}) (ServiceItemState, error)

	// TODO: work this out
	HandleEmitterError(error) error
}
