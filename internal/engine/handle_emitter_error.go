package engine

import "github.com/soyart/superwatcher/internal/emitterclient"

func (e *engine) handleEmitterError() error {
	return emitterclient.HandleEmitterError(e.emitterClient, e.serviceEngine, e.debugger) //nolint:wrapcheck
}
