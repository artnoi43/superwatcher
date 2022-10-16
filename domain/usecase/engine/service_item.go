package engine

// itemKey (K) is what used by ServiceFSM[K] to access ServiceItemState
type itemKey interface {
	BlockNumber() uint64
}

// ServiceItem is The service "domain"-type representation of the log
type ServiceItem[K itemKey] interface {
	ItemKey(...interface{}) K
	DebugString() string
}
