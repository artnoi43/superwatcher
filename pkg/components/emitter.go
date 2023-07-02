package components

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/soyart/gsl"

	"github.com/soyart/superwatcher"
	"github.com/soyart/superwatcher/internal/emitter"
)

func NewEmitter(
	conf *superwatcher.Config,
	client superwatcher.EthClient,
	stateDataGateway superwatcher.GetStateDataGateway,
	poller superwatcher.EmitterPoller,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	pollResultChan chan<- *superwatcher.PollerResult,
	errChan chan<- error,
) superwatcher.Emitter {
	return emitter.New(
		conf,
		client,
		stateDataGateway,
		poller,
		syncChan,
		pollResultChan,
		errChan,
	)
}

// NewWithPoller returns a new, default Emitter, with a default WatcherPoller.
// It is the preferred way to init a Emitter if you have not implement WatcherPoller yet yourself.
func NewEmitterWithPoller(
	conf *superwatcher.Config,
	client superwatcher.EthClient,
	stateDataGateway superwatcher.GetStateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	pollResultChan chan<- *superwatcher.PollerResult,
	errChan chan<- error,
) superwatcher.Emitter {
	return emitter.New(
		conf,
		client,
		stateDataGateway,
		NewPoller(addresses, topics, conf.DoReorg, conf.DoHeader, conf.FilterRange, client, conf.LogLevel, conf.Policy),
		syncChan,
		pollResultChan,
		errChan,
	)
}

func NewEmitterOptions(options ...Option) superwatcher.Emitter {
	var c componentConfig
	for _, opt := range options {
		opt(&c)
	}

	poller := NewPoller(
		c.addresses,
		c.topics,
		c.doReorg,
		c.doHeader,
		c.filterRange,
		c.ethClient,
		gsl.Max(c.logLevel, c.config.LogLevel),
		c.policy,
	)

	return emitter.New(
		c.config,
		c.ethClient,
		c.getStateDataGateway,
		poller,
		c.syncChan,
		c.pollResultChan,
		c.errChan,
	)
}
