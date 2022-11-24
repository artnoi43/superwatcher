package emitterclient

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// emitterClient is the actual implementation of Client.
// It uses channels to communicate with emitter.
type emitterClient struct {
	emitterConfig    *config.EmitterConfig
	emitterSyncChan  chan<- struct{}
	filterResultChan <-chan *superwatcher.FilterResult
	errChan          <-chan error

	debugger *debugger.Debugger
}

func New(
	emitterConfig *config.EmitterConfig,
	emitterSyncChan chan<- struct{},
	filterResultChan <-chan *superwatcher.FilterResult,
	errChan <-chan error,
	debug bool,
) superwatcher.EmitterClient {
	return &emitterClient{
		emitterConfig:    emitterConfig,
		filterResultChan: filterResultChan,
		emitterSyncChan:  emitterSyncChan,
		errChan:          errChan,
		debugger: &debugger.Debugger{
			Key:         "emitter-client",
			ShouldDebug: debug,
		},
	}
}

func (c *emitterClient) Shutdown() {
	c.debugger.Debug("emitterClient.Shutdown() called")

	if c.emitterSyncChan != nil {
		c.debugger.Debug("closing emitterClient.emitterSyncChan")
		close(c.emitterSyncChan)

	} else {
		c.debugger.Debug("emitterClient: emitterSyncChan was already closed")
	}
}

// WatcherNextFilterLogs sends a low-cost signal to emitter to return from emitter.filterLogs
func (c *emitterClient) WatcherEmitterSync() {
	c.emitterSyncChan <- struct{}{}
}

func (c *emitterClient) WatcherConfig() *config.EmitterConfig {
	return c.emitterConfig
}

func (c *emitterClient) WatcherResult() *superwatcher.FilterResult {
	result, ok := <-c.filterResultChan
	if ok {
		return result
	}

	c.debugger.Debug("WatcherCurrentLog - filterReorgChan is closed")
	return nil
}

func (c *emitterClient) WatcherError() error {
	err, ok := <-c.errChan
	if ok {
		return err
	}

	c.debugger.Debug("WatcherError - errChan is closed")
	return nil
}
