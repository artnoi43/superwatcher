package components

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/thinengine"
)

// NewThinEngine returns thinEngine implementation of superwatcher.Engine
func NewThinEngine(
	emitterClient superwatcher.EmitterClient,
	serviceEngine superwatcher.ThinServiceEngine,
	stateDataGateway superwatcher.SetStateDataGateway,
	logLevel uint8,
) superwatcher.Engine {
	return thinengine.New(emitterClient, serviceEngine, stateDataGateway, logLevel)
}

// NewThinEngine creates an EmitterPoller and and an EmitterClient,
// before using the created objects to create new thinEngine implementation
// of superwatcher.Engine. The channels shared by emitterClient and thinEngine
// are created within the functions and kept private by the implementation.
func NewThinEngineWithEmitter(
	conf *superwatcher.Config,
	getStateDataGateway superwatcher.GetStateDataGateway,
	setStateDataGateway superwatcher.SetStateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	client superwatcher.EthClient,
	policy superwatcher.Policy,
	serviceEngine superwatcher.ThinServiceEngine,
) (superwatcher.Emitter, superwatcher.Engine) {
	poller := NewPoller(addresses, topics, conf.DoReorg, conf.DoHeader, conf.FilterRange, client, conf.LogLevel, policy)

	syncChan := make(chan struct{})
	resultChan := make(chan *superwatcher.PollerResult)
	errChan := make(chan error)

	emitter := NewEmitter(conf, client, getStateDataGateway, poller, syncChan, resultChan, errChan)
	emitterClient := NewEmitterClient(conf, syncChan, resultChan, errChan)

	thinEngine := NewThinEngine(emitterClient, serviceEngine, setStateDataGateway, conf.LogLevel)

	return emitter, thinEngine
}
