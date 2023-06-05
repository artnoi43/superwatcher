package routerengine

import (
	"github.com/soyart/gsl"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/soyart/superwatcher"
	"github.com/soyart/superwatcher/pkg/logger"

	"github.com/soyart/superwatcher/examples/demoservice/internal/subengines"
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

func (e *routerEngine) mapBlocksToSubEngine(blocks []*superwatcher.Block) map[subengines.SubEngineEnum][]*superwatcher.Block {
	blocksMap := make(map[subengines.SubEngineEnum][]*superwatcher.Block)

	for _, block := range blocks {
		logsMap := e.mapLogsToSubEngine(block.Logs)

		for subEngine, logs := range logsMap {
			blocksMap[subEngine] = append(blocksMap[subEngine], &superwatcher.Block{
				Number: block.Number,
				Hash:   block.Hash,
				Header: block.Header,
				Logs:   logs,
			})
		}
	}

	return blocksMap
}

func (e *routerEngine) logToSubEngine(log *types.Log) (subengines.SubEngineEnum, bool) {
	for subEngine, addrTopics := range e.Routes {
		for address, topics := range addrTopics {
			if address == log.Address {
				if gsl.Contains(topics, log.Topics[0]) {
					return subEngine, true
				}
			}
		}
	}

	return subengines.SubEngineInvalid, false
}

func (e *routerEngine) logToService(log *types.Log) superwatcher.ServiceEngine { //nolint:unused
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
