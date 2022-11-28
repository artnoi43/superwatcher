package servicetest

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

type engine struct {
	reorgedAt          uint64
	emitterFilterRange uint64
	debugger           *debugger.Debugger
}

func (e *engine) HandleGoodLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	e.debugger.Debug(2, fmt.Sprintf("HandleGoodLogs: got %d logs", len(logs)))
	for _, log := range logs {
		e.debugger.Debug(
			1, "good log info",
			zap.Uint64("blockNumber", log.BlockNumber),
			zap.String("blockHash", gslutils.StringerToLowerString(log.BlockHash)),
		)
	}
	return nil, nil
}

func (e *engine) HandleReorgedLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
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

		// TODO: Polish test checks
		if log.BlockNumber != e.reorgedAt {
			if log.BlockNumber > e.reorgedAt+e.emitterFilterRange {
				return nil, fmt.Errorf(
					"reorgedAt is different from logs passed to HandleReorgedLogs: expecting %d, got %d",
					e.reorgedAt, log.BlockNumber,
				)
			}
		}
	}

	return nil, nil
}

func (e *engine) HandleEmitterError(err error) error {
	e.debugger.Debug(1, "got error", zap.Error(err))
	return err
}
