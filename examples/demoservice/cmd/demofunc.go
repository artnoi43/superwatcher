package main

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/emitter"
	"github.com/artnoi43/superwatcher/engine"
	"github.com/artnoi43/superwatcher/pkg/components"
	"github.com/artnoi43/superwatcher/poller"
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
// to init superwatcher components individually
func newSuperwatcherAdvanced( //nolint:unused
	conf *superwatcher.Config,
	ethClient superwatcher.EthClient,
	addresses []common.Address,
	topics []common.Hash,
	stateDataGateway superwatcher.StateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
) (superwatcher.Emitter, superwatcher.Engine) {
	syncChan := make(chan struct{})
	resultChan := make(chan *superwatcher.PollResult)
	errChan := make(chan error)

	emitter := components.NewEmitter(conf, ethClient, stateDataGateway, nil, syncChan, resultChan, errChan)
	emitterClient := components.NewEmitterClient(conf, syncChan, resultChan, errChan)
	engine := components.NewEngine(emitterClient, serviceEngine, stateDataGateway, conf.LogLevel)
	poller := components.NewPoller(nil, nil, conf.DoReorg, conf.DoHeader, conf.FilterRange, ethClient, conf.LogLevel)

	poller.SetAddresses(addresses)
	poller.SetTopics([][]common.Hash{topics})

	emitter.SetPoller(poller)

	return emitter, engine
}

// This demo function demonstrates how users can use OptionFunc to initiate superwatcher
func newSuperwacherSoyV1( //nolint:unused
	conf *superwatcher.Config,
	ethClient superwatcher.EthClient,
	addresses []common.Address,
	topics []common.Hash,
	stateDataGateway superwatcher.StateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
) (superwatcher.Emitter, superwatcher.Engine) {
	poller := poller.New(
		poller.WithDoReorg(conf.DoReorg),
		poller.WithDoHeader(conf.DoHeader),
		poller.WithEthClient(ethClient),
		poller.WithFilterRange(conf.FilterRange),
		poller.WithAddresses(addresses...),
		poller.WithTopics(topics),
		poller.WithLogLevel(conf.LogLevel),
	)

	syncChan := make(chan struct{})
	resultChan := make(chan *superwatcher.PollResult)
	errChan := make(chan error)

	emitter := emitter.New(
		emitter.WithConfig(conf),
		emitter.WithEmitterPoller(poller),
		emitter.WithEthClient(ethClient),
		emitter.WithGetStateDataGateway(stateDataGateway),
		emitter.WithSyncChan(syncChan),
		emitter.WithFilterResultChan(resultChan),
		emitter.WithErrChan(errChan),
	)

	engine := engine.New(
		engine.WithEmitterClient(nil),
		engine.WithServiceEngine(serviceEngine),
		engine.WithSetStateDataGateway(stateDataGateway),
		engine.WithLogLevel(conf.LogLevel),
	)

	return emitter, engine
}

func newSuperWatcherSoyV2( // nolint:unused
	conf *superwatcher.Config,
	ethClient superwatcher.EthClient,
	addresses []common.Address,
	topics []common.Hash,
	stateDataGateway superwatcher.StateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
) superwatcher.SuperWatcher {
	syncChan := make(chan struct{})
	resultChan := make(chan *superwatcher.PollResult)
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
