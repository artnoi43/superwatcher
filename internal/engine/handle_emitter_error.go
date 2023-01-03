package engine

import "github.com/artnoi43/superwatcher/internal/emitterclient"

func (e *engine) handleEmitterError() error {
	return emitterclient.HandleEmitterError(e.emitterClient, e.serviceEngine, e.debugger)
}
