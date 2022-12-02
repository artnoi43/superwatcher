package engine

import (
	"context"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

type engine struct {
	emitterClient    superwatcher.EmitterClient       // Interfaces with emitter
	stateDataGateway superwatcher.SetStateDataGateway // Saves lastRecordedBlock to Redis
	metadataTracker  MetadataTracker                  // Engine internal state machine

	serviceEngine superwatcher.ServiceEngine // Injected service code

	debugger *debugger.Debugger
	debug    bool // In case we need to debug within a loop with multiple
}

// newWatcherEngine returns default implementation of WatcherEngine
func New(
	client superwatcher.EmitterClient,
	serviceEngine superwatcher.ServiceEngine,
	stateDataGateway superwatcher.SetStateDataGateway,
	logLevel uint8,
) superwatcher.WatcherEngine {
	debug := logLevel > 0

	return &engine{
		emitterClient:    client,
		serviceEngine:    serviceEngine,
		stateDataGateway: stateDataGateway,
		metadataTracker:  NewTracker(logLevel),
		debugger:         debugger.NewDebugger("engine", logLevel),
		debug:            debug,
	}
}

func (e *engine) Loop(ctx context.Context) error {
	go func() {
		defer e.shutdown()

		if err := e.handleResults(ctx); err != nil {
			e.debugger.Debug(1, "engine.run exited", zap.Error(err))
		}
	}()

	return e.handleEmitterError()
}

// shutdown is not exported, and the user of the engine should not attempt to call it.
func (e *engine) shutdown() {
	// TODO: Should we close Redis or should the service does it?
	// e.stateDataGateway.Shutdown()
	e.emitterClient.Shutdown()
}
