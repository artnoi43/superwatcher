package initsuperwatcher

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/components/emitter"
	"github.com/artnoi43/superwatcher/pkg/components/engine"
)

// New returns default implementations of WatcherEmitter and WatcherEngine.
// The EmitterClient is initialized and embedded to the returned engine within this function.
// This is the preferred way for initializing superwatcher components.
func New(
	conf *config.EmitterConfig,
	ethClient superwatcher.EthClient,
	getStateDataGateway superwatcher.GetStateDataGateway,
	setStateDataGateway superwatcher.SetStateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	serviceEngine superwatcher.ServiceEngine,
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
		getStateDataGateway,
		addresses,
		topics,
		syncChan,
		filterResultChan,
		errChan,
	)

	watcherEngine := engine.NewWithClient(
		conf,
		serviceEngine,
		setStateDataGateway,
		syncChan,
		filterResultChan,
		errChan,
	)

	return watcherEmitter, watcherEngine
}
