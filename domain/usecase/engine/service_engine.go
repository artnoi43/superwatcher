package engine

import "github.com/ethereum/go-ethereum/core/types"

type Artifact any

type ServiceEngine interface {
	// Handle a block's logs
	HandleGoodBlock([]*types.Log, []Artifact) ([]Artifact, error)
	// Handle a block's reorged logs
	HandleReorgedBlock([]*types.Log, []Artifact) ([]Artifact, error)
	// Handle emitter error
	HandleEmitterError(error) error
}
