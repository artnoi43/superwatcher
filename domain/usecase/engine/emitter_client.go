package engine

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

type EmitterClient[T any] interface {
	WatcherResult() *emitter.FilterResult
	WatcherError() error
}

type emitterClient[T any] struct {
	filterResultChan <-chan *emitter.FilterResult
	errChan          <-chan error

	debug bool
}

func NewEmitterClient[T any](
	filterResultChan <-chan *emitter.FilterResult,
	errChan <-chan error,
	debug bool,
) EmitterClient[T] {
	return &emitterClient[T]{
		filterResultChan: filterResultChan,
		errChan:          errChan,
		debug:            debug,
	}
}

func (c *emitterClient[T]) WatcherResult() *emitter.FilterResult {
	result, ok := <-c.filterResultChan
	if ok {
		return result
	}

	if c.debug {
		logger.Debug("WatcherCurrentLog - filterReorgChan is closed")
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

func (c *emitterClient[T]) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(c.debug, msg, fields...)
}
