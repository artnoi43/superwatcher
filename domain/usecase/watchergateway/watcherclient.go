package watchergateway

import (
	"github.com/artnoi43/superwatcher/domain/usecase/watcher/reorg"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type WatcherClient[T any] interface {
	WatcherCurrentLog() *types.Log
	WatcherError() error
	WatcherReorg() *reorg.BlockInfo

	ToDomainData(*types.Log) (*T, error)
}

type watcherClient[T any] struct {
	logChan   <-chan *types.Log
	errChan   <-chan error
	reorgChan <-chan *reorg.BlockInfo

	adapter Adapter[T]
}

func NewWatcherClient[T any](
	logChan <-chan *types.Log,
	errChan <-chan error,
	reorgChan <-chan *reorg.BlockInfo,
	adapter Adapter[T],
) WatcherClient[T] {
	return &watcherClient[T]{
		logChan:   logChan,
		errChan:   errChan,
		reorgChan: reorgChan,
		adapter:   adapter,
	}
}

func (c *watcherClient[T]) WatcherCurrentLog() *types.Log {
	l, closed := <-c.logChan
	if !closed {
		return l
	}
	return nil
}

func (c *watcherClient[T]) WatcherError() error {
	err, closed := <-c.errChan
	if !closed {
		return err
	}
	return nil
}

func (c *watcherClient[T]) WatcherReorg() *reorg.BlockInfo {
	blockInfo, closed := <-c.reorgChan
	if !closed {
		return blockInfo
	}
	return nil
}

func (c *watcherClient[T]) ToDomainData(l *types.Log) (*T, error) {
	return nil, errors.New("not implemented")
}
