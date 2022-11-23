package servicetest

import (
	"fmt"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
	"github.com/ethereum/go-ethereum/core/types"
)

type engine struct {
	reorgedAt       uint64
	emitterLookBack uint64
	debugger        debugger.Debugger
}

func (e *engine) HandleGoodLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	return nil, nil
}

func (e *engine) HandleReorgedLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	e.debugger.Debug("GOT REORGED LOG IN SERVICETEST")
	for _, log := range logs {
		// TODO: Polish test checks
		if log.BlockNumber != e.reorgedAt {
			if log.BlockNumber > e.reorgedAt+e.emitterLookBack {
				return nil, fmt.Errorf("reorgedAt is different from logs passed to HandleReorgedLogs: expecting %d, got %d", e.reorgedAt, log.BlockNumber)
			}
		}
	}
	return nil, nil
}

func (e *engine) HandleEmitterError(err error) error {
	return err
}
