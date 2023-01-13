package uniswapv3factoryengine

import (
	"github.com/artnoi43/w3utils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/hardcode"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines"
)

type uniswapv3PoolFactoryEngine struct {
	poolFactoryContract w3utils.Contract
	dataGateway         datagateway.RepositoryPoolFactory
	debugger            *debugger.Debugger
}

type TestSuitePoolFactory struct {
	Engine   superwatcher.ServiceEngine // *ensEngine
	Routes   map[subengines.SubEngineEnum]map[common.Address][]common.Hash
	Services map[subengines.SubEngineEnum]superwatcher.ServiceEngine
}

func New(
	pooFactoryContract w3utils.Contract,
	dgw datagateway.RepositoryPoolFactory,
	logLevel uint8,
) superwatcher.ServiceEngine {
	return &uniswapv3PoolFactoryEngine{
		poolFactoryContract: pooFactoryContract,
		dataGateway:         dgw,
		debugger:            debugger.NewDebugger("poolFactEngine", logLevel),
	}
}

// NewTestSuitePoolFactory returns a convenient struct for injecting into routerengine.routerEngine
func NewTestSuitePoolFactory(dgw datagateway.RepositoryPoolFactory, logLevel uint8) *TestSuitePoolFactory {
	poolFactoryContract := hardcode.DemoContracts(hardcode.Uniswapv3Factory)[hardcode.Uniswapv3Factory]
	poolFactoryTopics := w3utils.CollectEventHashes(poolFactoryContract.ContractEvents)
	poolFactoryEngine := New(poolFactoryContract, dgw, logLevel)

	return &TestSuitePoolFactory{
		Engine: poolFactoryEngine,
		Routes: map[subengines.SubEngineEnum]map[common.Address][]common.Hash{
			subengines.SubEngineUniswapv3Factory: {
				poolFactoryContract.Address: poolFactoryTopics,
			},
		},
		Services: map[subengines.SubEngineEnum]superwatcher.ServiceEngine{
			subengines.SubEngineUniswapv3Factory: poolFactoryEngine,
		},
	}
}
