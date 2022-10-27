package demoengine

import "github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"

type demoArtifact interface {
	ForSubEngine() subengines.SubEngine
}
