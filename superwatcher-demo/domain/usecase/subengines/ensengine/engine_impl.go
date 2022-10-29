package ensengine

import (
	"fmt"

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
) (
	[]engine.Artifact,
	error,
) {
	logger.Debug("poolfactory.HandleGoodLog: got logs")

	var artifacts []ENSArtifact
	for _, log := range logs {
		resultArtifact, err := e.HandleGoodLog(log, artifacts)
		if err != nil {
			return nil, errors.Wrapf(err, "HandleGoodLog failed on log txHash %s", log.BlockHash.String())
		}

		artifacts = append(artifacts, resultArtifact)
	}

	return []engine.Artifact{artifacts}, nil
}

func (e *ensEngine) HandleGoodLog(
	log *types.Log,
	artifacts []ENSArtifact,
) (
	ENSArtifact,
	error,
) {
	logEventKey := log.Topics[0]

	artifact := ENSArtifact{BlockNumber: log.BlockNumber}

	switch log.Address {
	case e.ensRegistrar.Address:
		for _, event := range e.ensRegistrar.ContractEvents {
			// This engine is supposed to handle more than 1 event,
			// but it's not yet finished now.
			if logEventKey == event.ID {
				switch event.Name {
				// New domain registered from both contracts
				case "NameRegistered":
					var err error
					resultArtifact, err := e.handleNameRegisteredRegistrar(log, event.Name, artifacts)
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
				case "NameRegistered":
					resultArtifact, err := e.handleNameRegisteredController(log, event.Name, artifacts)
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
	logger.Debug("poolfactory.HandleReorgedLogs", zap.Any("input artifacts", artifacts))

	var outputArtifacts []engine.Artifact
	for _, log := range logs {
		ens, err := e.handleReorgedLog(log, artifacts)
		if err != nil {
			return nil, errors.Wrap(err, "ensEngine.handleReorgedLog failed")
		}

		outputArtifacts = append(outputArtifacts, ens)
	}

	// TODO: Fix engine.Artifact
	return outputArtifacts, nil
}

func (e *ensEngine) handleReorgedLog(
	log *types.Log,
	artifacts []engine.Artifact,
) (
	ENSArtifact,
	error,
) {

	var ensArtifacts []ENSArtifact
	for _, a := range artifacts {
		artifact, ok := a.(ENSArtifact)
		if !ok {
			logger.Debug("found non-pool artifact")
			continue
		}

		ensArtifacts = append(ensArtifacts, artifact)
	}

	artifact := ENSArtifact{BlockNumber: log.BlockNumber}

	logEventKey := log.Topics[0]
	for _, event := range e.ensRegistrar.ContractEvents {
		// This engine is supposed to handle more than 1 event,
		// but it's not yet finished now.
		if logEventKey == event.ID {
			switch event.Name {
			case "NameRegistered":
				logArtifact, err := e.handleNameRegisteredRegistrar(log, event.Name, ensArtifacts)
				if err != nil {
					return artifact, errors.Wrap(err, "failed to create new name from log")
				}
				artifact = *logArtifact
				return artifact, nil
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
