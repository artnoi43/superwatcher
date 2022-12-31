package components

import (
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter"
)

func NewEmitter(
	conf *config.Config,
	client superwatcher.EthClient,
	stateDataGateway superwatcher.GetStateDataGateway,
	poller superwatcher.EmitterPoller,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	filterResultChan chan<- *superwatcher.FilterResult,
	errChan chan<- error,
) superwatcher.Emitter {
	return emitter.New(
		conf,
		client,
		stateDataGateway,
		poller,
		syncChan,
		filterResultChan,
		errChan,
	)
}

// NewWithPoller returns a new, default WatcherEmitter, with a default WatcherPoller.
// It is the preferred way to init a WatcherEmitter if you have not implement WatcherPoller yet yourself.
func NewEmitterWithPoller(
	conf *config.Config,
	client superwatcher.EthClient,
	stateDataGateway superwatcher.GetStateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	filterResultChan chan<- *superwatcher.FilterResult,
	errChan chan<- error,
) superwatcher.Emitter {
	return emitter.New(
		conf,
		client,
		stateDataGateway,
		NewPoller(addresses, topics, conf.DoReorg, conf.FilterRange, client, conf.LogLevel),
		syncChan,
		filterResultChan,
		errChan,
	)
}

func NewEmitterOptions(options ...Option) superwatcher.Emitter {
	var c initConfig
	for _, opt := range options {
		opt(&c)
	}

	poller := NewPoller(
		c.addresses,
		c.topics,
		c.doReorg,
		c.filterRange,
		c.ethClient,
		gslutils.Max(c.logLevel, c.conf.LogLevel),
	)

	return emitter.New(
		c.conf,
		c.ethClient,
		c.getStateDataGateway,
		poller,
		c.syncChan,
		c.filterResultChan,
		c.errChan,
	)
}
