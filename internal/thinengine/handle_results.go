package thinengine

import "context"

func (e *thinEngine) handleResults(ctx context.Context) error {
	for {
		result := e.emitterClient.WatcherResult()
		err := e.serviceEngine.HandleFilterResult(result)

		if err != nil {
			return err
		}
	}
}
