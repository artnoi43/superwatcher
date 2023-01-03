package ensengine

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher/pkg/logger"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/lib/logutils"
)

func (e *ensEngine) revertNameRegisteredRegistrar(
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
	if prevArtifact == nil {
		panic("nil prevArtifact")
	}

	name := &prevArtifact.ENS

	// We'll only get ENS Name ID from contract Registrar
	switch logEvent {
	case eventNameRegistered:
		id := log.Topics[1].String()
		if id != name.ID {
			logger.Panic("prevArtifact ENS ID != log ENS ID", zap.String("prevArtifact ID", name.ID), zap.String("log ID", id))
		}

		name.Owner = gslutils.ToLower(log.Topics[2].Hex())
	}

	prevArtifact.LastEvent = RevertRegisterRegistrar
	prevArtifact.BlockEvents[RevertRegisterController] = log.BlockNumber

	return prevArtifact, nil
}

func (e *ensEngine) revertNameRegisteredController(
	log *types.Log,
	logEvent string,
	prevArtifact *ENSArtifact,
) (
	*ENSArtifact,
	error,
) {
	if prevArtifact == nil {
		logger.Panic("handleNameRegisteredController: got nil prevArtifact")
	}

	switch logEvent {
	case "NameRegistered":
		if len(log.Topics) < 3 {
			panic("bad log with < 3 topics - should not happen")
		}

		// Get name from previous artifact
		name := &prevArtifact.ENS
		if name == nil || prevArtifact == nil {
			return nil, errors.New("could not find ENS artifact from Registrar (controller NameRegistered)")
		}

		unpacked, err := logutils.UnpackLogDataIntoMap(e.ensController.ContractABI, logEvent, log.Data)
		if err != nil {
			return nil, err
		}
		domainName, err := logutils.ExtractFieldFromUnpacked[string](unpacked, "name")
		if err != nil {
			return nil, err
		}
		expire, err := logutils.ExtractFieldFromUnpacked[*big.Int](unpacked, "expires")
		if err != nil {
			return nil, err
		}

		name.Name = domainName
		name.Expires = time.Unix(expire.Int64(), 0)

		prevArtifact.RegisterBlockNumber = log.BlockNumber
		prevArtifact.BlockEvents[RegisteredController] = log.BlockNumber
		prevArtifact.LastEvent = RegisteredController
	}

	return prevArtifact, nil
}
