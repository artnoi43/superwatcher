package uniswapv3factoryengine

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
)

type uniswapv3PoolFactoryEngine struct {
	poolFactoryContract contracts.BasicContract
	debugger            *debugger.Debugger
}

func New(pooFactoryContract contracts.BasicContract, logLevel uint8) superwatcher.ServiceEngine {
	return &uniswapv3PoolFactoryEngine{
		poolFactoryContract: pooFactoryContract,
		debugger:            debugger.NewDebugger("poolFactEngine", logLevel),
	}
}
