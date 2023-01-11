package emitterclient

// This package maybe removed in favor of centralized pkg/components

import (
	"github.com/artnoi43/gsl"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/emitterclient"
)

type config struct {
	conf                *superwatcher.Config
	client              superwatcher.EthClient
	getStateDataGateway superwatcher.GetStateDataGateway
	syncChan            chan<- struct{}
	pollResultChan      <-chan *superwatcher.PollerResult
	errChan             <-chan error
	logLevel            uint8
}

type Option func(*config)

func WithConfig(conf *superwatcher.Config) Option {
	return func(c *config) {
		c.conf = conf
	}
}

func WithEthClient(client superwatcher.EthClient) Option {
	return func(c *config) {
		c.client = client
	}
}

func WithGetStateDataGateway(gateway superwatcher.StateDataGateway) Option {
	return func(c *config) {
		c.getStateDataGateway = gateway
	}
}

func WithSyncChan(syncChan chan<- struct{}) Option {
	return func(c *config) {
		c.syncChan = syncChan
	}
}

func WithFilterResultChan(resultChan <-chan *superwatcher.PollerResult) Option {
	return func(c *config) {
		c.pollResultChan = resultChan
	}
}

func WithErrChan(errChan <-chan error) Option {
	return func(c *config) {
		c.errChan = errChan
	}
}

func WithLogLevel(level uint8) Option {
	return func(c *config) {
		c.logLevel = level
	}
}

func New(options ...Option) superwatcher.EmitterClient {
	var c config
	for _, opt := range options {
		opt(&c)
	}

	return emitterclient.New(
		c.conf,
		c.syncChan,
		c.pollResultChan,
		c.errChan,
		gsl.Max(c.logLevel, c.conf.LogLevel),
	)
}
