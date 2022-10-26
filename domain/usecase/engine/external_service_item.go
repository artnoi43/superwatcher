package engine

// ItemKey is used by ServiceStateTracker to access state
// of ServiceItem
type ItemKey interface {
	BlockNumber() uint64
}

// ServiceItem is The service "domain"-type representation of the log
type ServiceItem interface {
	ItemKey(...interface{}) ItemKey
	DebugString() string
}
