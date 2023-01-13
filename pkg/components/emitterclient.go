package components

import (
	"github.com/artnoi43/gsl"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/emitterclient"
)

func NewEmitterClient(
	conf *superwatcher.Config,
	syncChan chan<- struct{},
	pollResultChan <-chan *superwatcher.PollerResult,
	errChan <-chan error,
) superwatcher.EmitterClient {
	return emitterclient.New(
		conf,
		syncChan,
		pollResultChan,
		errChan,
		conf.LogLevel,
	)
}

func NewEmitterClientOptions(options ...Option) superwatcher.EmitterClient {
	var c componentConfig
	for _, opt := range options {
		opt(&c)
	}

	return emitterclient.New(
		c.config,
		c.syncChan,
		c.pollResultChan,
		c.errChan,
		gsl.Max(c.logLevel, c.config.LogLevel),
	)
}
