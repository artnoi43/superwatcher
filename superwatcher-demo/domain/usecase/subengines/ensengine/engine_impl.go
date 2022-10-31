package ensengine

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

// The sub-engine uses entity.ENS as engine.Artifact

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *ensEngine) HandleGoodLogs(
	logs []*types.Log,
	artifacts []engine.Artifact,
) (
	[]engine.Artifact,
	error,
) {
	logger.Debug("ensengine.HandleGoodLogs: got logs")

	var outArtifacts []engine.Artifact
	for _, log := range logs {
		resultArtifact, err := e.HandleGoodLog(log, artifacts)
		if err != nil {
			return nil, errors.Wrapf(err, "HandleGoodLog failed on log txHash %s", log.BlockHash.String())
		}

		artifacts = append(artifacts, resultArtifact)
		outArtifacts = append(outArtifacts, resultArtifact)
	}

	return outArtifacts, nil
}

func (e *ensEngine) HandleGoodLog(
	log *types.Log,
	artifacts []engine.Artifact,
) (
	ENSArtifact,
	error,
) {
	// Artifact to return
	artifact := ENSArtifact{RegisterBlockNumber: log.BlockNumber}

	logEventKey := log.Topics[0]
	switch log.Address {
	case e.ensRegistrar.Address:
		for _, event := range e.ensRegistrar.ContractEvents {
			// This engine is supposed to handle more than 1 event,
			// but it's not yet finished now.
			if logEventKey == event.ID {
				switch event.Name {
				// New domain registered from both contracts
				case nameRegistered:
					var err error
					resultArtifact, err := e.handleNameRegisteredRegistrar(log, event.Name, nil)
					if err != nil {
						return artifact, errors.Wrap(err, "failed to create new name from log (registrar)")
					}
					artifact = *resultArtifact
				}
			}
		}
	case e.ensController.Address:
		for _, event := range e.ensController.ContractEvents {
			if logEventKey == event.ID {
				switch event.Name {
				case nameRegistered:
					// Previous artifacts
					var prevArtifact *ENSArtifact
					for _, artifact := range artifacts {
						a, ok := artifact.(ENSArtifact)
						if !ok {
							logger.Panic("found non-ENS artifact", zap.String("type", reflect.TypeOf(artifact).String()))
						}
						prevArtifact = &a
					}
					resultArtifact, err := e.handleNameRegisteredController(log, event.Name, prevArtifact)
					if err != nil {
						return artifact, errors.Wrap(err, "failed to create new name from log (controller)")
					}
					artifact = *resultArtifact
				}
			}
		}
	}

	return artifact, nil
}

func (e *ensEngine) HandleReorgedLogs(
	logs []*types.Log,
	artifacts []engine.Artifact,
) (
	[]engine.Artifact,
	error,
) {
	logger.Debug("ensengine.HandleReorgedLogs: got logs", zap.Any("input artifacts", artifacts))

	var outputArtifacts []engine.Artifact
	for _, log := range logs {
		ens, err := e.handleReorgedLog(log, artifacts)
		if err != nil {
			return nil, errors.Wrap(err, "ensEngine.handleReorgedLog failed")
		}

		outputArtifacts = append(outputArtifacts, ens)
	}

	return outputArtifacts, nil
}

// handleReorgedLog examines the log, get log's previous artifact, and handle chain reorg events
func (e *ensEngine) handleReorgedLog(
	log *types.Log,
	artifacts []engine.Artifact,
) (
	ENSArtifact,
	error,
) {

	// Previous artifacts
	var prevArtifact *ENSArtifact
	for _, artifact := range artifacts {
		a, ok := artifact.(ENSArtifact)
		if !ok {
			logger.Panic("found non-ENS artifact", zap.String("type", reflect.TypeOf(artifact).String()))
		}
		if a.TxHash == log.TxHash {
			prevArtifact = &a
		}
	}

	// Return artifact
	var artifact ENSArtifact

	logEventKey := log.Topics[0]
	switch log.Address {
	case e.ensRegistrar.Address:
		for _, event := range e.ensRegistrar.ContractEvents {
			// This engine is supposed to handle more than 1 event,
			// but it's not yet finished now.
			if logEventKey == event.ID {
				switch event.Name {
				case nameRegistered:
					reorgArtifact, err := e.revertNameRegisteredRegistrar(log, event.Name, prevArtifact)
					if err != nil {
						return artifact, errors.Wrap(err, "failed to create new name from log")
					}
					artifact = *reorgArtifact
					return artifact, nil
				}
			}
		}
	case e.ensController.Address:
		for _, event := range e.ensController.ContractEvents {
			// This engine is supposed to handle more than 1 event,
			// but it's not yet finished now.
			if logEventKey == event.ID {
				switch event.Name {
				case nameRegistered:
					reorgArtifact, err := e.revertNameRegisteredController(log, event.Name, prevArtifact)
					if err != nil {
						return artifact, errors.Wrap(err, "failed to create new name from log")
					}
					artifact = *reorgArtifact
					return artifact, nil
				}
			}
		}
	}

	return artifact, fmt.Errorf("event topic %s not found", logEventKey)
}

// Unused by this service
func (e *ensEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))
	return nil
}

func (e *ensEngine) createPool(pool *entity.Uniswapv3PoolCreated) error {
	logger.Debug("createPool: got pool", zap.Any("pool", pool))

	return nil
}
