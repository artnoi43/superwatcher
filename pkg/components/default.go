package components

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
)

// NewDefault returns default implementations of Emitter and Engine.
// The EmitterClient is initialized and embedded to the returned engine within this function.
// This is the preferred way for initializing superwatcher components. If you don't need to
// interact with these components, you can use `NewSuperWatcherDefault` instead.
func NewDefault(
	conf *config.Config,
	ethClient superwatcher.EthClient,
	getStateDataGateway superwatcher.GetStateDataGateway,
	setStateDataGateway superwatcher.SetStateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
	addresses []common.Address,
	topics [][]common.Hash,
) (
	superwatcher.Emitter,
	superwatcher.Engine,
) {
	syncChan := make(chan struct{})
	pollResultChan := make(chan *superwatcher.PollResult)
	errChan := make(chan error)

	watcherEmitter := NewEmitterWithPoller(
		conf,
		ethClient,
		getStateDataGateway,
		addresses,
		topics,
		syncChan,
		pollResultChan,
		errChan,
	)

	watcherEngine := NewEngineWithEmitterClient(
		conf,
		serviceEngine,
		setStateDataGateway,
		syncChan,
		pollResultChan,
		errChan,
	)

	return watcherEmitter, watcherEngine
}
