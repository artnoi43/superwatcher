package ensengine

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
)

type ENSEvent string

const (
	Null                     ENSEvent = "NULL"
	Revert                   ENSEvent = "REVERT"
	RegisteredRegistrar      ENSEvent = "RegisteredRegistrar"
	RevertRegisterRegistrar  ENSEvent = Revert + RegisteredRegistrar
	RegisteredController     ENSEvent = "RegisteredController"
	RevertRegisterController ENSEvent = Revert + RegisteredController
	Transferred              ENSEvent = "Transferred"
	RevertTransferred        ENSEvent = Revert + Transferred
)

type ENSArtifact struct {
	ID                  string              `json:"ensID"`
	RegisterBlockNumber uint64              `json:"registerBlockNumber"`
	LastEvent           ENSEvent            `json:"lastEvent"`
	BlockEvents         map[ENSEvent]uint64 `json:"events"` // One block may have >1 events
	ENS                 entity.ENS          `json:"ens"`
}

func (e ENSArtifact) ForSubEngine() subengines.SubEngineEnum {
	return subengines.SubEngineENS
}

func prevRegistrarArtifact(log *types.Log, artifacts []superwatcher.Artifact) *ENSArtifact {
	artifact := filterRegistrarArtifact(log, collectArtifactsENS(artifacts))

	var blankHash common.Hash
	if artifact.LastEvent == Null && common.HexToHash(artifact.ENS.TxHash) == blankHash {
		return nil
	}

	return &artifact
}

func collectArtifactsENS(artifacts []superwatcher.Artifact) []ENSArtifact {
	var ensArtifacts []ENSArtifact //nolint:prealloc
	for _, artifact := range artifacts {
		ensArtifact, ok := artifact.(ENSArtifact)
		if !ok {
			continue
		}

		ensArtifacts = append(ensArtifacts, ensArtifact)
	}

	return ensArtifacts
}

func filterRegistrarArtifact(log *types.Log, artifacts []ENSArtifact) ENSArtifact {
	for _, artifact := range artifacts {
		ensID := log.Topics[1].String()
		if artifact.ENS.ID == gslutils.ToLower(ensID) {
			return artifact
		}
	}

	panic(fmt.Sprintf("no such artifact for txHash %s", log.TxHash.String()))
}

func getPrevArtifactFromLogEntity(artifacts []ENSArtifact, ens entity.ENS) (*ENSArtifact, error) {
	return nil, errors.New("not implemented")
}
