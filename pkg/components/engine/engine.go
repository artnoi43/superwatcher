package engine

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/engine"
	"github.com/artnoi43/superwatcher/pkg/components/emitterclient"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
)

func New(
	emitterClient superwatcher.EmitterClient,
	serviceEngine superwatcher.ServiceEngine,
	stateDataGateway watcherstate.StateDataGateway,
	logLevel uint8,
) superwatcher.WatcherEngine {
	return engine.New(
		emitterClient,
		serviceEngine,
		stateDataGateway,
		logLevel,
	)
}

// New creates a new emitterClient, and pair it with an engine
func NewWithClient(
	conf *config.EmitterConfig,
	serviceEngine superwatcher.ServiceEngine,
	stateDataGateway watcherstate.StateDataGateway,
	syncChan chan<- struct{},
	filterResultChan <-chan *superwatcher.FilterResult,
	errChan <-chan error,
) superwatcher.WatcherEngine {
	// TODO: Do we still need EmitterClient?
	emitterClient := emitterclient.New(
		conf,
		syncChan,
		filterResultChan,
		errChan,
	)

	return New(
		emitterClient,
		serviceEngine,
		stateDataGateway,
		conf.LogLevel,
	)
}
