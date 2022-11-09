package superwatcher

import "context"

// Code that imports watcher should only use this method.
type WatcherEmitter interface {
	Loop(context.Context) error
	SyncsWithEngine()
	Shutdown()
}
