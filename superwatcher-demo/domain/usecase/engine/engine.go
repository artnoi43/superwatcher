package engine

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

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

func parseLogToUniswapv3Factory(unpacked map[string]interface{}) (*entity.Uniswapv3PoolCreated, error) {
	return nil, errors.New("not implemented")
}
