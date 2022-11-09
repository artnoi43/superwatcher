package engine

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
)

// newWatcherEngine returns default implementation of WatcherEngine
func New(
	client superwatcher.EmitterClient,
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
