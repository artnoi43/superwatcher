package emitterclient

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

// Client interfaces with emitter.WatcherEmitter via these methods
type Client[T any] interface {
	WatcherResult() *emitter.FilterResult
	WatcherEmitterSync()
	WatcherError() error

	Shutdown()
}

// emitterClient is the actual implementation of Client.
// It uses channels to communicate with emitter.
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
) Client[T] {
	return &emitterClient[T]{
		filterResultChan: filterResultChan,
		emitterSyncChan:  emitterSyncChan,
		errChan:          errChan,
		debug:            debug,
	}
}

func (c *emitterClient[T]) Shutdown() {
	if c.emitterSyncChan != nil {
		c.debugMsg("closing emitterClient.emitterSyncChan")
		close(c.emitterSyncChan)
	} else {
		c.debugMsg("emitterClient: emitterSyncChan was already closed")
	}
}

// WatcherNextFilterLogs sends a low-cost signal to emitter to return from emitter.filterLogs
func (c *emitterClient[T]) WatcherEmitterSync() {
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
