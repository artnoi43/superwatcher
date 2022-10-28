package demoengine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
)

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *demoEngine) HandleGoodLogs(
	logs []*types.Log,
	artifacts []engine.Artifact, // Ignore
) (
	[]engine.Artifact,
	error,
) {
	logsMap := e.mapLogsToSubEngine(logs)

	var retArtifacts []engine.Artifact
	for subEngine, logs := range logsMap {
		serviceEngine, ok := e.services[subEngine]
		if !ok {
			return nil, errors.Wrapf(errNoService, "subengine: %s", subEngine.String())
		}

		subArtifacts, err := serviceEngine.HandleGoodLogs(logs, artifacts)
		if err != nil {
			return nil, errors.Wrapf(err, "subengine %s HandleGoodBlock failed", subEngine.String())
		}

		retArtifacts = append(retArtifacts, subArtifacts)
	}

	return retArtifacts, nil
}

func (e *demoEngine) HandleReorgedLogs(
	logs []*types.Log,
	artifacts []engine.Artifact,

) ([]engine.Artifact, error) {
	logsMap := e.mapLogsToSubEngine(logs)

	var retArtifacts []engine.Artifact
	for subEngine, logs := range logsMap {
		serviceEngine, ok := e.services[subEngine]
		if !ok {
			return nil, errors.Wrapf(errNoService, "subengine", subEngine.String())
		}

		subArtifacts, err := serviceEngine.HandleReorgedLogs(logs, artifacts)
		if err != nil {
			return nil, errors.Wrapf(err, "subengine %s HandleReorgedBlock failed", subEngine.String())
		}

		retArtifacts = append(retArtifacts, subArtifacts)
	}

	return retArtifacts, nil
}

// Unused by this service
func (e *demoEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))

	return nil
}
