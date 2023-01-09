package demotest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"
	"github.com/artnoi43/superwatcher/pkg/testutils"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines/uniswapv3factoryengine"
)

var (
	logsPathPoolFactory    = testLogsPath + "/poolfactory"
	testCasesPoolFactoryV1 = []servicetest.TestCase{
		{
			LogsFiles: []string{
				logsPathPoolFactory + "/logs_reorg_test.json",
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
)

func TestServiceEnginePoolFactoryV1(t *testing.T) {
	err := testutils.RunTestCase(t, "testServiceEngineENSV1", testCasesPoolFactoryV1, testServiceEnginePoolFactoryV1)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testServiceEnginePoolFactoryV1(t *testing.T, caseNumber int) error {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		superwatcher.PolicyExpensive,
	} {
		testCase := testCasesPoolFactoryV1[caseNumber-1]
		testCase.Policy = policy
		b, _ := json.Marshal(testCase)
		t.Logf("testServiceEnginePoolFactoryV1 case %d, testCase %s", caseNumber, b)

		serviceDataGateway := datagateway.NewMockDataGatewayPoolFactory()
		stateDataGateway, err := runPoolFactory(testCase, serviceDataGateway)
		if err != nil {
			lastRecordedBlock, _ := stateDataGateway.GetLastRecordedBlock(nil)
			return errors.Wrapf(err, "error in full servicetest (poolfactory), lastRecordedBlock: %d ", lastRecordedBlock)
		}

		results, err := serviceDataGateway.GetPools(nil)
		if err != nil {
			return errors.Wrapf(err, "GetPools failed after servicetest returned")
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
				return fmt.Errorf("old block %d in the old chain (hash %s) has reorged hash %s", result.BlockCreated, result.BlockHash.String(), expectedReorgedHash.String())
			}
		}
	}

	return nil
}

func runPoolFactory(
	testCase servicetest.TestCase,
	lpStore datagateway.RepositoryPoolFactory,
) (
	superwatcher.GetStateDataGateway,
	error,
) {
	poolFactoryEngine := uniswapv3factoryengine.NewTestSuitePoolFactory(lpStore, 2).Engine

	components := servicetest.InitTestComponents(
		servicetest.DefaultServiceTestConfig(testCase.Param.StartBlock, 3, testCase.Policy),
		poolFactoryEngine,
		testCase.Param,
		testCase.Events,
		testCase.LogsFiles,
		testCase.DataGatewayFirstRun,
	)

	return servicetest.RunServiceTestComponents(components)
}
