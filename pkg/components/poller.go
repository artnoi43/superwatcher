package components

import (
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/poller"
)

func NewPoller(
	addresses []common.Address,
	topics [][]common.Hash,
	doReorg bool,
	doHeader bool,
	filterRange uint64,
	client superwatcher.EthClient,
	logLevel uint8,
	pollLevel superwatcher.PollLevel,
) superwatcher.EmitterPoller {
	return poller.New(
		addresses,
		topics,
		doReorg,
		doHeader,
		filterRange,
		client,
		logLevel,
		pollLevel,
	)
}

func NewPollerOptions(options ...Option) superwatcher.EmitterPoller {
	var c componentConfig
	for _, opt := range options {
		opt(&c)
	}

	return poller.New(
		c.addresses,
		c.topics,
		c.doReorg,
		c.doHeader,
		c.filterRange,
		c.ethClient,
		gslutils.Max(c.logLevel, c.config.LogLevel),
		gslutils.Max(c.pollLevel, c.config.PollLevel),
	)
}
