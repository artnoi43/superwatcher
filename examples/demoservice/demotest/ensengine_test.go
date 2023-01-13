package demotest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/artnoi43/gsl"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"
	"github.com/artnoi43/superwatcher/pkg/testutils"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines/ensengine"
)

var (
	ensCasesAlreadyRunV1 = false
	logsPathENS          = testLogsPath + "/ens"
	testCasesENSV1       = []servicetest.TestCase{
		{
			LogsFiles: []string{
				logsPathENS + "/logs_reorg_test.json",
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
				logsPathENS + "/logs_servicetest_16054000_16054100.json",
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
)

func invalidateNullFieldsENS(ens *entity.ENS) error {
	if ens.ID == "" {
		return errors.New("ens.ID is null")
	}

	if ens.Name == "" {
		return fmt.Errorf("ens %s Name is null (txHash %s)", ens.ID, ens.TxHash)
	}

	if ens.BlockNumber == 0 {
		return fmt.Errorf("ens %s (Name %s) BlockNumber is 0 (txHash %s)", ens.ID, ens.Name, ens.TxHash)
	}

	return nil
}

// TestServiceEngineENSV1 is full tests for SubEngineENS with only 1 reorg event.
func TestServiceEngineENSV1(t *testing.T) {
	err := testutils.RunTestCase(t, "TestServiceEngineENSV1", testCasesENSV1, testServiceEngineENSV1)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testServiceEngineENSV1(t *testing.T, caseNumber int) error {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		// superwatcher.PolicyExpensive,
	} {
		testCase := testCasesENSV1[caseNumber-1]
		testCase.Policy = policy
		b, _ := json.Marshal(testCase)
		t.Logf("testServiceEngineENSV1 case %d, testCase %s", caseNumber, b)

		// We'll later use |ensStore| to check for saved results
		ensStore := datagateway.NewMockDataGatewayENS()
		fakeRedis, err := runENS(testCase, ensStore)
		if err != nil {
			lastRecordedBlock, _ := fakeRedis.GetLastRecordedBlock(nil)
			return errors.Wrapf(err, "error in full servicetest (ens), lastRecordedBlock %d", lastRecordedBlock)
		}

		results, err := ensStore.GetENSes(nil)
		if err != nil {
			return errors.Wrap(err, "GetENSes failed after servicetest")
		}

		// Test if moved logs were properly removed
		movedHashes, logsPark, logsDst := reorgsim.LogsReorgPaths(testCase.Events)
		debugDB := ensStore.(datagateway.DebugDataGateway)
		for _, txHash := range movedHashes {
			parks := logsPark[txHash]

			if err := findDeletionFromParks(parks, debugDB); err != nil {
				t.Error(err.Error())
			}
		}

		for _, result := range results {
			if result.BlockNumber >= testCase.Events[0].ReorgBlock {

				t.Log("checking block", result.BlockNumber)

				expectedHash := gsl.StringerToLowerString(
					reorgsim.ReorgHash(result.BlockNumber, 0),
				)

				if result.BlockHash != expectedHash {
					t.Fatalf("unexpected block %d hash (ens): expecting %s, got %s", result.BlockNumber, expectedHash, result.BlockHash)
				}

				if err := invalidateNullFieldsENS(result); err != nil {
					t.Error("result has invalid ENS values", err.Error())
				}

				if h := common.HexToHash(result.TxHash); gsl.Contains(movedHashes, h) {
					expectedFinalBlock := logsDst[h]
					if expectedFinalBlock != result.BlockNumber {
						t.Fatalf("expecting moved blockNumber %d, got %d", expectedFinalBlock, result.BlockNumber)
					}
				}
			}
		}
	}

	return nil
}

func runENS(
	testCase servicetest.TestCase,
	ensStore datagateway.RepositoryENS,
) (
	superwatcher.GetStateDataGateway,
	error,
) {
	ensEngine := ensengine.NewTestSuiteENS(ensStore, 2).Engine

	components := servicetest.InitTestComponents(
		servicetest.DefaultServiceTestConfig(testCase.Param.StartBlock, 4, testCase.Policy),
		ensEngine,
		testCase.Param,
		testCase.Events,
		testCase.LogsFiles,
		testCase.DataGatewayFirstRun,
	)

	return servicetest.RunServiceTestComponents(components)
}
