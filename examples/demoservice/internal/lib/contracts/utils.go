package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

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
