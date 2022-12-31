package components

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/poller"
)

func NewPoller(
	addresses []common.Address,
	topics [][]common.Hash,
	doReorg bool,
	filterRange uint64,
	client superwatcher.EthClient,
	logLevel uint8,
) superwatcher.EmitterPoller {
	return poller.New(
		addresses,
		topics,
		doReorg,
		filterRange,
		client,
		logLevel,
	)
}

func NewPollerOptions(options ...Option) superwatcher.EmitterPoller {
	var c initConfig
	for _, opt := range options {
		opt(&c)
	}

	return poller.New(
		c.addresses,
		c.topics,
		c.doReorg,
		c.filterRange,
		c.ethClient,
		gslutils.Max(c.logLevel, c.conf.LogLevel),
	)
}
