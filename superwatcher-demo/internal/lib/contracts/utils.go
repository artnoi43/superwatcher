package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// ContractInfo reads abiStr and returns the Go ABI,
// as well as all `abi.Event`s whose name matches eventKeys.
func ContractInfo(contractABI abi.ABI, eventKeys ...string) (abi.ABI, []abi.Event, error) {
	var events []abi.Event
	for _, eventKey := range eventKeys {
		event, found := contractABI.Events[eventKey]
		if !found {
			return abi.ABI{}, nil, errors.Wrapf(ErrNoSuchABIEvent, "eventKey %s not found", eventKey)
		}
		events = append(events, event)
	}

	return contractABI, events, nil
}

func accrueEvents(contractABI abi.ABI, eventKeys ...string) []abi.Event {
	var events []abi.Event
	for key, event := range contractABI.Events {
		for _, eventKey := range eventKeys {
			if key == eventKey {
				events = append(events, event)
			}
		}
	}

	return events
}

func CollectEventHashes(events []abi.Event) []common.Hash {
	hashes := make([]common.Hash, len(events))
	for i, event := range events {
		hashes[i] = event.ID
	}

	return hashes
}
