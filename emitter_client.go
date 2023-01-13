package superwatcher

// EmitterClient interfaces with Emitter. It can help abstract the complexity
// of receiving of channel data away from Engine.
// It can be ignored by superwatcher users if they are not implementing their own Engine.
type EmitterClient interface {
	// WatcherResult returns result from Emitter to caller
	WatcherResult() *PollerResult
	// WatcherError returns error sent by Emitter
	WatcherError() error
	// WatcherConfig returns config used to create its Emitter
	WatcherConfig() *Config
	// SyncsEmitter sends sync signal to Emitter so it can continue
	SyncsEmitter()
	// Shutdown closes Emitter comms channels
	Shutdown()
}
