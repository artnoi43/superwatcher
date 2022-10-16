package engine

// itemKey (K) is what used by EngineFSM[K] as key for accessing ServiceItemState[K]
type itemKey interface {
	// BlockNumber returns the block number the item was first seen.
	// This helps when clearing old states from the state machine.
	BlockNumber() uint64
}

// ServiceItem is The service "domain"-type representation of the log
type ServiceItem[K itemKey] interface {
	ItemKey(...interface{}) K
	DebugString() string
}
