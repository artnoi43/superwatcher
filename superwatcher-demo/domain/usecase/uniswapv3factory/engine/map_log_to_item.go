package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	watcherengine "github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/logutils"
)

func (e *uniswapv3FactoryEngine) MapLogToItem(
	log *types.Log,
) (
	watcherengine.ServiceItem[entity.Uniswapv3FactoryWatcherKey],
	error,
) {
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
	for _, event := range contractInterestingEvents {
		if logEventKey == event.ID {
			return mapLogToItem(contractABI, event.Name, log)
		}
	}

	return nil, fmt.Errorf("event topic %s not found", logEventKey)
}

// mapLogToItem maps *types.Log to entity.Uniswapv3PoolCreated.
// It can be updated to handle more events from this contract.
func mapLogToItem(
	contractABI abi.ABI,
	eventName string,
	l *types.Log,
) (
	watcherengine.ServiceItem[entity.Uniswapv3FactoryWatcherKey],
	error,
) {
	unpacked, err := logutils.UnpackLogDataIntoMap(contractABI, eventName, l.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack uniswapv3factory logs (event %s)", eventName)
	}
	poolCreated, err := parseLogDataToUniswapv3Factory(unpacked)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert unpacked log data to *entity.Uniswapv3PoolCreated")
	}

	poolCreated.Token0 = common.HexToAddress(l.Topics[1].String())
	poolCreated.Token1 = common.HexToAddress(l.Topics[2].String())
	poolCreated.Fee = l.Topics[3].Big().Uint64()
	poolCreated.BlockCreated = l.BlockNumber

	return poolCreated, nil
}
