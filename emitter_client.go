package superwatcher

import "github.com/artnoi43/superwatcher/config"

// EmitterClient interfaces with WatcherEmitter
type EmitterClient interface {
	WatcherResult() *FilterResult         // Returns result from WatherEmitter to caller
	WatcherEmitterSync()                  // Sends sync signal to WatcherEmitter so it can continue
	WatcherError() error                  // Returns error sent by WatcherEmitter
	WatcherConfig() *config.EmitterConfig // Returns config used to create its WatcherEmitter
	Shutdown()                            // Closes WatcherEmitter comms channels
}
