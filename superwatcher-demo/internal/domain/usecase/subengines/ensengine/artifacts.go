package ensengine

import (
	"errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	nameRegistered string = "NameRegistered"
)

type ENSArtifact struct {
	ID                  string              `json:"ensID"`
	RegisterBlockNumber uint64              `json:"registerBlockNumber"`
	LastEvent           ENSEvent            `json:"lastEvent"`
	BlockEvents         map[ENSEvent]uint64 `json:"events"` // One block may have >1 events
	TxHash              common.Hash         `json:"txHash"`
	ENS                 entity.ENS          `json:"ens"`
}

func (e *ENSArtifact) ForSubEngine() subengines.SubEngineEnum {
	return subengines.SubEngineENS
}

func getPrevENSArtifactFromLogTxHash(artifacts []ENSArtifact, log *types.Log) *ENSArtifact {
	for _, artifact := range artifacts {
		if artifact.TxHash == log.TxHash {
			return &artifact
		}
	}

	return nil
}

func getPrevArtifactFromLogEntity(artifacts []ENSArtifact, ens entity.ENS) (*ENSArtifact, error) {
	return nil, errors.New("not implemented")
}
