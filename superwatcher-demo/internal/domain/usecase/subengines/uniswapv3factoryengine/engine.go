package uniswapv3factoryengine

import (
	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
)

type uniswapv3PoolFactoryEngine struct {
	poolFactoryContract contracts.BasicContract
}

func New(pooFactoryContract contracts.BasicContract) superwatcher.ServiceEngine {
	return &uniswapv3PoolFactoryEngine{
		poolFactoryContract: pooFactoryContract,
	}
}
