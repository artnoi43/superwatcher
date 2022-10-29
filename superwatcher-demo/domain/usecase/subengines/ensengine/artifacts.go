package ensengine

import (
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
)

type ENSLastEvent string

const (
	Null                 ENSLastEvent = "NULL"
	RegisteredRegistrar  ENSLastEvent = "RegisteredRegistrar"
	RegisteredController ENSLastEvent = "RegisteredController"
	Transferred          ENSLastEvent = "Transferred"
)

type ENSArtifact struct {
	BlockNumber uint64       `json:"blockNumber"`
	LastEvent   ENSLastEvent `json:"lastEvent"`
	ENS         entity.ENS   `json:"ens"`
}

func (e *ENSArtifact) ForSubEngine() subengines.SubEngineEnum {
	return subengines.SubEngineENS
}
