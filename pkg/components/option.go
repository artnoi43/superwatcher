package components

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/soyart/superwatcher"
)

type componentConfig struct {
	// poller        superwatcher.EmitterPoller
	// emitter       superwatcher.Emitter
	// emitterClient superwatcher.EmitterClient
	// engine        superwatcher.Engine

	config              *superwatcher.Config
	serviceEngine       superwatcher.ServiceEngine
	ethClient           superwatcher.EthClient
	addresses           []common.Address
	topics              [][]common.Hash
	doReorg             bool
	doHeader            bool
	filterRange         uint64
	policy              superwatcher.Policy
	logLevel            uint8 // redundant in conf, but users may want to set this separately
	syncChan            chan struct{}
	pollResultChan      chan *superwatcher.PollerResult
	errChan             chan error
	getStateDataGateway superwatcher.GetStateDataGateway
	setStateDataGateway superwatcher.SetStateDataGateway
}

type Option func(*componentConfig)

func WithConfig(conf *superwatcher.Config) Option {
	return func(c *componentConfig) {
		c.config = conf
	}
}

func WithEthClient(client superwatcher.EthClient) Option {
	return func(c *componentConfig) {
		c.ethClient = client
	}
}

func WithSyncChan(syncChan chan struct{}) Option {
	return func(c *componentConfig) {
		c.syncChan = syncChan
	}
}

func WithFilterResultChan(resultChan chan *superwatcher.PollerResult) Option {
	return func(c *componentConfig) {
		c.pollResultChan = resultChan
	}
}

func WithErrChan(errChan chan error) Option {
	return func(c *componentConfig) {
		c.errChan = errChan
	}
}

func WithServiceEngine(service superwatcher.ServiceEngine) Option {
	return func(c *componentConfig) {
		c.serviceEngine = service
	}
}

func WithGetStateDataGateway(gateway superwatcher.GetStateDataGateway) Option {
	return func(c *componentConfig) {
		c.getStateDataGateway = gateway
	}
}

func WithSetStateDataGateway(gateway superwatcher.SetStateDataGateway) Option {
	return func(c *componentConfig) {
		c.setStateDataGateway = gateway
	}
}

func WithLogLevel(level uint8) Option {
	return func(c *componentConfig) {
		c.logLevel = level
	}
}

func WithFilterRange(filterRange uint64) Option {
	return func(c *componentConfig) {
		c.filterRange = filterRange
	}
}

func WithAddresses(addresses ...common.Address) Option {
	return func(c *componentConfig) {
		c.addresses = addresses
	}
}

func WithTopics(topics ...[]common.Hash) Option {
	return func(c *componentConfig) {
		c.topics = topics
	}
}

func WithDoReorg(doReorg bool) Option {
	return func(c *componentConfig) {
		c.doReorg = doReorg
	}
}

func WithDoHeader(doHeader bool) Option {
	return func(c *componentConfig) {
		c.doHeader = doHeader
	}
}

func WithPolicy(level superwatcher.Policy) Option {
	return func(c *componentConfig) {
		c.policy = level
	}
}

// func WithEngine(engine superwatcher.Engine) Option {
// 	return func(c *initConfig) {
// 		c.engine = engine
// 	}
// }
//
// func WithEmitter(emitter superwatcher.Emitter) Option {
// 	return func(c *initConfig) {
// 		c.emitter = emitter
// 	}
// }
//
// func WithEmitterClient(client superwatcher.EmitterClient) Option {
// 	return func(c *initConfig) {
// 		c.emitterClient = client
// 	}
// }
//
// func WithEmitterPoller(poller superwatcher.EmitterPoller) Option { //nolint:revive
// 	return func(c *initConfig) {
// 		c.poller = poller
// 	}
// }
