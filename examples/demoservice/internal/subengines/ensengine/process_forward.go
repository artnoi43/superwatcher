package ensengine

import (
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
)

func (e *ensEngine) handleNameRegisteredRegistrar(
	log *types.Log,
	logEvent string,
	prevArtifact *ENSArtifact,
) (
	*ENSArtifact,
	error,
) {
	if len(log.Topics) < 2 {
		panic("bad log with < 2 topics")
	}

	name := entity.ENS{
		TxHash:      gslutils.StringerToLowerString(log.TxHash),
		BlockHash:   gslutils.StringerToLowerString(log.BlockHash),
		BlockNumber: log.BlockNumber,
	}

	// We'll only get ENS Name ID from contract Registrar
	switch logEvent {
	case eventNameRegistered:
		if err := e.unmarshalLogToENS(logEvent, log, &name); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal ENS registrar log to entity.ENS")
		}
	}

	return &ENSArtifact{
		ID:                  name.ID,
		RegisterBlockNumber: log.BlockNumber,
		BlockEvents: map[ENSEvent]uint64{
			RegisteredRegistrar: log.BlockNumber,
		},
		LastEvent: RegisteredRegistrar,
		ENS:       name,
	}, nil
}

// handleNameRegisteredController is called after handleNameRegisteredRegistrar
// during event |NameRegistered|, because the Registrar TX comes before Controller during registration,
// therefore we'll meed some artifact to merge ENS data from 2 logs
func (e *ensEngine) handleNameRegisteredController(
	log *types.Log,
	logEvent string,
	prevArtifact *ENSArtifact,
) (
	*ENSArtifact,
	error,
) {
	if prevArtifact == nil {
		panic("handleNameRegisteredController: got nil prevArtifact")
	}

	if prevArtifact.BlockEvents == nil {
		prevArtifact.BlockEvents = make(map[ENSEvent]uint64)
	}
	if prev, curr := common.HexToHash(prevArtifact.ENS.BlockHash), log.BlockHash; prev != curr {
		panic(fmt.Sprintf("controller prevArtifact from registrar has different blockHash: %s vs %s", prev.String(), curr.String()))
	}
	if prev, curr := common.HexToHash(prevArtifact.ENS.TxHash), log.TxHash; prev != curr {
		panic(fmt.Sprintf("controller prevArtifact from registrar has different txHash: %s vs %s", prev.String(), curr.String()))
	}

	switch logEvent {
	case eventNameRegistered:
		if len(log.Topics) < 3 {
			panic("bad log with < 3 topics - should not happen")
		}

		// Get name from previous artifact
		name := &prevArtifact.ENS
		if name == nil || prevArtifact == nil {
			return nil, errors.New("could not find ENS artifact from Registrar (controller NameRegistered)")
		}

		if err := e.unmarshalLogToENS(logEvent, log, name); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal ENS Controller log to entity.ENS")
		}

		prevArtifact.RegisterBlockNumber = log.BlockNumber
		prevArtifact.BlockEvents[RegisteredController] = log.BlockNumber
		prevArtifact.LastEvent = RegisteredController
	}

	return prevArtifact, nil
}
