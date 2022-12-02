package servicetest

import (
	"testing"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/uniswapv3factoryengine"
)

func TestServiceEnginePoolFactory(t *testing.T) {
	logsPath := "../assets/poolfactory"
	testCases := []testCase{
		{
			startBlock: 16054000,
			reorgBlock: 16054066,
			exitBlock:  16054200,
			logsFiles: []string{
				logsPath + "/logs_reorg_test.json",
			},
		},
	}

	for _, testCase := range testCases {
		dgw := datagateway.NewMockDataGatewayPoolFactory()
		if err := testServiceEnginePoolFactory(
			testCase.startBlock,
			testCase.reorgBlock,
			testCase.exitBlock,
			testCase.logsFiles,
			dgw,
		); err != nil {
			t.Errorf("error in full servicetest (poolfactory): %s", err.Error())
		}

		results, err := dgw.GetPools(nil)
		if err != nil {
			t.Errorf("GetPools failed after service returned: %s", err.Error())
		}

		for _, result := range results {
			expectedReorgedHash := reorgsim.PRandomHash(result.BlockCreated)
			if result.BlockCreated >= testCase.reorgBlock {
				if result.BlockHash != expectedReorgedHash {
					t.Fatalf("blockHash not reorged")
				}

				continue
			}

			if result.BlockHash == expectedReorgedHash {
				t.Fatal("old block in the old chain has reorged hash")
			}
		}
	}
}

func testServiceEnginePoolFactory(
	startBlock uint64,
	reorgedAt uint64,
	exitBlock uint64,
	logsFiles []string,
	lpStore datagateway.DataGatewayPoolFactory,
) error {
	conf := &config.EmitterConfig{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		StartBlock:    startBlock,
		FilterRange:   10,
		GoBackRetries: 2,
		LoopInterval:  0,
	}

	poolFactoryEngine := uniswapv3factoryengine.NewTestSuitePoolFactory(lpStore, 2).Engine
	components, param := initTestComponents(
		conf,
		poolFactoryEngine,
		logsFiles,
		conf.StartBlock,
		reorgedAt,
		exitBlock,
	)

	return serviceEngineTestTemplate(components, param)
}
