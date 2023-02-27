package thinengine

import (
	"context"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// thinEngine is a thin implementation of superwatcher.Engine.
// It does not manage states for the service, and all handling of PollerResult is
// managed directly by serviceEngine.
type thinEngine struct { //nolint:unused
	emitterClient    superwatcher.EmitterClient
	serviceEngine    superwatcher.ThinServiceEngine
	stateDataGateway superwatcher.SetStateDataGateway

	debug    bool
	debugger *debugger.Debugger
}

func New(
	emitterClient superwatcher.EmitterClient,
	serviceEngine superwatcher.ThinServiceEngine,
	stateDataGateway superwatcher.SetStateDataGateway,
	debugLevel uint8,
) superwatcher.Engine {
	return &thinEngine{
		emitterClient:    emitterClient,
		serviceEngine:    serviceEngine,
		stateDataGateway: stateDataGateway,
		debug:            debugLevel > 0,
		debugger:         debugger.NewDebugger("thinEngine", debugLevel),
	}
}

func (e *thinEngine) shutdown() { // nolint:unused
	// TODO: Should we close Redis or should the service does it?
	// e.stateDataGateway.Shutdown()
	e.emitterClient.Shutdown()
}

func (e *thinEngine) Loop(ctx context.Context) error { // nolint:unused
	go func() {
		defer e.shutdown()

		if err := e.handleResults(ctx); err != nil {
			e.debugger.Debug(1, "engine.run exited", zap.Error(err))
		}
	}()

	return e.handleEmitterError()
}
