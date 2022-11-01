package demoengine

import (
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines"
)

func (e *demoEngine) mapLogsToSubEngine(logs []*types.Log) map[subengines.SubEngineEnum][]*types.Log {
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

func (e *demoEngine) logToSubEngine(log *types.Log) (subengines.SubEngineEnum, bool) {
	se, ok := e.usecases[log.Address]
	if !ok {
		return subengines.SubEngineInvalid, false
	}

	return se, true
}

func (e *demoEngine) logToService(log *types.Log) superwatcher.ServiceEngine {
	subEngine, ok := e.logToSubEngine(log)
	if !ok {
		logger.Panic("log address not mapped to subengine - should not happen", zap.String("address", log.Address.String()))
	}

	serviceEngine, ok := e.services[subEngine]
	if !ok {
		logger.Panic(
			"usecase has no service",
			zap.String("subengine usecase", subEngine.String()),
		)
	}

	return serviceEngine
}
