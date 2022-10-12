package engine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
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

	debug bool
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

func NewWatcherClientDebug[T any](
	logChan <-chan *types.Log,
	errChan <-chan error,
	reorgChan <-chan *reorg.BlockInfo,
	adapter Adapter[T],
) WatcherClient[T] {
	client := NewWatcherClient(logChan, errChan, reorgChan, adapter)
	client.(*watcherClient[T]).debug = true

	return client
}

func (c *watcherClient[T]) WatcherCurrentLog() *types.Log {
	l, ok := <-c.logChan
	if ok {
		return l
	}

	if c.debug {
		logger.Debug("WatcherCurrentLog - logChan is closed")
	}
	return nil
}

func (c *watcherClient[T]) WatcherError() error {
	err, ok := <-c.errChan
	if ok {
		return err
	}

	if c.debug {
		logger.Debug("WatcherError - errChan is closed")
	}
	return nil
}

func (c *watcherClient[T]) WatcherReorg() *reorg.BlockInfo {
	blockInfo, ok := <-c.reorgChan
	if ok {
		return blockInfo
	}

	if c.debug {
		logger.Debug("WatcherReorg - reorgChan is closed")
	}
	return nil
}

func (c *watcherClient[T]) ToDomainData(l *types.Log) (*T, error) {
	return nil, errors.New("not implemented")
}
