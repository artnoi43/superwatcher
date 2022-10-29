package ensengine

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/logutils"
)

// Multiple event logs from ENS have owner or new owner address in log.Data, not in topics.
func extractOwnerAddressFromUnpacked(unpacked map[string]interface{}) (common.Address, error) {
	return logutils.ExtractFieldFromUnpacked[common.Address](unpacked, "owner")
}

func extractTTLFromUnpacked(unpacked map[string]interface{}) (uint64, error) {
	return logutils.ExtractFieldFromUnpacked[uint64](unpacked, "ttl")
}

// TODO: this is broken. To create a new ENS name, we need logs from 2 contracts,
// but 1 log can only come from 1 contract
func (e *ensEngine) handleNameRegisteredRegistrar(
	log *types.Log,
	logEvent string,
	artifacts []ENSArtifact,
) (
	*ENSArtifact,
	error,
) {
	if len(log.Topics) < 2 {
		panic("bad log with < 2 topics")
	}

	var name entity.ENS
	// We'll only get ENS Name ID from contract Registrar
	switch logEvent {
	case "NameRegistered":
		name.ID = common.HexToHash(log.Topics[1].Hex()).String()
		name.Owner = common.HexToAddress(log.Topics[2].Hex())
	}

	return &ENSArtifact{
		BlockNumber: log.BlockNumber,
		LastEvent:   Registered,
		ENS:         name,
	}, nil
}

func (e *ensEngine) handleNameRegisteredController(
	log *types.Log,
	logEvent string,
	artifacts []ENSArtifact,
) (
	*ENSArtifact,
	error,
) {
	if len(log.Topics) < 3 {
		panic("bad log with < 3 topics")
	}

	var name entity.ENS
	// We'll only get ENS Name ID from contract Registrar
	switch logEvent {
	case "NameRegistered":
		owner := common.HexToAddress(log.Topics[2].Hex())

		// Find targetENS from
		var targetENS entity.ENS
		for _, artifact := range artifacts {
			ens := artifact.ENS
			if ens.Owner == owner {
				targetENS = ens
			}
		}

		unpacked, err := logutils.UnpackLogDataIntoMap(e.ensController.ContractABI, logEvent, log.Data)
		if err != nil {
			return nil, err
		}
		name, err := logutils.ExtractFieldFromUnpacked[string](unpacked, "name")
		if err != nil {
			return nil, err
		}
		expire, err := logutils.ExtractFieldFromUnpacked[*big.Int](unpacked, "expire")
		if err != nil {
			return nil, err
		}

		targetENS.Name = name
		targetENS.Expire = time.Unix(expire.Int64(), 0)
	}

	return &ENSArtifact{
		BlockNumber: log.BlockNumber,
		LastEvent:   Registered,
		ENS:         name,
	}, nil
}
