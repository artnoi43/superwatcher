package engine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// ServiceItem is The service "domain"-type representation of the log
type ServiceItem interface {
	ItemKey() string
}

// ServiceFSM[T] is the service's implementation of chain reorganization state machine
// that operates on T ServiceItem
type ServiceFSM[T ServiceItem] interface {
	SetServiceState(T, ServiceItemState)                            // Overwrites state blindly
	GetServiceState(T) ServiceItemState                             // Gets current item state
	FireServiceEvent(T, ServiceItemEvent) (ServiceItemState, error) // Traverses FSM
}

// ServiceEngine[T] defines what service should implement and inject into engine.
type ServiceEngine[T ServiceItem] interface {
	ServiceStateTracker() (ServiceFSM[T], error)
	MapLogToItem(l *types.Log) (T, error)
	ItemAction(T) error
	HandleReorg(T) error
	HandleEmitterError(error) error
}

type engine[T ServiceItem] struct {
	client        watcherClient[T]
	serviceEngine ServiceEngine[T]
	engineFSM     EngineFSM[T]
}

func (e *engine[T]) handleLog() error {
	return errors.New("not implemented")
}

func (e *engine[T]) handleBlock() error {
	return errors.New("not implemented")
}

func (e *engine[T]) handleReorg() error {
	return errors.New("not implemented")
}

func (e *engine[T]) handleError() error {
	return errors.New("not implemented")
}
