package superwatcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TODO: Artifact is too generic and this makes it hard
// for services to separate relevant artifacts.

type Artifact any

// BaseServiceEngine is shared by ServiceEngine and ServiceThinEngine
type BaseServiceEngine interface {
	HandleEmitterError(error) error
}

// ServiceEngine is embedded and injected into WatcherEngine to perform business logic.
// It is the preferred way to use superwatcher
type ServiceEngine interface {
	BaseServiceEngine

	// Handle logs from emitter (filter ranged blocks) -- the return type is map of blockHash to []Artifact
	HandleGoodLogs([]*types.Log, []Artifact) (map[common.Hash][]Artifact, error)
	// Handle reorged logs from emitter (filter ranged blocks) -- the return type is map of blockHash to []Artifact
	HandleReorgedLogs([]*types.Log, []Artifact) (map[common.Hash][]Artifact, error)
}

// ThinServiceEngine is embedded and injected into thinEngine, a thin implementation of WatcherEngine without managed states.
// It is recommended for niche use cases and advanced users
type ThinServiceEngine interface {
	BaseServiceEngine

	HandleFilterResult(*FilterResult) error
}
