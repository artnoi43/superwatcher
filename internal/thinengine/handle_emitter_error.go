package thinengine

import "github.com/soyart/superwatcher/internal/emitterclient"

func (e *thinEngine) handleEmitterError() error { // nolint:unused
	return emitterclient.HandleEmitterError(e.emitterClient, e.serviceEngine, e.debugger) //nolint:wrapcheck
}
