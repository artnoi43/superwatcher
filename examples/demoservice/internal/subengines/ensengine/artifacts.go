package ensengine

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines"
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
	for i, artifact := range artifacts {
		if artifact == nil {
			panic(fmt.Sprintf("nil artifact at index %d", i))
		}

		ensArtifact, ok := artifact.(ENSArtifact)
		if ok {
			ensArtifacts = append(ensArtifacts, ensArtifact)
			continue
		}

		for _, seArtifact := range artifact.([]superwatcher.Artifact) {
			ensArtifact, ok := seArtifact.(ENSArtifact)
			if !ok {
				continue
			}

			ensArtifacts = append(ensArtifacts, ensArtifact)
		}
	}

	return ensArtifacts
}

func filterRegistrarArtifact(log *types.Log, artifacts []ENSArtifact) ENSArtifact {
	if artifacts == nil {
		panic("nil artifacts")
	}
	if len(artifacts) == 0 {
		panic(fmt.Sprintf("0 artifacts from blockHash %s", log.BlockHash.String()))
	}

	for _, artifact := range artifacts {
		ensID := gslutils.StringerToLowerString(log.Topics[1])
		if artifact.ENS.ID == ensID {
			return artifact
		}

		txHash := gslutils.StringerToLowerString(log.TxHash)
		if artifact.ENS.TxHash == txHash {
			return artifact
		}
	}

	panic(fmt.Sprintf("no such artifact for txHash %s", log.TxHash.String()))
}
