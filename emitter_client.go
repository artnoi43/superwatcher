package superwatcher

import "github.com/artnoi43/superwatcher/config"

// EmitterClient interfaces with Emitter. It can help abstract the complexity of receiving of channel data away from Engine.
// It can be ignored by superwatcher users if they are not implementing their own Engine.
type EmitterClient interface {
	WatcherResult() *FilterResult  // Returns result from WatherEmitter to caller
	WatcherEmitterSync()           // Sends sync signal to WatcherEmitter so it can continue
	WatcherError() error           // Returns error sent by WatcherEmitter
	WatcherConfig() *config.Config // Returns config used to create its WatcherEmitter
	Shutdown()                     // Closes WatcherEmitter comms channels
}
