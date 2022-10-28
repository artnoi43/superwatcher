package uniswapv3factoryengine

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts"
)

type uniswapv3PoolFactoryEngine struct {
	poolFactoryContract contracts.BasicContract
}

func NewUniswapV3Engine(
	contractABI abi.ABI,
	contractEvents []abi.Event,
) engine.ServiceEngine {
	return &uniswapv3PoolFactoryEngine{
		poolFactoryContract: contracts.BasicContract{
			ContractABI:    contractABI,
			ContractEvents: contractEvents,
		},
	}
}
