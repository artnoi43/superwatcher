package engine

// State and Event are very generic interfaces that get reused
type (
	State interface {
		String() string
		IsValid() bool

		// Fire traverses the state transition table,
		// sets the calling state to new state,
		// and returns that new state for other code to use.
		Fire(Event) State
	}

	Event interface {
		String() string
		IsValid() bool
	}
)
