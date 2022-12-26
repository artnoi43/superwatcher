package demotest

import (
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines/ensengine"
)

// TestServiceEngineENSV1 is full tests for SubEngineENS with only 1 reorg event.
func TestServiceEngineENSV1(t *testing.T) {
	logsPath := testLogsPath + "/ens"
	testCases := []servicetest.TestCase{
		{
			LogsFiles: []string{
				logsPath + "/logs_reorg_test.json",
			},
			DataGatewayFirstRun: false, // Normal run
			Param: reorgsim.BaseParam{
				StartBlock:    15984000,
				BlockProgress: 20,
				ExitBlock:     15984200,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 15984033,
				},
			},
		},
		{
			LogsFiles: []string{
				logsPath + "/logs_servicetest_16054000_16054100.json",
			},
			DataGatewayFirstRun: false,
			Param: reorgsim.BaseParam{
				StartBlock:    16054000,
				BlockProgress: 20,
				ExitBlock:     16054200,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 16054078,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("testCase for ENS: %+v", testCase)
		// We'll later use |ensStore| to check for saved results
		ensStore := datagateway.NewMockDataGatewayENS()

		fakeRedis, err := testServiceEngineENSV1(testCase, ensStore)
		if err != nil {
			lastRecordedBlock, _ := fakeRedis.GetLastRecordedBlock(nil)
			t.Errorf("lastRecordedBlock %d - error in full servicetest (ens): %s", lastRecordedBlock, err.Error())
		}

		results, err := ensStore.GetENSes(nil)
		if err != nil {
			t.Error("error from ensStore (ens):", err.Error())
		}

		for _, result := range results {
			if result.BlockNumber >= testCase.Events[0].ReorgBlock {
				t.Log("checking block", result.BlockNumber)

				expectedHash := gslutils.StringerToLowerString(
					reorgsim.ReorgHash(result.BlockNumber, 0),
				)

				if result.BlockHash != expectedHash {
					t.Fatalf("unexpected block %d hash (ens): expecting %s, got %s", result.BlockNumber, expectedHash, result.BlockHash)
				}

				if result.ID == "" {
					t.Fatal("empty ENS ID -- should not happen")
				}
				if result.Name == "" {
					t.Fatal("empty ENS Name -- should not happen")
				}
			}
		}
	}
}

func testServiceEngineENSV1(
	testCase servicetest.TestCase,
	ensStore datagateway.RepositoryENS,
) (
	superwatcher.GetStateDataGateway,
	error,
) {
	ensEngine := ensengine.NewTestSuiteENS(ensStore, 2).Engine

	components := servicetest.InitTestComponents(
		servicetest.DefaultServiceTestConfig(testCase.Param.StartBlock, 4),
		ensEngine,
		testCase.Param,
		testCase.Events,
		testCase.LogsFiles,
		testCase.DataGatewayFirstRun,
	)

	return servicetest.RunServiceTestComponents(components)
}
