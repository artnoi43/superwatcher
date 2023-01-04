package uniswapv3factoryengine

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines"
)

// This sub-engine uses a hash map as superwatcher.Artifact. Because it's a hash map,
// only 1 instance of this struct is needed for all the lofs
type PoolFactoryArtifact map[entity.Uniswapv3PoolCreated]uniswapv3PoolFactoryState

// Implements demoengine.demoArtifact
func (a PoolFactoryArtifact) ForSubEngine() subengines.SubEngineEnum {
	return subengines.SubEngineUniswapv3Factory
}

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *uniswapv3PoolFactoryEngine) HandleGoodBlocks(
	blocks []*superwatcher.BlockInfo,
	artifacts []superwatcher.Artifact,
) (
	map[common.Hash][]superwatcher.Artifact,
	error,
) {
	retArtifacts := make(map[common.Hash][]superwatcher.Artifact)
	for _, block := range blocks {
		// New artifact is created for new logs
		for _, log := range block.Logs {
			logArtifact, err := e.handleGoodLog(log)
			if err != nil {
				return nil, errors.Wrapf(err, "poolfactory.HandleGoodLog failed on log txHash %s", log.BlockHash.String())
			}

			retArtifacts[log.BlockHash] = append(retArtifacts[log.BlockHash], logArtifact)
		}
	}

	// poolArtifact is a map, use one instance returned from HandleGoodLog
	return retArtifacts, nil
}

func (e *uniswapv3PoolFactoryEngine) HandleReorgedBlocks(
	blocks []*superwatcher.BlockInfo,
	artifacts []superwatcher.Artifact,
) (
	map[common.Hash][]superwatcher.Artifact,
	error,
) {
	e.debugger.Debug(1, "poolfactory.HandleReorgedLogs", zap.Any("input artifacts", artifacts))

	retArtifacts := make(map[common.Hash][]superwatcher.Artifact)
	for _, block := range blocks {
		for _, log := range block.Logs {
			logArtifact, err := e.handleReorgedLog(log, artifacts)
			if err != nil {
				return nil, errors.Wrap(err, "uniswapv3PoolFactoryEngine.handleReorgedLog failed")
			}

			retArtifacts[log.BlockHash] = append(retArtifacts[log.BlockHash], logArtifact)
		}
	}

	return retArtifacts, nil
}

func (e *uniswapv3PoolFactoryEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))
	return err
}
