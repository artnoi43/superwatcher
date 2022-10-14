package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	watcherengine "github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/logutils"
)

func (e *uniswapv3FactoryEngine) MapLogToItem(log *types.Log) (watcherengine.ServiceItem[string], error) {
	contractAddr := log.Address

	contractABI, ok := e.mapAddrToABI[contractAddr]
	if !ok {
		return nil, fmt.Errorf("abi not found for address %s", contractAddr.String())
	}

	contractInterestingEvents, ok := e.mapAddrToEvents[contractAddr]
	if !ok {
		return nil, fmt.Errorf("events not found for address %s", contractAddr.String())
	}

	logEventKey := log.Topics[0]
	unpacked := make(map[string]interface{})
	var err error
	for _, event := range contractInterestingEvents {
		if logEventKey == event.ID {
			unpacked, err = logutils.UnpackIntoMap(contractABI, event.Name, log)
			if err != nil {
				return nil, errors.New("failed to unpack uniswapv3factory logs")
			}
		}
	}

	poolCreated, err := parseLogToUniswapv3Factory(unpacked)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert unpacked log data to *entity.UniswapSwap")
	}

	return poolCreated, nil
}
