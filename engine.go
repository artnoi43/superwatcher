package superwatcher

import "context"

// WatcherEngine executes business service logic with ServiceEngine
type WatcherEngine interface {
	Loop(context.Context) error
}
