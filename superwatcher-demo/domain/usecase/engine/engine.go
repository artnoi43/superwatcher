package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/logutils"
)

type uniswapv3FactoryEngine struct {
	mapAddrToABI    map[common.Address]abi.ABI
	mapAddrToEvents map[common.Address][]abi.Event
}

func NewUniswapV3Engine(
	mapAddrABI map[common.Address]abi.ABI,
	mapAddrEvents map[common.Address][]abi.Event,
) *uniswapv3FactoryEngine {
	return &uniswapv3FactoryEngine{
		mapAddrToABI:    mapAddrABI,
		mapAddrToEvents: mapAddrEvents,
	}
}

func (e *uniswapv3FactoryEngine) MapLogToItem(log *types.Log) (*entity.Uniswapv3PoolCreated, error) {
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

func parseLogToUniswapv3Factory(unpacked map[string]interface{}) (*entity.Uniswapv3PoolCreated, error) {
	return nil, errors.New("not implemented")
}
