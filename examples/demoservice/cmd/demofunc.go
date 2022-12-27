package main

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	spwconf "github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/components"
)

// This demo function calls components.NewDefault, which is the preferred way to init superwatcher for most cases.
func newSuperWatcherPreferred( //nolint:unused
	conf *spwconf.Config,
	ethClient superwatcher.EthClient,
	addresses []common.Address,
	topics []common.Hash,
	stateDataGateway superwatcher.StateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
) (superwatcher.Emitter, superwatcher.Engine) {
	return components.NewDefault(
		conf,
		ethClient,
		// Both data gateways can be implemented separately by different structs,
		// but here in the demo it's just using default implementation.
		superwatcher.GetStateDataGatewayFunc(stateDataGateway.GetLastRecordedBlock), // stateDataGateway alone is fine too
		superwatcher.SetStateDataGatewayFunc(stateDataGateway.SetLastRecordedBlock), // stateDataGateway alone is fine too
		serviceEngine,
		addresses,
		[][]common.Hash{topics},
	)
}

// This demo function demonstrates how users can use the components package
// to init superwatcher components individually
func newSuperwatcherAdvanced( //nolint:unused
	conf *spwconf.Config,
	ethClient superwatcher.EthClient,
	addresses []common.Address,
	topics []common.Hash,
	stateDataGateway superwatcher.StateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
) (superwatcher.Emitter, superwatcher.Engine) {
	errChan := make(chan error)
	syncChan := make(chan struct{})
	resultChan := make(chan *superwatcher.FilterResult)

	emitter := components.NewEmitter(conf, ethClient, stateDataGateway, nil, syncChan, resultChan, errChan)
	emitterClient := components.NewEmitterClient(conf, syncChan, resultChan, errChan)
	engine := components.NewEngine(emitterClient, serviceEngine, stateDataGateway, conf.LogLevel)
	poller := components.NewPoller(nil, nil, true, conf.FilterRange, ethClient, conf.LogLevel)

	poller.SetAddresses(addresses)
	poller.SetTopics([][]common.Hash{topics})

	emitter.SetPoller(poller)

	return emitter, engine
}
