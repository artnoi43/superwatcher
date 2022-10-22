package uniswapv3factoryengine

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

type uniswapv3FactoryEngine struct {
	mapAddrToABI    map[common.Address]abi.ABI
	mapAddrToEvents map[common.Address][]abi.Event
	serviceFSM      *poolFactoryFSM
}

func NewUniswapV3Engine(
	mapAddrABI map[common.Address]abi.ABI,
	mapAddrEvents map[common.Address][]abi.Event,
) engine.ServiceEngine[entity.Uniswapv3FactoryWatcherKey, *entity.Uniswapv3PoolCreated] {
	return &uniswapv3FactoryEngine{
		mapAddrToABI:    mapAddrABI,
		mapAddrToEvents: mapAddrEvents,
		// TODO: Should we add func `NewPoolFactoryFSM`?
		serviceFSM: &poolFactoryFSM{
			states: make(map[entity.Uniswapv3FactoryWatcherKey]engine.ServiceItemState),
		},
	}
}

func (e *uniswapv3FactoryEngine) ServiceStateTracker() (
	engine.ServiceFSM[entity.Uniswapv3FactoryWatcherKey],
	error,
) {
	if e.serviceFSM == nil {
		return nil, errors.New("nil uniswapv3FactoryEngine.serviceFSM")
	}

	return e.serviceFSM, nil
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
