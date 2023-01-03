package servicetest

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

type DebugServiceEngine struct {
	ReorgedAt          uint64
	EmitterFilterRange uint64

	// Users can inject these extra debugging functions
	HandleFuncGoodLog      func(*types.Log) error
	HandleFuncReorgedLog   func(*types.Log) error
	HandleFuncEmitterError func(error) error

	debugger *debugger.Debugger
}

// Implements superwatcher.ThinServiceEngine
func (e *DebugServiceEngine) HandleResult(result *superwatcher.PollResult) error {
	e.debugger.Debug(2, fmt.Sprintf("Got result: %d GoodBlocks, %d ReorgedBlocks", len(result.GoodBlocks), len(result.ReorgedBlocks)))

	for _, block := range result.ReorgedBlocks {
		for _, log := range block.Logs {
			if e.HandleFuncReorgedLog != nil {
				if err := e.HandleFuncReorgedLog(log); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (e *DebugServiceEngine) HandleGoodLogs(
	logs []*types.Log,
	artifacts []superwatcher.Artifact,
) (
	[]superwatcher.Artifact,
	error,
) {
	e.debugger.Debug(2, fmt.Sprintf("HandleGoodLogs: got %d logs", len(logs)))
	for _, log := range logs {
		e.debugger.Debug(
			1, "good log info",
			zap.Uint64("blockNumber", log.BlockNumber),
			zap.String("blockHash", gslutils.StringerToLowerString(log.BlockHash)),
		)

		// Calling injected func
		if e.HandleFuncGoodLog != nil {
			if err := e.HandleFuncGoodLog(log); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func (e *DebugServiceEngine) HandleReorgedLogs(
	logs []*types.Log,
	artifacts []superwatcher.Artifact,
) (
	[]superwatcher.Artifact,
	error,
) {
	e.debugger.Debug(
		1, "GOT REORGED LOG IN SERVICETEST",
		zap.Int("len(artifacts)", len(artifacts)),
		zap.Any("artifacts", artifacts),
	)

	for _, log := range logs {
		e.debugger.Debug(
			1, "reorged log info",
			zap.Uint64("blockNumber", log.BlockNumber),
			zap.String("blockHash", gslutils.StringerToLowerString(log.BlockHash)),
		)

		if e.HandleFuncReorgedLog != nil {
			err := e.HandleFuncReorgedLog(log)
			if err != nil {
				return nil, err
			}
		}

		// TODO: Polish test checks
		if log.BlockNumber != e.ReorgedAt {
			if log.BlockNumber > e.ReorgedAt+e.EmitterFilterRange {
				return nil, fmt.Errorf(
					"reorgedAt is different from logs passed to HandleReorgedLogs: expecting %d, got %d",
					e.ReorgedAt, log.BlockNumber,
				)
			}
		}
	}

	return nil, nil
}

func (e *DebugServiceEngine) HandleEmitterError(err error) error {
	e.debugger.Debug(1, "got error", zap.Error(err))

	if e.HandleFuncEmitterError != nil {
		err = e.HandleFuncEmitterError(err)
	}

	return err
}
