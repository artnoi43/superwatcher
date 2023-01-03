package components

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
)

type initConfig struct {
	// poller        superwatcher.EmitterPoller
	// emitter       superwatcher.Emitter
	// emitterClient superwatcher.EmitterClient
	// engine        superwatcher.Engine

	conf                *config.Config
	serviceEngine       superwatcher.ServiceEngine
	ethClient           superwatcher.EthClient
	addresses           []common.Address
	topics              [][]common.Hash
	doReorg             bool
	filterRange         uint64
	logLevel            uint8 // redundant in conf, but users may want to set this separately
	syncChan            chan struct{}
	pollResultChan      chan *superwatcher.PollResult
	errChan             chan error
	getStateDataGateway superwatcher.GetStateDataGateway
	setStateDataGateway superwatcher.SetStateDataGateway
}

type Option func(*initConfig)

func WithConfig(conf *config.Config) Option {
	return func(c *initConfig) {
		c.conf = conf
	}
}

func WithEthClient(client superwatcher.EthClient) Option {
	return func(c *initConfig) {
		c.ethClient = client
	}
}

func WithSyncChan(syncChan chan struct{}) Option {
	return func(c *initConfig) {
		c.syncChan = syncChan
	}
}

func WithFilterResultChan(resultChan chan *superwatcher.PollResult) Option {
	return func(c *initConfig) {
		c.pollResultChan = resultChan
	}
}

func WithErrChan(errChan chan error) Option {
	return func(c *initConfig) {
		c.errChan = errChan
	}
}

func WithServiceEngine(service superwatcher.ServiceEngine) Option {
	return func(c *initConfig) {
		c.serviceEngine = service
	}
}

func WithGetStateDataGateway(gateway superwatcher.GetStateDataGateway) Option {
	return func(c *initConfig) {
		c.getStateDataGateway = gateway
	}
}

func WithSetStateDataGateway(gateway superwatcher.SetStateDataGateway) Option {
	return func(c *initConfig) {
		c.setStateDataGateway = gateway
	}
}

func WithLogLevel(level uint8) Option {
	return func(c *initConfig) {
		c.logLevel = level
	}
}

func WithFilterRange(filterRange uint64) Option {
	return func(c *initConfig) {
		c.filterRange = filterRange
	}
}

func WithAddresses(addresses ...common.Address) Option {
	return func(c *initConfig) {
		c.addresses = addresses
	}
}

func WithTopics(topics ...[]common.Hash) Option {
	return func(c *initConfig) {
		c.topics = topics
	}
}

func WithDoReorg(doReorg bool) Option {
	return func(c *initConfig) {
		c.doReorg = doReorg
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
