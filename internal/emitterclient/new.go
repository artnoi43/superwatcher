package emitterclient

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
)

func New(
	emitterConfig *config.Config,
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
		debug:            debug,
	}
}
