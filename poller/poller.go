package poller

// This package maybe removed in favor of centralized pkg/components

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/poller"
)

type config struct {
	doReorg     bool
	doHeader    bool
	logLevel    uint8
	policy      superwatcher.Policy
	filterRange uint64
	client      superwatcher.EthClient
	addresses   []common.Address
	topics      [][]common.Hash
}

type Option func(*config)

func WithLogLevel(level uint8) Option {
	return func(c *config) {
		c.logLevel = level
	}
}

func WithFilterRange(filterRange uint64) Option {
	return func(c *config) {
		c.filterRange = filterRange
	}
}

func WithEthClient(client superwatcher.EthClient) Option {
	return func(c *config) {
		c.client = client
	}
}

func WithAddresses(addresses ...common.Address) Option {
	return func(c *config) {
		c.addresses = addresses
	}
}

func WithTopics(topics ...[]common.Hash) Option {
	return func(c *config) {
		c.topics = topics
	}
}

func WithDoReorg(doReorg bool) Option {
	return func(c *config) {
		c.doReorg = doReorg
	}
}

func WithDoHeader(doHeader bool) Option {
	return func(c *config) {
		c.doHeader = doHeader
	}
}

func WithPolicy(level superwatcher.Policy) Option {
	return func(c *config) {
		c.policy = level
	}
}

func New(options ...Option) superwatcher.EmitterPoller {
	var c config
	for _, opt := range options {
		opt(&c)
	}

	return poller.New(
		c.addresses,
		c.topics,
		c.doReorg,
		c.doHeader,
		c.filterRange,
		c.client,
		c.logLevel,
		c.policy,
	)
}
