package engine

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/engine"
	"github.com/artnoi43/superwatcher/pkg/components/emitterclient"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

func New(
	emitterClient superwatcher.EmitterClient,
	serviceEngine superwatcher.ServiceEngine,
	stateDataGateway watcherstate.StateDataGateway,
	debug bool,
) superwatcher.WatcherEngine {
	if debug {
		logger.Debug("initializing watcherEngine")
	}

	return engine.New(
		emitterClient,
		serviceEngine,
		stateDataGateway,
		debug,
	)
}

// New creates a new emitterClient, and pair it with an engine
func NewWithClient(
	emitterConfig *config.EmitterConfig,
	serviceEngine superwatcher.ServiceEngine,
	stateDataGateway watcherstate.StateDataGateway,
	syncChan chan<- struct{},
	filterResultChan <-chan *superwatcher.FilterResult,
	errChan <-chan error,
	debug bool,
) superwatcher.WatcherEngine {
	// TODO: Do we still need EmitterClient?
	emitterClient := emitterclient.New(
		emitterConfig,
		syncChan,
		filterResultChan,
		errChan,
		debug,
	)

	return New(
		emitterClient,
		serviceEngine,
		stateDataGateway,
		debug,
	)
}
