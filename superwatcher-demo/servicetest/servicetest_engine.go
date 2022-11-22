package servicetest

import (
	"github.com/artnoi43/superwatcher"
	"github.com/ethereum/go-ethereum/core/types"
)

type engine struct{}

func (e *engine) HandleGoodLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	return nil, nil
}

func (e *engine) HandleReorgedLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	panic("got reorged logs")
}

func (e *engine) HandleEmitterError(err error) error {
	return err
}
