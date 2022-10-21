package engine

type State interface {
	String() string
	IsValid() bool

	// Fire traverses the state transition table,
	// sets the calling state to new state,
	// and returns that new state for other code to use.
	Fire(Event) State
}

type Event interface {
	String() string
	IsValid() bool
}

type (
	ServiceItemState State
	ServiceItemEvent Event
)

// ServiceFSM[T] is the service's implementation of chain reorganization state machine
// that operates on T ServiceItem
type ServiceFSM[K ItemKey] interface {
	SetServiceState(K, ServiceItemState) // Overwrites state blindly
	GetServiceState(K) ServiceItemState  // Gets current item state
}
