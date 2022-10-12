package engine

type State interface {
	String() string
	IsValid() bool
}

type Event interface {
	String() string
}

type ServiceItemState State
type ServiceItemEvent Event
