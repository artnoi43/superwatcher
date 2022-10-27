package engine

import (
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter"
	"github.com/artnoi43/superwatcher/domain/usecase/emitterclient"
)

// newWatcherEngine returns default implementation of WatcherEngine
func newWatcherEngine(
	client emitterclient.Client,
	serviceEngine ServiceEngine,
	statDataGateway datagateway.StateDataGateway,
	debug bool,
) WatcherEngine {
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
	serviceEngine ServiceEngine,
	stateDataGateway datagateway.StateDataGateway,
	syncChan chan<- struct{},
	filterResultChan <-chan *emitter.FilterResult,
	errChan <-chan error,
	debug bool,
) WatcherEngine {

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
