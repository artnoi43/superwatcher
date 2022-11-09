package superwatcher

import "github.com/ethereum/go-ethereum/core/types"

// TODO: Artifact is too generic and this makes it hard
// for services to separate relevant artifacts.

type Artifact any

// ServiceEngine is embedded and injected into WatcherEngine
// to perform business logic.
type ServiceEngine interface {
	// Handle a block's logs || Maybe changed to just logs, not a block's logs
	HandleGoodLogs([]*types.Log, []Artifact) ([]Artifact, error)
	// Handle a block's reorged logs || Maybe changed to just logs, not a block's logs
	HandleReorgedLogs([]*types.Log, []Artifact) ([]Artifact, error)
	// Handle emitter error
	HandleEmitterError(error) error
}
