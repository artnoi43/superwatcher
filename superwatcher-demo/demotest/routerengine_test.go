package demotest

import (
	"testing"

	"github.com/artnoi43/gsl/gslutils"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/routerengine"
)

func TestServiceEngineRouter(t *testing.T) {
	logsPath := "../assets/servicetest"
	testCases := []servicetest.TestCase{
		{
			StartBlock: 16054000,
			ReorgBlock: 16054078,
			ExitBlock:  16054100,
			LogsFiles: []string{
				logsPath + "/logs_servicetest_16054000_16054100.json",
			},
		},
	}

	logLevel := uint8(3)
	for _, testCase := range testCases {
		dgwENS := datagateway.NewMockDataGatewayENS()
		dgwPoolFactory := datagateway.NewMockDataGatewayPoolFactory()
		router := routerengine.NewMockRouter(logLevel, dgwENS, dgwPoolFactory)

		conf := &config.EmitterConfig{
			// We use fakeRedis and fakeEthClient, so no need for token strings.
			StartBlock:    testCase.StartBlock,
			FilterRange:   10,
			GoBackRetries: 2,
			LoopInterval:  0,
			LogLevel:      logLevel,
		}

		components, _ := servicetest.InitTestComponents(
			conf,
			router,
			testCase.LogsFiles,
			testCase.StartBlock,
			testCase.ReorgBlock,
			testCase.ExitBlock,
		)

		_, err := servicetest.RunServiceTestComponents(components)
		if err != nil {
			t.Error("error in full servicetest (ens):", err.Error())
		}

		resultsENS, err := dgwENS.GetENSes(nil)
		if err != nil {
			t.Errorf("error getting results from dgwENS: %s", err.Error())
		}
		if len(resultsENS) == 0 {
			t.Fatalf("0 results from dgwENS")
		}
		resultsPoolFactory, err := dgwPoolFactory.GetPools(nil)
		if err != nil {
			t.Errorf("error getting results from dgwPoolFactory: %s", err.Error())
		}
		if len(resultsPoolFactory) == 0 {
			t.Fatalf("0 results from dgwPoolFactory")
		}

		for _, result := range resultsENS {
			if result.DomainString() == "" {
				t.Errorf("emptyDomain name for resultENS id: %s", result.ID)
			}

			expectedReorgedHash := gslutils.StringerToLowerString(reorgsim.PRandomHash(result.BlockNumber))

			if result.BlockNumber < testCase.ReorgBlock {
				if result.BlockHash == expectedReorgedHash {
					t.Errorf("good block resultENS has reorged blockHash: %s", expectedReorgedHash)
				}

				continue
			}

			if result.BlockHash != expectedReorgedHash {
				t.Errorf("reorged block resultENS has unexpected blockHash: %s", result.BlockHash)
			}
		}

		for _, result := range resultsPoolFactory {
			expectedReorgedHash := gslutils.StringerToLowerString(reorgsim.PRandomHash(result.BlockCreated))
			resultBlockHash := gslutils.StringerToLowerString(result.BlockHash)

			if result.BlockCreated < testCase.ReorgBlock {
				if resultBlockHash == expectedReorgedHash {
					t.Errorf("good block resultPoolFactory has reorged blockHash: %s", resultBlockHash)
				}

				continue
			}

			if resultBlockHash != expectedReorgedHash {
				t.Errorf("reorged block resultPoolFactory has unexpected blockHash: %s", resultBlockHash)
			}
		}
	}
}
