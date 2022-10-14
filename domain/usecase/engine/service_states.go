package engine

type State interface {
	String() string
	IsValid() bool
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
