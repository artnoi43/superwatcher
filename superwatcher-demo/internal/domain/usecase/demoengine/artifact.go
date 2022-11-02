package demoengine

import (
	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines"
)

type subEngineArtifact interface {
	ForSubEngine() subengines.SubEngineEnum
}

func artifactIsFor(artifact superwatcher.Artifact, subEngine subengines.SubEngineEnum) bool {
	demoArtifact, ok := artifact.(subEngineArtifact)
	if !ok {
		return false
	}

	return demoArtifact.ForSubEngine() == subEngine
}

// // TODO: maybe removed - not used yet
// // Exported for debugging
// type DemoArtifact struct {
// 	// PoolFactoryArtifact is a hash map - so it can have 1 artifact for multiple logs
// 	PoolFactoryArtifact subEngineArtifact
// 	// EnsArtifact is a struct, so it needs an array to represent multiple logs
// 	EnsArtifact []subEngineArtifact
// }
//
// func (a *DemoArtifact) AddSubEngineArtifact(
// 	subEngineArtifacts []superwatcher.Artifact,
// ) error {
// 	for _, artifact := range subEngineArtifacts {
// 		seArtifact, ok := artifact.(subEngineArtifact)
// 		if !ok {
// 			debug.DebugMsg(true, "artifact is not subEngineArtifact", zap.String("actual type", reflect.TypeOf(artifact).String()))
// 			continue
// 		}
//
// 		switch seArtifact.ForSubEngine() {
// 		case subengines.SubEngineENS:
// 			a.EnsArtifact = append(a.EnsArtifact, seArtifact)
// 		case subengines.SubEngineUniswapv3Factory:
// 			a.PoolFactoryArtifact = seArtifact
// 		}
// 	}
//
// 	return errors.New("unknown subengineArtifact")
// }
