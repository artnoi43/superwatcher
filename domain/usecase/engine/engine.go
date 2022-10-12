package engine

import "github.com/ethereum/go-ethereum/core/types"

// ServiceItem is The service "domain"-type representation of the log
type ServiceItem interface{}

// ServiceFSM[T] is the service's implementation of chain reorganization state machine
// that operates on T ServiceItem
type ServiceFSM[T ServiceItem] interface {
	SetServiceState(T, ServiceItemState)                            // Overwrites state blindly
	GetServiceState(T) ServiceItemState                             // Gets current item state
	FireServiceEvent(T, ServiceItemEvent) (ServiceItemState, error) // Traverses FSM
}

type ServiceEngine[T ServiceItem] interface {
	MapLogToItem(l *types.Log) (T, error)
	ItemAction(T) error
	HandleReorg(T) error
	HandleEmitterError(error) error
}

type engine[T ServiceItem] struct {
	serviceEngine ServiceEngine[T]
}
