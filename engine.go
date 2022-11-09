package superwatcher

import "context"

type WatcherEngine interface {
	Loop(context.Context) error
}
