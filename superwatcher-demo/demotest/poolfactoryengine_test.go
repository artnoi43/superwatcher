package demotest

import (
	"testing"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/uniswapv3factoryengine"
)

func TestServiceEnginePoolFactory(t *testing.T) {
	logsPath := "../assets/poolfactory"
	testCases := []servicetest.TestCase{
		{
			StartBlock: 16054000,
			ReorgBlock: 16054066,
			ExitBlock:  16054200,
			LogsFiles: []string{
				logsPath + "/logs_reorg_test.json",
			},
		},
	}

	for _, testCase := range testCases {
		serviceDataGateway := datagateway.NewMockDataGatewayPoolFactory()
		stateDataGateway, err := testServiceEnginePoolFactory(
			testCase.StartBlock,
			testCase.ReorgBlock,
			testCase.ExitBlock,
			testCase.LogsFiles,
			serviceDataGateway,
		)
		if err != nil {
			lastRecordedBlock, _ := stateDataGateway.GetLastRecordedBlock(nil)
			t.Errorf("lastRecordedBlock: %d error in full servicetest (poolfactory): %s", lastRecordedBlock, err.Error())
		}

		results, err := serviceDataGateway.GetPools(nil)
		if err != nil {
			t.Errorf("GetPools failed after service returned: %s", err.Error())
		}

		for _, result := range results {
			expectedReorgedHash := reorgsim.PRandomHash(result.BlockCreated)
			if result.BlockCreated >= testCase.ReorgBlock {
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
	lpStore datagateway.RepositoryPoolFactory,
) (
	superwatcher.StateDataGateway,
	error,
) {
	conf := &config.EmitterConfig{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		StartBlock:    startBlock,
		FilterRange:   10,
		GoBackRetries: 2,
		LoopInterval:  0,
	}

	poolFactoryEngine := uniswapv3factoryengine.NewTestSuitePoolFactory(lpStore, 2).Engine
	components, _ := servicetest.InitTestComponents(
		conf,
		poolFactoryEngine,
		logsFiles,
		conf.StartBlock,
		reorgedAt,
		exitBlock,
	)

	return servicetest.RunServiceTestComponents(components)
}
