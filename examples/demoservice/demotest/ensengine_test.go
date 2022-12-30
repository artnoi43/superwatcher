package demotest

import (
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"
	"github.com/ethereum/go-ethereum/common"

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
			Param: reorgsim.Param{
				StartBlock:    15984000,
				BlockProgress: 20,
				ExitBlock:     15984200,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 15984040,
					MovedLogs: map[uint64][]reorgsim.MoveLogs{
						15984040: {
							{
								NewBlock: 15984043,
								TxHashes: []common.Hash{
									common.HexToHash("0xd5d1beffbfe5fbf4d8dee6284d291a0a11260085c9fc6074e56ca4ed44491263"),
								},
							},
						},
					},
				},
			},
		},
		{
			LogsFiles: []string{
				logsPath + "/logs_servicetest_16054000_16054100.json",
			},
			DataGatewayFirstRun: false,
			Param: reorgsim.Param{
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

		// Test if moved logs were properly removed
		movedHashes, logsPark, logsDst := reorgsim.LogsReorgPaths(testCase.Events)
		ensMockDB := ensStore.(*datagateway.MockDataGatewayENS)
		for _, txHash := range movedHashes {
			parks := logsPark[txHash]

			for _, park := range parks {
				var foundDel bool
				for _, writeLog := range ensMockDB.WriteLogs {
					method, _, blockNumber, _, err := writeLog.Unmarshal()
					if err != nil {
						t.Fatal("bad writeLog", err.Error())
					}

					if method != "DEL" {
						continue
					}

					if blockNumber == park {
						foundDel = true
					}
				}

				if !foundDel {
					t.Errorf("moved log did not produce writeLog DEL for txHash %s", txHash.String())
					t.Log(ensMockDB.WriteLogs)
				}
			}
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

				if h := common.HexToHash(result.TxHash); gslutils.Contains(movedHashes, h) {
					expectedFinalBlock := logsDst[h]
					if expectedFinalBlock != result.BlockNumber {
						t.Fatalf("expecting moved blockNumber %d, got %d", expectedFinalBlock, result.BlockNumber)
					}
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
