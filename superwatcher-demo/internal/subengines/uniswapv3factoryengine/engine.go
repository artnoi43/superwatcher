package uniswapv3factoryengine

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
)

type uniswapv3PoolFactoryEngine struct {
	poolFactoryContract contracts.BasicContract
	dataGateway         datagateway.DataGatewayPoolFactory
	debugger            *debugger.Debugger
}

func New(
	pooFactoryContract contracts.BasicContract,
	dgw datagateway.DataGatewayPoolFactory,
	logLevel uint8,
) superwatcher.ServiceEngine {
	return &uniswapv3PoolFactoryEngine{
		poolFactoryContract: pooFactoryContract,
		dataGateway:         dgw,
		debugger:            debugger.NewDebugger("poolFactEngine", logLevel),
	}
}
