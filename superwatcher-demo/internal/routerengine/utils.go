package routerengine

import (
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
)

func (e *routerEngine) mapLogsToSubEngine(logs []*types.Log) map[subengines.SubEngineEnum][]*types.Log {
	logsMap := make(map[subengines.SubEngineEnum][]*types.Log)

	for _, log := range logs {
		subEngine, ok := e.logToSubEngine(log)
		if !ok {
			continue
		}
		logsMap[subEngine] = append(logsMap[subEngine], log)
	}

	return logsMap
}

func (e *routerEngine) logToSubEngine(log *types.Log) (subengines.SubEngineEnum, bool) {
	for subEngine, addrTopics := range e.Routes {
		for address, topics := range addrTopics {
			if address == log.Address {
				if gslutils.Contains(topics, log.Topics[0]) {
					return subEngine, true
				}
			}
		}

	}

	return subengines.SubEngineInvalid, false
}

func (e *routerEngine) logToService(log *types.Log) superwatcher.ServiceEngine {
	subEngine, ok := e.logToSubEngine(log)
	if !ok {
		logger.Panic("log address not mapped to subengine - should not happen", zap.String("address", log.Address.String()))
	}

	serviceEngine, ok := e.Services[subEngine]
	if !ok {
		logger.Panic(
			"usecase has no service",
			zap.String("subengine usecase", subEngine.String()),
		)
	}

	return serviceEngine
}

func mergeArtifacts(
	routerArtifacts map[common.Hash][]superwatcher.Artifact, // Usually empty
	subEngineArtifacts map[common.Hash][]superwatcher.Artifact,
) {
	for blockHash := range subEngineArtifacts {
		routerArtifacts[blockHash] = append(routerArtifacts[blockHash], subEngineArtifacts[blockHash]...)
	}
}
