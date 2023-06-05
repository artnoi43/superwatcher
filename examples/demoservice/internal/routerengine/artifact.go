package routerengine

import (
	"github.com/soyart/superwatcher"

	"github.com/soyart/superwatcher/examples/demoservice/internal/subengines"
)

type routerArtifact interface {
	ForSubEngine() subengines.SubEngineEnum
}

func artifactIsFor(
	subEngine subengines.SubEngineEnum,
	artifact superwatcher.Artifact,
) bool {
	routerArtifact, ok := artifact.(routerArtifact)
	if !ok {
		panic("not routerArtifact")
	}

	return routerArtifact.ForSubEngine() == subEngine
}

func filterArtifacts(
	subEngine subengines.SubEngineEnum,
	artifacts []superwatcher.Artifact,
) []superwatcher.Artifact {
	// Each sub-engine already returns `[]superwatcher.Artifact`,
	// and method routerEngine.HandleLogs treats each returned `[]superwatcher.Artifact` as `superwatcher.Artifact`,
	// i.e. if we have 3 sub-engines A, B, and C -- then the artifacts returned by routerEngine.HandleLogs
	// will have 3 elements: [ []superwatcher.Artifact from A, []superwatcher.Artifact from B, []superwatcher.Artifact from C ]

	var subEngineArtifacts []superwatcher.Artifact
	for _, seArtifacts := range artifacts {
		for _, seArtifact := range seArtifacts.([]superwatcher.Artifact) {
			if artifactIsFor(subEngine, seArtifact) {
				subEngineArtifacts = append(subEngineArtifacts, seArtifact)
			}
		}
	}

	if len(subEngineArtifacts) > 0 {
		return subEngineArtifacts
	}

	return nil
}
