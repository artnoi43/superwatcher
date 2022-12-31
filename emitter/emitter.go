package emitter

import (
	"github.com/artnoi43/superwatcher"
	spwconf "github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter"
)

type config struct {
	conf                *spwconf.Config
	poller              superwatcher.EmitterPoller
	ethClient           superwatcher.EthClient
	getStateDataGateway superwatcher.GetStateDataGateway
	syncChan            <-chan struct{}
	filterResultChan    chan<- *superwatcher.FilterResult
	errChan             chan<- error
}

type Option func(*config)

func WithConfig(conf *spwconf.Config) Option {
	return func(c *config) {
		c.conf = conf
	}
}

func WithEmitterPoller(poller superwatcher.EmitterPoller) Option { //nolint:revive
	return func(c *config) {
		c.poller = poller
	}
}

func WithEthClient(client superwatcher.EthClient) Option {
	return func(c *config) {
		c.ethClient = client
	}
}

func WithGetStateDataGateway(gateway superwatcher.StateDataGateway) Option {
	return func(c *config) {
		c.getStateDataGateway = gateway
	}
}

func WithSyncChan(syncChan <-chan struct{}) Option {
	return func(c *config) {
		c.syncChan = syncChan
	}
}

func WithFilterResultChan(resultChan chan<- *superwatcher.FilterResult) Option {
	return func(c *config) {
		c.filterResultChan = resultChan
	}
}

func WithErrChan(errChan chan<- error) Option {
	return func(c *config) {
		c.errChan = errChan
	}
}

func New(options ...Option) superwatcher.Emitter {
	var c config
	for _, opt := range options {
		opt(&c)
	}

	return emitter.New(
		c.conf,
		c.ethClient,
		c.getStateDataGateway,
		c.poller,
		c.syncChan,
		c.filterResultChan,
		c.errChan,
	)
}
