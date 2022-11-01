package engine

import (
	"context"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/internal/domain/usecase/emitterclient"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/logger/debug"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"
)

type engine struct {
	emitterClient    emitterclient.Client          // Interfaces with emitter
	stateDataGateway watcherstate.StateDataGateway // Saves lastRecordedBlock to Redis
	metadataTracker  MetadataTracker               // Engine internal state machine

	serviceEngine superwatcher.ServiceEngine // Injected service code

	debug bool
}

func (e *engine) Loop(ctx context.Context) error {
	go func() {
		if err := e.handleLogs(ctx); err != nil {
			e.debugMsg("*engine.run exited", zap.Error(err))
		}

		defer e.shutdown()
	}()

	return e.handleEmitterError()
}

// shutdown is not exported, and the user of the engine should not attempt to call it.
func (e *engine) shutdown() {
	// TODO: Should we close Redis or should the service does it?
	// e.stateDataGateway.Shutdown()
	e.emitterClient.Shutdown()
}

func (e *engine) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(e.debug, msg, fields...)
}
