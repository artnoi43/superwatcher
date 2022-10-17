package engine

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
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

func parseLogDataToUniswapv3Factory(unpacked map[string]interface{}) (*entity.Uniswapv3PoolCreated, error) {
	var poolCreated entity.Uniswapv3PoolCreated
	for k, v := range unpacked {
		switch k {
		case "pool":
			poolAddr, ok := v.(common.Address)
			if !ok {
				return nil, errors.New("type assertion to common.Address failed")
			}
			poolCreated.Address = poolAddr
		}
	}
	return &poolCreated, nil
}
