package uniswapv3factoryengine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
)

// This sub-engine uses a hash map as superwatcher.Artifact. Because it's a hash map,
// only 1 instance of this struct is needed for all the lofs
type PoolFactoryArtifact map[entity.Uniswapv3PoolCreated]uniswapv3PoolFactoryState

// Implements demoengine.demoArtifact
func (a PoolFactoryArtifact) ForSubEngine() subengines.SubEngineEnum {
	return subengines.SubEngineUniswapv3Factory
}

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *uniswapv3PoolFactoryEngine) HandleGoodLogs(
	logs []*types.Log,
	artifacts []superwatcher.Artifact,
) (
	[]superwatcher.Artifact,
	error,
) {
	// New artifact is created for new logs
	var logArtifact PoolFactoryArtifact
	var err error
	for _, log := range logs {
		logArtifact, err = e.handleGoodLog(log)
		if err != nil {
			return nil, errors.Wrapf(err, "poolfactory.HandleGoodLog failed on log txHash %s", log.BlockHash.String())
		}
	}

	// poolArtifact is a map, use one instance returned from HandleGoodLog
	return []superwatcher.Artifact{logArtifact}, nil
}

func (e *uniswapv3PoolFactoryEngine) HandleReorgedLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	e.debugger.Debug(1, "poolfactory.HandleReorgedLogs", zap.Any("input artifacts", artifacts))

	var logArtifact PoolFactoryArtifact
	var err error
	for _, log := range logs {
		logArtifact, err = e.handleReorgedLog(log, artifacts)
		if err != nil {
			return nil, errors.Wrap(err, "uniswapv3PoolFactoryEngine.handleReorgedLog failed")
		}

	}

	return []superwatcher.Artifact{logArtifact}, nil
}

func (e *uniswapv3PoolFactoryEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))
	return err
}
