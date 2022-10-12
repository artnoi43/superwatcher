package fsm

type State interface {
	String() string
	IsValid() bool
}

type Event interface {
	String() string
}

type ServiceItemState State
type ServiceItemEvent Event

// ServiceItem is The service "domain"-type representation of the log
type ServiceItem interface{}

// ServiceFSM[T] is the service's implementation of chain reorganization state machine
// that operates on T ServiceItem
type ServiceFSM[T ServiceItem] interface {
	SetServiceState(T, ServiceItemState)                            // Overwrites state blindly
	GetServiceState(T) ServiceItemState                             // Gets current item state
	FireServiceEvent(T, ServiceItemEvent) (ServiceItemState, error) // Traverses FSM
}
