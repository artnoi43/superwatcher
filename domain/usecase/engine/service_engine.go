package engine

import "github.com/ethereum/go-ethereum/core/types"

// TODO: Artifact is too generic and this makes it hard
// for services to separate relevant artifacts.

type Artifact any

type ServiceEngine interface {
	// Handle a block's logs
	HandleGoodLogs([]*types.Log) ([]Artifact, error)
	// Handle a block's reorged logs
	HandleReorgedLogs([]*types.Log, []Artifact) ([]Artifact, error)
	// Handle emitter error
	HandleEmitterError(error) error
}
