package engine

import "github.com/pkg/errors"

func (e *engine) handleEmitterError() error {
	e.debugger.Debug(3, "handleError started")
	for {
		err := e.emitterClient.WatcherError()
		if err != nil {
			err = e.serviceEngine.HandleEmitterError(err)
			if err != nil {
				return errors.Wrap(err, "serviceEngine.HandleEmitterError returned non-nil error")
			}

			// Emitter error handled in service without error
			continue
		}

		e.debugger.Debug(3, "got nil error from emitter - should not happen unless errChan was closed")
		break
	}
	return nil
}
