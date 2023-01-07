package demotest

import (
	"testing"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines/uniswapv3factoryengine"
)

func TestServiceEnginePoolFactoryV1(t *testing.T) {
	logsPath := testLogsPath + "/poolfactory"
	testCases := []servicetest.TestCase{
		{
			LogsFiles: []string{
				logsPath + "/logs_reorg_test.json",
			},
			DataGatewayFirstRun: false,
			Param: reorgsim.Param{
				StartBlock:    16054000,
				ExitBlock:     16054200,
				BlockProgress: 20,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 16054066,
				},
			},
		},
	}

	for _, pollLevel := range []superwatcher.PollLevel{
		superwatcher.PollLevelFast,
		superwatcher.PollLevelNormal,
		superwatcher.PollLevelExpensive,
	} {
		for _, testCase := range testCases {
			testCase.PollLevel = pollLevel
			serviceDataGateway := datagateway.NewMockDataGatewayPoolFactory()
			stateDataGateway, err := testServiceEnginePoolFactoryV1(testCase, serviceDataGateway)
			if err != nil {
				lastRecordedBlock, _ := stateDataGateway.GetLastRecordedBlock(nil)
				t.Errorf("lastRecordedBlock: %d error in full servicetest (poolfactory): %s", lastRecordedBlock, err.Error())
			}

			results, err := serviceDataGateway.GetPools(nil)
			if err != nil {
				t.Errorf("GetPools failed after service returned: %s", err.Error())
			}

			for _, result := range results {
				expectedReorgedHash := reorgsim.ReorgHash(result.BlockCreated, 0)
				if result.BlockCreated >= testCase.Events[0].ReorgBlock {
					if result.BlockHash != expectedReorgedHash {
						t.Fatalf(
							"unexpected reorgedBlockHash - expecting %s, got %s",
							expectedReorgedHash.String(), result.BlockHash.String(),
						)
					}

					continue
				}

				if result.BlockHash == expectedReorgedHash {
					t.Fatal("old block in the old chain has reorged hash")
				}
			}
		}
	}
}

func testServiceEnginePoolFactoryV1(
	testCase servicetest.TestCase,
	lpStore datagateway.RepositoryPoolFactory,
) (
	superwatcher.GetStateDataGateway,
	error,
) {
	poolFactoryEngine := uniswapv3factoryengine.NewTestSuitePoolFactory(lpStore, 2).Engine

	components := servicetest.InitTestComponents(
		servicetest.DefaultServiceTestConfig(testCase.Param.StartBlock, 3, testCase.PollLevel),
		poolFactoryEngine,
		testCase.Param,
		testCase.Events,
		testCase.LogsFiles,
		testCase.DataGatewayFirstRun,
	)

	return servicetest.RunServiceTestComponents(components)
}
