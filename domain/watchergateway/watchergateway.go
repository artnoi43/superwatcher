package watchergateway

import "github.com/ethereum/go-ethereum/core/types"

type WatcherClient interface {
	WatcherCurrentLog() *types.Log
	WatcherError() error
	WatcherReorg() bool
}

type watcherClient[T any] struct {
	logChan   <-chan *types.Log
	errChan   <-chan error
	reorgChan <-chan *struct{}

	adapter Adapter[T]
}

func NewWatcherClient[T any](
	logChan <-chan *types.Log,
	errChan <-chan error,
	reorgChan <-chan *struct{},
	adapter Adapter[T],
) WatcherClient {
	return &watcherClient[T]{
		logChan:   logChan,
		errChan:   errChan,
		reorgChan: reorgChan,
		adapter:   adapter,
	}
}

func (c *watcherClient[T]) WatcherCurrentLog() *types.Log {
	return <-c.logChan
}

func (c *watcherClient[T]) WatcherError() error {
	return <-c.errChan
}

func (c *watcherClient[T]) WatcherReorg() bool {
	return <-c.reorgChan != nil
}
