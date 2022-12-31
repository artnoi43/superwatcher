package components

import (
	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitterclient"
)

func NewEmitterClient(
	conf *config.Config,
	syncChan chan<- struct{},
	pollResultChan <-chan *superwatcher.PollResult,
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
	var c initConfig
	for _, opt := range options {
		opt(&c)
	}

	return emitterclient.New(
		c.conf,
		c.syncChan,
		c.pollResultChan,
		c.errChan,
		gslutils.Max(c.logLevel, c.conf.LogLevel),
	)
}
