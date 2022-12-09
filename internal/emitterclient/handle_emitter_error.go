package emitterclient

import (
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

func HandleEmitterError(
	client superwatcher.EmitterClient,
	serviceEngine superwatcher.BaseServiceEngine,
	debugger *debugger.Debugger,
) error {
	for {
		err := client.WatcherError()
		if err != nil {
			err = serviceEngine.HandleEmitterError(err)
			if err != nil {
				return errors.Wrap(err, "serviceEngine.HandleEmitterError returned non-nil error")
			}

			// Emitter error handled in service without error
			continue
		}

		debugger.Debug(3, "got nil error from emitter - should not happen unless errChan was closed")
		break
	}

	return nil
}
