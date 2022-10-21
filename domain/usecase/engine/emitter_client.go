package engine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

type EmitterClient[T any] interface {
	WatcherCurrentLog() *types.Log
	WatcherCurrentBlock() *reorg.BlockInfo
	WatcherReorg() *reorg.BlockInfo
	WatcherError() error

	ToDomainData(*types.Log) (*T, error)
}

type emitterClient[T any] struct {
	logChan   <-chan *types.Log
	blockChan <-chan *reorg.BlockInfo
	reorgChan <-chan *reorg.BlockInfo
	errChan   <-chan error

	debug bool
}

func NewWatcherClient[T any](
	logChan <-chan *types.Log,
	blockChan <-chan *reorg.BlockInfo,
	reorgChan <-chan *reorg.BlockInfo,
	errChan <-chan error,
) EmitterClient[T] {
	return &emitterClient[T]{
		logChan:   logChan,
		blockChan: blockChan,
		errChan:   errChan,
		reorgChan: reorgChan,
	}
}

func NewEmitterClientDebug[T any](
	logChan <-chan *types.Log,
	blockChan <-chan *reorg.BlockInfo,
	reorgChan <-chan *reorg.BlockInfo,
	errChan <-chan error,
) EmitterClient[T] {
	client := NewWatcherClient[T](logChan, blockChan, reorgChan, errChan)
	client.(*emitterClient[T]).debug = true

	return client
}

func (c *emitterClient[T]) WatcherCurrentLog() *types.Log {
	l, ok := <-c.logChan
	if ok {
		return l
	}

	if c.debug {
		logger.Debug("WatcherCurrentLog - logChan is closed")
	}
	return nil
}

func (c *emitterClient[T]) WatcherCurrentBlock() *reorg.BlockInfo {
	b, ok := <-c.blockChan
	if ok {
		return b
	}

	if c.debug {
		logger.Debug("WatcherCurrentBlock - blockChan is closed")
	}
	return nil
}

func (c *emitterClient[T]) WatcherReorg() *reorg.BlockInfo {
	blockInfo, ok := <-c.reorgChan
	if ok {
		return blockInfo
	}

	if c.debug {
		logger.Debug("WatcherReorg - reorgChan is closed")
	}
	return nil
}

func (c *emitterClient[T]) WatcherError() error {
	err, ok := <-c.errChan
	if ok {
		return err
	}

	if c.debug {
		logger.Debug("WatcherError - errChan is closed")
	}
	return nil
}

func (c *emitterClient[T]) ToDomainData(l *types.Log) (*T, error) {
	return nil, errors.New("not implemented")
}

func (c *emitterClient[T]) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(c.debug, msg, fields...)
}
