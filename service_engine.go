package superwatcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TODO: Artifact is too generic and this makes it hard
// for services to separate relevant artifacts.

type Artifact any

// ServiceEngine is embedded and injected into WatcherEngine
// to perform business logic.
type ServiceEngine interface {
	// Handle logs from emitter (filter ranged blocks) -- the return type is map of blockHash to []Artifact
	HandleGoodLogs([]*types.Log, []Artifact) (map[common.Hash][]Artifact, error)
	// Handle reorged logs from emitter (filter ranged blocks) -- the return type is map of blockHash to []Artifact
	HandleReorgedLogs([]*types.Log, []Artifact) (map[common.Hash][]Artifact, error)
	// Handle emitter error. If the returned error is not nil, WatcherEngine.HandleEmitterError returns,
	// and the engine shutdowns.
	HandleEmitterError(error) error
}
