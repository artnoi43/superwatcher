package demotest

import (
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/ensengine"
)

func TestServiceEngineENS(t *testing.T) {
	logsPath := "../../test_logs/ens"
	testCases := []servicetest.TestCase{
		{
			StartBlock: 15984000,
			ReorgBlock: 15984033,
			ExitBlock:  15984100,
			LogsFiles: []string{
				logsPath + "/logs_reorg_test.json",
			},
			DataGatewayFirstRun: false, // Normal run
		},
		{
			StartBlock: 16054000,
			ReorgBlock: 16054078,
			ExitBlock:  16054100,
			LogsFiles: []string{
				logsPath + "/logs_servicetest_16054000_16054100.json",
			},
			DataGatewayFirstRun: false,
		},
	}

	for _, testCase := range testCases {
		t.Logf("testCase for ENS: %+v", testCase)
		// We'll later use |ensStore| to check for saved results
		ensStore := datagateway.NewMockDataGatewayENS()

		fakeRedis, err := testServiceEngineENS(testCase, ensStore)
		if err != nil {
			lastRecordedBlock, _ := fakeRedis.GetLastRecordedBlock(nil)
			t.Errorf("lastRecordedBlock %d - error in full servicetest (ens): %s", lastRecordedBlock, err.Error())
		}

		results, err := ensStore.GetENSes(nil)
		if err != nil {
			t.Error("error from ensStore (ens):", err.Error())
		}

		for _, result := range results {
			if result.BlockNumber >= testCase.ReorgBlock {
				t.Log("checking block", result.BlockNumber)
				// Since reorged block uses hash from deterministic PRandomHash,
				// we can check for equality this way
				expectedHash := gslutils.StringerToLowerString(reorgsim.PRandomHash(result.BlockNumber))

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

func testServiceEngineENS(
	testCase servicetest.TestCase,
	ensStore datagateway.RepositoryENS,
) (
	superwatcher.GetStateDataGateway,
	error,
) {
	conf := &config.EmitterConfig{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		StartBlock:    testCase.StartBlock,
		FilterRange:   10,
		GoBackRetries: 2,
		LoopInterval:  0,
		LogLevel:      4,
	}

	ensEngine := ensengine.NewTestSuiteENS(ensStore, 2).Engine

	components, _ := servicetest.InitTestComponents(
		conf,
		ensEngine,
		testCase.LogsFiles,
		testCase.StartBlock,
		testCase.ReorgBlock,
		testCase.ReorgBlock+conf.FilterRange*conf.GoBackRetries,
		testCase.DataGatewayFirstRun,
	)

	return servicetest.RunServiceTestComponents(components)
}
