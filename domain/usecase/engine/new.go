package engine

import (
	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter"
)

func New[K ItemKey, T ServiceItem[K]](
	serviceEngine ServiceEngine[K, T],
	stateDataGateway datagateway.StateDataGateway,
	syncChan chan<- struct{},
	filterResultChan <-chan *emitter.FilterResult,
	errChan <-chan error,
	debug bool,
) WatcherEngine {

	// TODO: Do we still need EmitterClient?
	emitterClient := NewEmitterClient[T](
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
