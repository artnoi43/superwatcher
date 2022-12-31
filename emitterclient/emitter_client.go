package emitterclient

import (
	"github.com/artnoi43/gsl/gslutils"

	"github.com/artnoi43/superwatcher"
	spwconf "github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitterclient"
)

type config struct {
	conf                *spwconf.Config
	client              superwatcher.EthClient
	getStateDataGateway superwatcher.GetStateDataGateway
	syncChan            chan<- struct{}
	filterResultChan    <-chan *superwatcher.FilterResult
	errChan             <-chan error
	logLevel            uint8
}

type Option func(*config)

func WithConfig(conf *spwconf.Config) Option {
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

func WithFilterResultChan(resultChan <-chan *superwatcher.FilterResult) Option {
	return func(c *config) {
		c.filterResultChan = resultChan
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
		c.filterResultChan,
		c.errChan,
		gslutils.Max(c.logLevel, c.conf.LogLevel),
	)
}
