package main

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/soyart/superwatcher"
	"github.com/soyart/superwatcher/pkg/components"
)

// This demo function calls components.NewDefault, which is the preferred way to init superwatcher for most cases.
func newSuperWatcherPreferred( //nolint:unused
	conf *superwatcher.Config,
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
// to init superwatcher components individuall and customize the components via factory function arguments
func newSuperwatcherFromArgs( //nolint:unused
	conf *superwatcher.Config,
	ethClient superwatcher.EthClient,
	addresses []common.Address,
	topics []common.Hash,
	stateDataGateway superwatcher.StateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
) (superwatcher.Emitter, superwatcher.Engine) {
	syncChan := make(chan struct{})
	resultChan := make(chan *superwatcher.PollerResult)
	errChan := make(chan error)

	emitter := components.NewEmitter(conf, ethClient, stateDataGateway, nil, syncChan, resultChan, errChan)
	emitterClient := components.NewEmitterClient(conf, syncChan, resultChan, errChan)
	engine := components.NewEngine(emitterClient, serviceEngine, stateDataGateway, conf.LogLevel)
	poller := components.NewPoller(nil, nil, conf.DoReorg, conf.DoHeader, conf.FilterRange, ethClient, conf.LogLevel, conf.Policy)

	poller.SetAddresses(addresses)
	poller.SetTopics([][]common.Hash{topics})

	emitter.SetPoller(poller)

	return emitter, engine
}

// newSuperWatcherOptions demonstrates how users can use components.NewSuperWatcherOptions
// to create the top-level wrapper superwatcher.SuperWatcher.
// For simple applications, this might be the most effective and error-free way to initialize superwatcher service.
//
// Note that components.*Options funcs do not validate the option.
func newSuperWatcherOptions( // nolint:unused
	conf *superwatcher.Config,
	ethClient superwatcher.EthClient,
	addresses []common.Address,
	topics []common.Hash,
	stateDataGateway superwatcher.StateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
) superwatcher.SuperWatcher {
	syncChan := make(chan struct{})
	resultChan := make(chan *superwatcher.PollerResult)
	errChan := make(chan error)

	return components.NewSuperWatcherOptions(
		components.WithDoReorg(conf.DoReorg),
		components.WithDoHeader(conf.DoHeader),
		components.WithConfig(conf),
		components.WithAddresses(addresses...),
		components.WithTopics(topics),
		components.WithGetStateDataGateway(stateDataGateway),
		components.WithSetStateDataGateway(stateDataGateway),
		components.WithServiceEngine(serviceEngine),
		components.WithSyncChan(syncChan),
		components.WithFilterResultChan(resultChan),
		components.WithErrChan(errChan),
	)
}
