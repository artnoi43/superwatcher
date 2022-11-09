package superwatcher

import "github.com/artnoi43/superwatcher/config"

// EmitterClient interfaces with emitter.WatcherEmitter via these methods
type EmitterClient interface {
	WatcherResult() *FilterResult
	WatcherEmitterSync()
	WatcherError() error
	WatcherConfig() *config.Config

	Shutdown()
}
