package routerengine

import (
	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
)

type routerArtifact interface {
	ForSubEngine() subengines.SubEngineEnum
}

func artifactIsFor(artifact superwatcher.Artifact, subEngine subengines.SubEngineEnum) bool {
	routerArtifact, ok := artifact.(routerArtifact)
	if !ok {
		return false
	}

	return routerArtifact.ForSubEngine() == subEngine
}

type artifactStore map[subengines.SubEngineEnum][]superwatcher.Artifact
