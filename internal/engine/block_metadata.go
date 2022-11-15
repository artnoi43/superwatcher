package engine

import (
	"fmt"
	"strings"

	"github.com/artnoi43/superwatcher"
)

type blockMetadata struct {
	blockNumber uint64
	blockHash   string // Must be all lowercase
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
	return fmt.Sprintf("%d:%s", k.blockNumber, strings.ToLower(k.blockHash))
}
