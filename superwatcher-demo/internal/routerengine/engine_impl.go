package routerengine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal"
)

func (e *routerEngine) HandleGoodLogs(
	logs []*types.Log,
	artifacts []superwatcher.Artifact, // Ignored
) (
	[]superwatcher.Artifact,
	error,
) {
	// Artifacts to return - we don't know its size
	var retArtifacts []superwatcher.Artifact //nolint:prealloc
	logsMap := e.mapLogsToSubEngine(logs)

	// ensLogs := logsMap[subengines.SubEngineENS]
	for subEngine, logs := range logsMap {
		serviceEngine, ok := e.services[subEngine]
		if !ok {
			return nil, errors.Wrapf(errNoService, "subengine: %s", subEngine.String())
		}

		resultArtifacts, err := serviceEngine.HandleGoodLogs(logs, filterArtifacts(subEngine, artifacts))
		if err != nil {
			if errors.Is(err, internal.ErrNoNeedHandle) {
				e.debugger.Debug(2, "routerEngine: got ErrNoNeedHandle", zap.String("subEngine", subEngine.String()))
				continue
			}
			return nil, errors.Wrapf(err, "subengine %s HandleGoodBlock failed", subEngine.String())
		}

		// Only append non-nil subengine artifacts
		if resultArtifacts != nil {
			retArtifacts = append(retArtifacts, resultArtifacts)
		}
	}

	return retArtifacts, nil
}

func (e *routerEngine) HandleReorgedLogs(
	logs []*types.Log,
	artifacts []superwatcher.Artifact,
) ([]superwatcher.Artifact, error) {
	e.debugger.Debug(2, "HandleReorgedLogs called", zap.Int("len(logs)", len(logs)), zap.Any("artifacts", artifacts))
	logsMap := e.mapLogsToSubEngine(logs)

	var retArtifacts []superwatcher.Artifact //nolint:all Artifacts to return - we dont know the size
	for subEngine, logs := range logsMap {
		serviceEngine, ok := e.services[subEngine]
		if !ok {
			return nil, errors.Wrap(errNoService, "subengine: "+subEngine.String())
		}

		// Aggregate subEngine-specific artifacts
		outputArtifacts, err := serviceEngine.HandleReorgedLogs(logs, filterArtifacts(subEngine, artifacts))
		if err != nil {
			return nil, errors.Wrapf(err, "subengine %s HandleReorgedBlock failed", subEngine.String())
		}

		retArtifacts = append(retArtifacts, outputArtifacts)
	}

	return retArtifacts, nil
}

// Unused by this service
func (e *routerEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))

	return err
}
