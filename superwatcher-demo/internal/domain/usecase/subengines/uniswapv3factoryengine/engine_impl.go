package uniswapv3factoryengine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines"
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
		logArtifact, err = e.HandleGoodLog(log)
		if err != nil {
			return nil, errors.Wrapf(err, "poolfactory.HandleGoodLog failed on log txHash %s", log.BlockHash.String())
		}
	}

	// poolArtifact is a map, use one instance returned from HandleGoodLog
	return []superwatcher.Artifact{logArtifact}, nil
}

func (e *uniswapv3PoolFactoryEngine) HandleGoodLog(log *types.Log) (PoolFactoryArtifact, error) {
	artifact := make(PoolFactoryArtifact)
	logEventKey := log.Topics[0]

	for _, event := range e.poolFactoryContract.ContractEvents {
		// This engine is supposed to handle more than 1 event,
		// but it's not yet finished now.
		if logEventKey == event.ID || event.Name == "PoolCreated" {
			pool, err := mapLogToPoolCreated(e.poolFactoryContract.ContractABI, event.Name, log)
			if err != nil {
				return nil, errors.Wrap(err, "failed to map PoolCreated log to domain struct")
			}
			if pool == nil {
				logger.Panic("nil pool mapped - should not happen")
			}
			if err := e.handlePoolCreated(pool); err != nil {
				return nil, errors.Wrap(err, "failed to process poolCreated")
			}

			// Saves engine artifact
			artifact[*pool] = PoolFactoryStateCreated
		}
	}

	return artifact, nil
}

func (e *uniswapv3PoolFactoryEngine) HandleReorgedLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	logger.Debug("poolfactory.HandleReorgedLogs", zap.Any("input artifacts", artifacts))

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

func (e *uniswapv3PoolFactoryEngine) handleReorgedLog(log *types.Log, artifacts []superwatcher.Artifact) (PoolFactoryArtifact, error) {

	var returnArtifacts []superwatcher.Artifact
	logEventKey := log.Topics[0]

	// Find poolFactory artifact here
	var poolArtifact PoolFactoryArtifact
	for _, a := range artifacts {
		pa, ok := a.(PoolFactoryArtifact)
		if !ok {
			logger.Debug("found non-pool artifact")
			continue
		}

		poolArtifact = pa
	}

	for _, event := range e.poolFactoryContract.ContractEvents {
		// This engine is supposed to handle more than 1 event,
		// but it's not yet finished now.
		if logEventKey == event.ID || event.Name == "PoolCreated" {
			pool, err := mapLogToPoolCreated(e.poolFactoryContract.ContractABI, event.Name, log)
			if err != nil {
				return nil, errors.Wrap(err, "failed to map PoolCreated log to domain struct")
			}

			processArtifacts, err := e.handleReorgedPool(pool, poolArtifact)
			if err != nil {
				return nil, errors.Wrap(err, "failed to handle reorged PoolCreated")
			}

			returnArtifacts = append(returnArtifacts, processArtifacts)
		}
	}

	return nil, fmt.Errorf("event topic %s not found", logEventKey)
}

// In uniswapv3poolfactory case, we only revert PoolCreated in the db.
// Other service may need more elaborate HandleReorg.
func (e *uniswapv3PoolFactoryEngine) handleReorgedPool(
	pool *entity.Uniswapv3PoolCreated,
	poolArtifact PoolFactoryArtifact,
) (
	PoolFactoryArtifact,
	error,
) {
	poolState := poolArtifact[*pool]

	switch poolState {
	case PoolFactoryStateCreated:
		if err := e.revertPoolCreated(pool); err != nil {
			return nil, errors.Wrapf(err, "failed to revert poolCreated for pool %s", pool.Address.String())
		}
	}

	poolArtifact[*pool] = PoolFactoryStateNull
	return poolArtifact, nil
}

// Unused by this service
func (e *uniswapv3PoolFactoryEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))
	return nil
}

func (e *uniswapv3PoolFactoryEngine) createPool(pool *entity.Uniswapv3PoolCreated) error {
	logger.Debug("createPool: got pool", zap.Any("pool", pool))

	return nil
}
