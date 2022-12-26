package ensengine

import (
	"context"
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal"
)

// The sub-engine uses entity.ENS as superwatcher.Artifact

// MapLogToItem wraps mapLogToItem, so the latter can be unit tested.
func (e *ensEngine) HandleGoodLogs(
	logs []*types.Log,
	artifacts []superwatcher.Artifact,
) (
	map[common.Hash][]superwatcher.Artifact,
	error,
) {
	e.debugger.Debug(3, "HandleGoodLogs called")

	outArtifacts := make(map[common.Hash][]superwatcher.Artifact)
	for _, log := range logs {
		logArtifact, err := e.HandleGoodLog(log, artifacts)
		if err != nil {
			if errors.Is(err, internal.ErrNoNeedHandle) {
				continue
			}
			return nil, errors.Wrapf(err, "HandleGoodLog failed on log txHash %s", log.BlockHash.String())
		}

		artifacts = append(artifacts, logArtifact)
		outArtifacts[log.BlockHash] = append(outArtifacts[log.BlockHash], logArtifact)
	}

	return outArtifacts, nil
}

func (e *ensEngine) HandleGoodLog(
	log *types.Log,
	artifacts []superwatcher.Artifact,
) (
	ENSArtifact,
	error,
) {
	// Artifact to return
	artifact := ENSArtifact{RegisterBlockNumber: log.BlockNumber}

	var handleFunc func(*types.Log, string, *ENSArtifact) (*ENSArtifact, error)
	var eventName string
	var prevArtifact *ENSArtifact

	logEventKey := log.Topics[0]
	switch log.Address {
	case e.ensRegistrar.Address:
		for _, event := range e.ensRegistrar.ContractEvents {
			// This engine is supposed to handle more than 1 event,
			// but it's not yet finished now.
			if logEventKey == event.ID {
				switch event.Name {
				// New domain registered from both contracts
				case eventNameRegistered:
					handleFunc = e.handleNameRegisteredRegistrar
					eventName = eventNameRegistered
				default:
					continue
				}
			}
		}
	case e.ensController.Address:
		for _, event := range e.ensController.ContractEvents {
			if logEventKey == event.ID {
				switch event.Name {
				case eventNameRegistered:
					handleFunc = e.handleNameRegisteredController
					eventName = eventNameRegistered
					// Get previous artifacts
					prevArtifact = prevRegistrarArtifact(log, artifacts)
					if prevArtifact == nil {
						panic("nil prevArtifact")
					}
				default:
					continue
				}
			}
		}
	default:
		panic("ensEngine.handleGoodLog: found unexpected contract address: " + log.Address.String())
	}

	if handleFunc == nil {
		e.debugger.Debug(3, "ensEngine: handleFunc is nil, probably because uninteresting topics", zap.Any("artifact", artifacts))
		return artifact, internal.ErrNoNeedHandle
	}

	resultArtifact, err := handleFunc(log, eventName, prevArtifact)
	if err != nil {
		return artifact, errors.Wrapf(err, "failed to create new name from log (event %s)", eventName)
	}

	if resultArtifact == nil {
		panic("nil resultArtifact")
	}
	artifact = *resultArtifact

	err = e.dataGateway.SetENS(context.Background(), &artifact.ENS)
	if err != nil {
		err = errors.Wrapf(err, "Save ENS failed on log txHash %s", log.BlockHash.String())
	}

	return artifact, err
}

func (e *ensEngine) HandleReorgedLogs(
	logs []*types.Log,
	artifacts []superwatcher.Artifact,
) (
	map[common.Hash][]superwatcher.Artifact,
	error,
) {
	e.debugger.Debug(1, fmt.Sprintf("got %d reorged logs and %d artifacts", len(logs), len(artifacts)), zap.Any("artifacts", artifacts))

	outputArtifacts := make(map[common.Hash][]superwatcher.Artifact)
	for _, log := range logs {
		ens, err := e.handleReorgedLog(log, artifacts)
		if err != nil {
			return nil, errors.Wrap(err, "ensEngine.handleReorgedLog failed")
		}

		outputArtifacts[log.BlockHash] = append(outputArtifacts[log.BlockHash], ens)

		err = e.dataGateway.DelENS(context.Background(), &ens.ENS)
		if err != nil {
			return nil, errors.Wrapf(err, "Remove ENS failed on log txHash %s", gslutils.StringerToLowerString(log.BlockHash))
		}
	}

	return outputArtifacts, nil
}

// handleReorgedLog examines the log, get log's previous artifact, and handle chain reorg events
func (e *ensEngine) handleReorgedLog(
	log *types.Log,
	artifacts []superwatcher.Artifact,
) (
	ENSArtifact,
	error,
) {
	// Previous artifacts
	prevArtifact := prevRegistrarArtifact(log, artifacts)
	if prevArtifact == nil {
		panic("nil prevArtifact")
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
				case eventNameRegistered:
					reorgArtifact, err := e.revertNameRegisteredRegistrar(log, event.Name, prevArtifact)
					if err != nil {
						return artifact, errors.Wrap(err, "failed to create new name from log")
					}
					artifact = *reorgArtifact
					return artifact, nil

				default:
					continue
				}
			}
		}
	case e.ensController.Address:
		for _, event := range e.ensController.ContractEvents {
			// This engine is supposed to handle more than 1 event,
			// but it's not yet finished now.
			if logEventKey == event.ID {
				switch event.Name {
				case eventNameRegistered:
					reorgArtifact, err := e.revertNameRegisteredController(log, event.Name, prevArtifact)
					if err != nil {
						return artifact, errors.Wrap(err, "failed to create new name from log")
					}
					artifact = *reorgArtifact
					return artifact, nil

				default:
					continue
				}
			}
		}
	}

	return artifact, fmt.Errorf("event topic %s not found", logEventKey)
}

// Unused by this service
func (e *ensEngine) HandleEmitterError(err error) error {
	logger.Warn("emitter error", zap.Error(err))
	return err
}
