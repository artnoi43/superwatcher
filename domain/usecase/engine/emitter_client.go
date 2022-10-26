package engine

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

type EmitterClient[T any] interface {
	WatcherResult() *emitter.FilterResult
	WatcherNextFilterLogs()
	WatcherError() error
}

type emitterClient[T any] struct {
	emitterSyncChan  chan<- struct{}
	filterResultChan <-chan *emitter.FilterResult
	errChan          <-chan error

	debug bool
}

func NewEmitterClient[T any](
	emitterSyncChan chan<- struct{},
	filterResultChan <-chan *emitter.FilterResult,
	errChan <-chan error,
	debug bool,
) EmitterClient[T] {
	return &emitterClient[T]{
		filterResultChan: filterResultChan,
		emitterSyncChan:  emitterSyncChan,
		errChan:          errChan,
		debug:            debug,
	}
}

// WatcherNextFilterLogs sends a low-cost signal to emitter to return from emitter.filterLogs
func (c *emitterClient[T]) WatcherNextFilterLogs() {
	c.emitterSyncChan <- struct{}{}
}

func (c *emitterClient[T]) WatcherResult() *emitter.FilterResult {
	result, ok := <-c.filterResultChan
	if ok {
		return result
	}

	c.debugMsg("WatcherCurrentLog - filterReorgChan is closed")
	return nil
}

func (c *emitterClient[T]) WatcherError() error {
	err, ok := <-c.errChan
	if ok {
		return err
	}

	c.debugMsg("WatcherError - errChan is closed")
	return nil
}

func (c *emitterClient[T]) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(c.debug, msg, fields...)
}
