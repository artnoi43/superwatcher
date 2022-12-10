package thinengine

import "github.com/artnoi43/superwatcher/internal/emitterclient"

func (e *thinEngine) handleEmitterError() error {
	return emitterclient.HandleEmitterError(e.emitterClient, e.serviceEngine, e.debugger)
}
