package servicetest

import (
	"fmt"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debug"
	"github.com/ethereum/go-ethereum/core/types"
)

type engine struct {
	reorgedAt       uint64
	emitterLookBack uint64
}

func (e *engine) HandleGoodLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	return nil, nil
}

func (e *engine) HandleReorgedLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	debug.DebugMsg(true, "GOT REORG LOGS IN SERVICE ENGINE")

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
