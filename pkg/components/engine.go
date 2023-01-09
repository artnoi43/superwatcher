package components

import (
	"github.com/artnoi43/gsl/gslutils"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/engine"
)

func NewEngine(
	emitterClient superwatcher.EmitterClient,
	serviceEngine superwatcher.ServiceEngine,
	stateDataGateway superwatcher.SetStateDataGateway,
	logLevel uint8,
) superwatcher.Engine {
	return engine.New(
		emitterClient,
		serviceEngine,
		stateDataGateway,
		logLevel,
	)
}

// NewEngineWithEmitterClient creates a new superwatcher.Engine, and pair it with an superwatcher.EmitterClient.
// This is the preferred way of creating a new superwatcher.Engine
func NewEngineWithEmitterClient(
	conf *superwatcher.Config,
	serviceEngine superwatcher.ServiceEngine,
	stateDataGateway superwatcher.SetStateDataGateway,
	syncChan chan<- struct{},
	pollResultChan <-chan *superwatcher.PollerResult,
	errChan <-chan error,
) superwatcher.Engine {
	// TODO: Do we still need EmitterClient?
	emitterClient := NewEmitterClient(
		conf,
		syncChan,
		pollResultChan,
		errChan,
	)

	return NewEngine(
		emitterClient,
		serviceEngine,
		stateDataGateway,
		conf.LogLevel,
	)
}

func NewEngineOptions(options ...Option) superwatcher.Engine {
	var c componentConfig
	for _, opt := range options {
		opt(&c)
	}

	emitterClient := NewEmitterClient(
		c.config,
		c.syncChan,
		c.pollResultChan,
		c.errChan,
	)

	return engine.New(
		emitterClient,
		c.serviceEngine,
		c.setStateDataGateway,
		gslutils.Max(c.logLevel, c.config.LogLevel),
	)
}
