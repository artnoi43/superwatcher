package thinengine

import (
	"context"

	"github.com/pkg/errors"
)

func (e *thinEngine) handleResults(ctx context.Context) error { //nolint:unused
	for {
		result := e.emitterClient.WatcherResult()
		err := e.serviceEngine.HandleFilterResult(result)
		if err != nil {
			return errors.Wrap(errors.WithStack(err), "serviceEngine returned error")
		}
	}
}
