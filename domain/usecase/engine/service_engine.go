package engine

import "github.com/ethereum/go-ethereum/core/types"

type Artifact any
type Artifacts []Artifact

type ServiceEngine interface {
	// Handle a block's logs
	HandleGoodBlockLogs([]*types.Log, Artifacts) (Artifacts, error)
	// Handle a block's reorged logs
	HandleReorgedBlockLogs([]*types.Log, Artifacts) (Artifacts, error)
	// Handle emitter error
	HandleEmitterError(error) error
}
