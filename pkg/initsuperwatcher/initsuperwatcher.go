package initsuperwatcher

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/components/emitter"
	"github.com/artnoi43/superwatcher/pkg/components/engine"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
)

// New returns default implementations of WatcherEmitter and WatcherEngine.
// The EmitterClient is initialized and embedded to the returned engine within this function.
// This is the preferred way for initializing superwatcher components.
func New[H superwatcher.EmitterBlockHeader](
	conf *config.Config,
	ethClient superwatcher.EthClient[H],
	stateDataGateway watcherstate.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	serviceEngine superwatcher.ServiceEngine,
	debug bool,
) (
	superwatcher.WatcherEmitter,
	superwatcher.WatcherEngine,
) {
	syncChan := make(chan struct{})
	filterResultChan := make(chan *superwatcher.FilterResult)
	errChan := make(chan error)

	watcherEmitter := emitter.New(
		conf,
		ethClient,
		stateDataGateway,
		addresses,
		topics,
		syncChan,
		filterResultChan,
		errChan,
		debug,
	)

	watcherEngine := engine.NewWithClient(
		conf,
		serviceEngine,
		stateDataGateway,
		syncChan,
		filterResultChan,
		errChan,
		debug,
	)

	return watcherEmitter, watcherEngine
}
