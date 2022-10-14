package contracts

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

// ContractInfo reads abiStr and returns the Go ABI,
// as well as all `abi.Event`s whose name matches eventKeys.
func ContractInfo(abiStr string, eventKeys ...string) (abi.ABI, []abi.Event, error) {
	contractABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return abi.ABI{}, nil, errors.Wrap(err, "read ABI failed")
	}

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
