package servicetest

import (
	"testing"

	"github.com/artnoi43/gsl/gslutils"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/routerengine"
)

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
		dgwPoolFactory := datagateway.NewMockDataGatewayPoolFactory()
		router := routerengine.NewMockRouter(logLevel, dgwENS, dgwPoolFactory)

		conf := &config.EmitterConfig{
			// We use fakeRedis and fakeEthClient, so no need for token strings.
			Chain:         string(enums.ChainEthereum),
			StartBlock:    testCase.startBlock,
			FilterRange:   10,
			GoBackRetries: 2,
			LoopInterval:  0,
			LogLevel:      logLevel,
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

			if result.BlockNumber < testCase.reorgBlock {
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

			if result.BlockCreated < testCase.reorgBlock {
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
