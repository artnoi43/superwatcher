package engine

import (
	"fmt"

	"github.com/artnoi43/superwatcher/pkg/superwatcher"
)

type blockMetadata struct {
	blockNumber uint64
	state       EngineBlockState

	// artifacts maybe removed - I see no use case yet
	artifacts []superwatcher.Artifact
}

func (k blockMetadata) BlockNumber() uint64 {
	// TODO: Here for debugging
	if k.blockNumber == 0 {
		panic("got blockNumber 0 from a serviceLogStateKey")
	}

	return k.blockNumber
}

func (k blockMetadata) String() string {
	return fmt.Sprintf("%d", k.blockNumber)
}
