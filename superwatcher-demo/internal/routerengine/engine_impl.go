package routerengine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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
	map[common.Hash][]superwatcher.Artifact,
	error,
) {
	// Artifacts to return - we don't know its size
	retArtifacts := make(map[common.Hash][]superwatcher.Artifact) //nolint:prealloc
	logsMap := e.mapLogsToSubEngine(logs)

	// ensLogs := logsMap[subengines.SubEngineENS]
	for subEngine, logs := range logsMap {
		serviceEngine, ok := e.Services[subEngine]
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
			mergeArtifacts(retArtifacts, resultArtifacts)
		}
	}

	return retArtifacts, nil
}

func (e *routerEngine) HandleReorgedLogs(
	logs []*types.Log,
	artifacts []superwatcher.Artifact,
) (
	map[common.Hash][]superwatcher.Artifact,
	error,
) {
	e.debugger.Debug(
		2, fmt.Sprintf("got %d reorged logs and %d artifacts", len(logs), len(artifacts)),
		zap.Any("artifacts", artifacts),
	)

	var retArtifacts = make(map[common.Hash][]superwatcher.Artifact) //nolint:all Artifacts to return - we dont know the size
	logsMap := e.mapLogsToSubEngine(logs)

	for subEngine, logs := range logsMap {
		serviceEngine, ok := e.Services[subEngine]
		if !ok {
			return nil, errors.Wrap(errNoService, "subengine: "+subEngine.String())
		}

		// Aggregate subEngine-specific artifacts
		resultArtifacts, err := serviceEngine.HandleReorgedLogs(logs, filterArtifacts(subEngine, artifacts))
		if err != nil {
			return nil, errors.Wrapf(err, "subengine %s HandleReorgedBlock failed", subEngine.String())
		}

		if resultArtifacts != nil {
			mergeArtifacts(retArtifacts, resultArtifacts)
		}
	}

	return retArtifacts, nil
}

// Unused by this service
func (e *routerEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))

	return err
}
