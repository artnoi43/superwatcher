package emitterclient

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/logger/debug"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"
)

// Client interfaces with emitter.WatcherEmitter via these methods
type Client interface {
	WatcherResult() *superwatcher.FilterResult
	WatcherEmitterSync()
	WatcherError() error
	WatcherConfig() *config.Config

	Shutdown()
}

// emitterClient is the actual implementation of Client.
// It uses channels to communicate with emitter.
type emitterClient struct {
	emitterConfig    *config.Config
	emitterSyncChan  chan<- struct{}
	filterResultChan <-chan *superwatcher.FilterResult
	errChan          <-chan error

	debug bool
}

func NewEmitterClient(
	emitterConfig *config.Config,
	emitterSyncChan chan<- struct{},
	filterResultChan <-chan *superwatcher.FilterResult,
	errChan <-chan error,
	debug bool,
) Client {
	return &emitterClient{
		emitterConfig:    emitterConfig,
		filterResultChan: filterResultChan,
		emitterSyncChan:  emitterSyncChan,
		errChan:          errChan,
		debug:            debug,
	}
}

func (c *emitterClient) Shutdown() {
	if c.emitterSyncChan != nil {
		c.debugMsg("closing emitterClient.emitterSyncChan")
		close(c.emitterSyncChan)
	} else {
		c.debugMsg("emitterClient: emitterSyncChan was already closed")
	}
}

// WatcherNextFilterLogs sends a low-cost signal to emitter to return from emitter.filterLogs
func (c *emitterClient) WatcherEmitterSync() {
	c.emitterSyncChan <- struct{}{}
}

func (c *emitterClient) WatcherConfig() *config.Config {
	return c.emitterConfig
}

func (c *emitterClient) WatcherResult() *superwatcher.FilterResult {
	result, ok := <-c.filterResultChan
	if ok {
		return result
	}

	c.debugMsg("WatcherCurrentLog - filterReorgChan is closed")
	return nil
}

func (c *emitterClient) WatcherError() error {
	err, ok := <-c.errChan
	if ok {
		return err
	}

	c.debugMsg("WatcherError - errChan is closed")
	return nil
}

func (c *emitterClient) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(c.debug, msg, fields...)
}
