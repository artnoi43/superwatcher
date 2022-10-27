package uniswapv3factoryengine

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

type uniswapv3PoolFactoryEngine struct {
	contractABI    abi.ABI
	contractEvents []abi.Event
}

func NewUniswapV3Engine(
	contractABI abi.ABI,
	contractEvents []abi.Event,
) engine.ServiceEngine {
	return &uniswapv3PoolFactoryEngine{
		contractABI:    contractABI,
		contractEvents: contractEvents,
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
