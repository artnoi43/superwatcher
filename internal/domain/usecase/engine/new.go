package engine

import (
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/domain/usecase/emitterclient"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"
)

// newWatcherEngine returns default implementation of WatcherEngine
func newWatcherEngine(
	client emitterclient.Client,
	serviceEngine superwatcher.ServiceEngine,
	statDataGateway watcherstate.StateDataGateway,
	debug bool,
) superwatcher.WatcherEngine {
	return &engine{
		emitterClient:    client,
		serviceEngine:    serviceEngine,
		stateDataGateway: statDataGateway,
		metadataTracker:  NewTracker(debug),
		debug:            debug,
	}
}

// New creates a new emitter.emitter, and pair it with an engine
func New(
	emitterConfig *config.Config,
	serviceEngine superwatcher.ServiceEngine,
	stateDataGateway watcherstate.StateDataGateway,
	syncChan chan<- struct{},
	filterResultChan <-chan *superwatcher.FilterResult,
	errChan <-chan error,
	debug bool,
) superwatcher.WatcherEngine {

	// TODO: Do we still need EmitterClient?
	emitterClient := emitterclient.NewEmitterClient(
		emitterConfig,
		syncChan,
		filterResultChan,
		errChan,
		debug,
	)

	return newWatcherEngine(
		emitterClient,
		serviceEngine,
		stateDataGateway,
		debug,
	)
}
