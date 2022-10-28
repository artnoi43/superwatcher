package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

type BasicContract struct {
	ContractABI    abi.ABI
	ContractEvents []abi.Event
}

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

func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}
