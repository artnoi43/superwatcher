package servicetest

import (
	"testing"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/routerengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/ensengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/uniswapv3factoryengine"
	"github.com/ethereum/go-ethereum/common"
)

func TestFoo(t *testing.T) {
	logsPath := "../assets/servicetest"
	testCases := []testCase{
		{
			startBlock: 16054000,
			reorgBlock: 16054078,
			exitBlock:  16054100,
			logsFiles: []string{
				logsPath + "/logs_servicetest_16054000_16054100.json",
			},
		},
	}

	for _, testCase := range testCases {
		conf := &config.EmitterConfig{
			// We use fakeRedis and fakeEthClient, so no need for token strings.
			Chain:         string(enums.ChainEthereum),
			StartBlock:    testCase.startBlock,
			FilterRange:   10,
			GoBackRetries: 2,
			LoopInterval:  0,
		}

		components, param := initTestComponents(
			conf,
			&engine{
				reorgedAt:          testCase.reorgBlock,
				emitterFilterRange: conf.FilterRange,
				debugger:           debugger.NewDebugger("serviceTest", 4),
			},
			testCase.logsFiles,
			testCase.startBlock,
			testCase.reorgBlock,
			testCase.exitBlock,
		)

		err := serviceEngineTestTemplate(components, param)
		if err != nil {
			t.Error("error in full servicetest (ens):", err.Error())
		}

	}
}

func TestServiceEngineRouter(t *testing.T) {
	logsPath := "../assets/servicetest"
	testCases := []testCase{
		{
			startBlock: 16054000,
			reorgBlock: 16054078,
			exitBlock:  16054100,
			logsFiles: []string{
				logsPath + "/logs_servicetest_16054000_16054100.json",
			},
		},
	}

	logLevel := uint8(3)

	for _, testCase := range testCases {
		dgwENS := datagateway.NewMockDataGatewayENS()
		suiteENS := ensengine.NewTestSuiteENS(dgwENS, logLevel)
		dgwPoolFactory := datagateway.NewMockDataGatewayPoolFactory()
		suitePoolFactory := uniswapv3factoryengine.NewTestSuitePoolFactory(dgwPoolFactory, logLevel)

		routes := make(map[subengines.SubEngineEnum]map[common.Address][]common.Hash)
		routes[subengines.SubEngineENS] = suiteENS.Routes[subengines.SubEngineENS]
		routes[subengines.SubEngineUniswapv3Pool] = suiteENS.Routes[subengines.SubEngineUniswapv3Pool]

		services := map[subengines.SubEngineEnum]superwatcher.ServiceEngine{
			subengines.SubEngineENS:           suiteENS.Engine,
			subengines.SubEngineUniswapv3Pool: suitePoolFactory.Engine,
		}

		router := routerengine.New(routes, services, logLevel)

		conf := &config.EmitterConfig{
			// We use fakeRedis and fakeEthClient, so no need for token strings.
			Chain:         string(enums.ChainEthereum),
			StartBlock:    testCase.startBlock,
			FilterRange:   10,
			GoBackRetries: 2,
			LoopInterval:  0,
		}

		components, param := initTestComponents(
			conf,
			router,
			testCase.logsFiles,
			testCase.startBlock,
			testCase.reorgBlock,
			testCase.exitBlock,
		)

		err := serviceEngineTestTemplate(components, param)
		if err != nil {
			t.Error("error in full servicetest (ens):", err.Error())
		}
	}
}
