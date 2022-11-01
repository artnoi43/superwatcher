package uniswapv3factoryengine

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
)

type uniswapv3PoolFactoryEngine struct {
	poolFactoryContract contracts.BasicContract
}

func NewUniswapV3Engine(
	contractAddress common.Address,
	contractABI abi.ABI,
	contractEvents []abi.Event,
) superwatcher.ServiceEngine {
	return &uniswapv3PoolFactoryEngine{
		poolFactoryContract: contracts.BasicContract{
			Address:        contractAddress,
			ContractABI:    contractABI,
			ContractEvents: contractEvents,
		},
	}
}
