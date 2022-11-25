package servicetest

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

type engine struct {
	reorgedAt          uint64
	emitterFilterRange uint64
	debugger           debugger.Debugger
}

func (e *engine) HandleGoodLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	return nil, nil
}

func (e *engine) HandleReorgedLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	e.debugger.Debug(1, "GOT REORGED LOG IN SERVICETEST")
	for _, log := range logs {
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
	return err
}
