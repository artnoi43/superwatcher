package watchergateway

// Adapter[T] represents a type whose Serialize method can convert
// Ethereum event log data from bytes to T
type Adapter[T any] interface {
	// Serialize parses Ethereum log data bytes into ServiceData.
	Serialize(
		eventKey string,
		logData []byte,
		opts ...any,
	) (
		T,
		error,
	)
}
