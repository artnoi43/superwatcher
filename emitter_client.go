package superwatcher

import "github.com/artnoi43/superwatcher/config"

// EmitterClient interfaces with Emitter. It can help abstract the complexity
// of receiving of channel data away from Engine.
// It can be ignored by superwatcher users if they are not implementing their own Engine.
type EmitterClient interface {
	// Returns result from WatherEmitter to caller
	WatcherResult() *FilterResult
	// Sends sync signal to WatcherEmitter so it can continue
	WatcherEmitterSync()
	// Returns error sent by WatcherEmitter
	WatcherError() error
	// Returns config used to create its WatcherEmitter
	WatcherConfig() *config.Config
	// Closes WatcherEmitter comms channels
	Shutdown()
}
