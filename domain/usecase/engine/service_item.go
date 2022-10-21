package engine

// ItemKey (K) is what used by ServiceFSM[K] to access ServiceItemState
type ItemKey interface {
	BlockNumber() uint64
}

// ServiceItem is The service "domain"-type representation of the log
type ServiceItem[K ItemKey] interface {
	ItemKey(...interface{}) K
	DebugString() string
}
