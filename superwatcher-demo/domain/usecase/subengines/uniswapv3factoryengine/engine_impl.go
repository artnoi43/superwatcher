package uniswapv3factoryengine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *uniswapv3PoolFactoryEngine) HandleGoodLogs(
	logs []*types.Log,
	artifacts []engine.Artifact, // Ignore
) (
	[]engine.Artifact,
	error,
) {
	logger.Debug("poolfactory.HandleGoodLog", zap.Any("input artifacts", artifacts))

	var logArtifact poolFactoryArtifact
	var err error
	for _, log := range logs {
		logArtifact, err = e.HandleGoodLog(log)
		if err != nil {
			return nil, errors.Wrapf(err, "HandleGoodLog failed on log txHash %s", log.BlockHash.String())
		}
	}

	return []engine.Artifact{logArtifact}, nil
}

func (e *uniswapv3PoolFactoryEngine) HandleGoodLog(log *types.Log) (poolFactoryArtifact, error) {
	artifact := make(poolFactoryArtifact)
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

func (e *uniswapv3PoolFactoryEngine) HandleReorgedLogs(logs []*types.Log, artifacts []engine.Artifact) ([]engine.Artifact, error) {
	logger.Debug("poolfactory.HandleReorgedLogs", zap.Any("input artifacts", artifacts))

	var logArtifact poolFactoryArtifact
	var err error
	for _, log := range logs {
		logArtifact, err = e.handleReorgedLog(log, artifacts)
		if err != nil {
			return nil, errors.Wrap(err, "uniswapv3PoolFactoryEngine.handleReorgedLog failed")
		}
	}

	return []engine.Artifact{logArtifact}, nil
}

func (e *uniswapv3PoolFactoryEngine) handleReorgedLog(log *types.Log, artifacts []engine.Artifact) (poolFactoryArtifact, error) {

	var returnArtifacts []engine.Artifact
	logEventKey := log.Topics[0]

	// Find poolFactory artifact here
	var poolArtifact poolFactoryArtifact
	for _, a := range artifacts {
		pa, ok := a.(poolFactoryArtifact)
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
	poolArtifact poolFactoryArtifact,
) (
	poolFactoryArtifact,
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
