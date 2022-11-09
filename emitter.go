package superwatcher

import "context"

type WatcherEmitter interface {
	Loop(context.Context) error // Main emitter loop
	SyncsWithEngine()           // Waits until engine is done processing the last batch
	Shutdown()                  // Shutdown and closing Go channels
}
