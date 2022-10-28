package ensengine

import (
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
)

type ENSLastEvent string

const (
	Null        ENSLastEvent = ""
	Registered  ENSLastEvent = "Registered"
	Transferred ENSLastEvent = "Transferred"
)

type ENSArtifact struct {
	BlockNumber uint64     `json:"blockNumber"`
	LastEvent   string     `json:"lastEvent"`
	ENS         entity.ENS `json:"ens"`
}

func (e *ENSArtifact) ForSubEngine() subengines.SubEngineEnum {
	return subengines.SubEngineENS
}
