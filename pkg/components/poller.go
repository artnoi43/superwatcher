package components

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/soyart/gsl"

	"github.com/soyart/superwatcher"
	"github.com/soyart/superwatcher/internal/poller"
)

func NewPoller(
	addresses []common.Address,
	topics [][]common.Hash,
	doReorg bool,
	doHeader bool,
	filterRange uint64,
	client superwatcher.EthClient,
	logLevel uint8,
	policy superwatcher.Policy,
) superwatcher.EmitterPoller {
	return poller.New(
		addresses,
		topics,
		doReorg,
		doHeader,
		filterRange,
		client,
		logLevel,
		policy,
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
		gsl.Max(c.logLevel, c.config.LogLevel),
		gsl.Max(c.policy, c.config.Policy),
	)
}
