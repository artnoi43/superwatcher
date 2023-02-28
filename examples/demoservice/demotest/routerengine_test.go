package demotest

import (
	"testing"

	"github.com/artnoi43/gsl"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"
	"github.com/artnoi43/superwatcher/pkg/testutils"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/routerengine"
)

var (
	logsPathRouter  = testLogsPath + "/servicetest"
	testCasesRouter = []servicetest.TestCase{
		{
			LogsFiles: []string{
				logsPathRouter + "/logs_servicetest_16054000_16054100.json",
			},
			Param: reorgsim.Param{
				StartBlock:    16054000,
				BlockProgress: 20,
				ExitBlock:     16054200,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 16054078,
					MovedLogs: map[uint64][]reorgsim.MoveLogs{
						16054078: {
							{
								NewBlock: 16054093,
								TxHashes: []common.Hash{
									common.HexToHash("0xed2520a4168f1d26a8c5a0081403711415087218555cfc61fc0192432912ff1c"),
									common.HexToHash("0xa77589c6e436e85a99dbccd1cddaf13766148c740a5c0972260a4e90a742c6d5"),
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestServiceEngineRouterV1(t *testing.T) {
	err := testutils.RunTestCase(t, "TestServiceEngineRouterV1", testCasesRouter, testServiceEngineRouterV1)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testServiceEngineRouterV1(t *testing.T, caseNumber int) error {
	logLevel := uint8(3)
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		superwatcher.PolicyExpensive,
	} {
		testCase := testCasesRouter[caseNumber-1]
		testCase.Policy = policy
		dgwENS := datagateway.NewMockDataGatewayENS()
		dgwPoolFactory := datagateway.NewMockDataGatewayPoolFactory()
		router := routerengine.NewMockRouter(logLevel, dgwENS, dgwPoolFactory)

		components := servicetest.InitTestComponents(
			servicetest.DefaultServiceTestConfig(testCase.Param.StartBlock, 4, testCase.Policy),
			router,
			testCase.Param,
			testCase.Events,
			testCase.LogsFiles,
			testCase.DataGatewayFirstRun,
		)

		stateDgw, err := servicetest.RunServiceTestComponents(components)
		if err != nil {
			lastRecordedBlock, _ := stateDgw.GetLastRecordedBlock(nil)
			return errors.Wrapf(err, "error in full servicetest (ens) test case %d, lastRecordedBlock %d", caseNumber, lastRecordedBlock)
		}

		resultsENS, err := dgwENS.GetENSes(nil)
		if err != nil {
			return errors.Wrap(err, "GetENSes failed")
		}
		if len(resultsENS) == 0 {
			return errors.New("len resultsENS = 0")
		}
		resultsPoolFactory, err := dgwPoolFactory.GetPools(nil)
		if err != nil {
			t.Errorf("error getting results from dgwPoolFactory: %s", err.Error())
		}
		if len(resultsPoolFactory) == 0 {
			t.Fatalf("0 results from dgwPoolFactory")
		}

		movedHashes, _, logsDst := reorgsim.LogsReorgPaths(testCase.Events)

		var someENS bool
		for _, result := range resultsENS {
			someENS = true
			if result.DomainString() == "" {
				t.Errorf("emptyDomain name for resultENS id: %s", result.ID)
			}

			expectedReorgedHash := gsl.StringerToLowerString(reorgsim.ReorgHash(result.BlockNumber, 0))

			if result.BlockNumber < testCase.Events[0].ReorgBlock {
				if result.BlockHash == expectedReorgedHash {
					t.Errorf("good block %d resultENS has reorged blockHash: %s", result.BlockNumber, result.BlockHash)
				}

				continue
			}

			if result.BlockHash != expectedReorgedHash {
				t.Errorf("reorged block %d resultENS has unexpected blockHash: %s", result.BlockNumber, result.BlockHash)
			}

			if h := common.HexToHash(result.TxHash); gsl.Contains(movedHashes, h) {
				if expected := logsDst[h]; result.BlockNumber != expected {
					t.Fatalf("expecting moved blockNumber %d, got %d", expected, result.BlockNumber)
				}
			}
		}

		if !someENS {
			return errors.New("got no ENS result")
		}

		var somePoolFactory bool
		for _, result := range resultsPoolFactory {
			somePoolFactory = true
			expectedReorgedHash := gsl.StringerToLowerString(reorgsim.PRandomHash(result.BlockCreated))
			resultBlockHash := gsl.StringerToLowerString(result.BlockHash)

			if result.BlockCreated < testCase.Events[0].ReorgBlock {
				if resultBlockHash == expectedReorgedHash {
					t.Errorf("resultPoolFactory from good block %d has reorged blockHash: %s", result.BlockCreated, resultBlockHash)
				}

				continue
			}

			if resultBlockHash != expectedReorgedHash {
				t.Errorf("resultPoolFactory from reorged block %d has unexpected blockHash: %s", result.BlockCreated, resultBlockHash)
			}
		}

		if !somePoolFactory {
			return errors.New("got no poolFactory result")
		}
	}

	return nil
}
