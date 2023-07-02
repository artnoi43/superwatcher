package demotest

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/soyart/superwatcher"
	"github.com/soyart/superwatcher/pkg/reorgsim"
	"github.com/soyart/superwatcher/pkg/servicetest"
	"github.com/soyart/superwatcher/pkg/testutils"

	"github.com/soyart/superwatcher/examples/demoservice/internal/domain/datagateway"
)

var (
	logsPathENSV2  = testLogsPath + "/ens"
	testCasesENSV2 = []servicetest.TestCase{
		{
			LogsFiles: []string{
				logsPathENSV2 + "/logs_servicetest_16054000_16054100.json",
			},
			DataGatewayFirstRun: false, // Normal run
			Param: reorgsim.Param{
				StartBlock:    16054000,
				BlockProgress: 17,
				ExitBlock:     16054600,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgTrigger: 16054020,
					ReorgBlock:   16054014,
					MovedLogs: map[uint64][]reorgsim.MoveLogs{
						16054014: {
							{
								NewBlock: 16054022,
								TxHashes: []common.Hash{
									common.HexToHash("0x7a3bdb4ec3bef7a532a7b215fffb147c05d828750cd601ebc8e3959ab6e2d8b1"),
								},
							},
						},
					},
				},
				{
					ReorgTrigger: 16054035,
					ReorgBlock:   16054026,
					MovedLogs: map[uint64][]reorgsim.MoveLogs{
						16054027: {
							{
								NewBlock: 16054026,
								TxHashes: []common.Hash{
									common.HexToHash("0x2de80c99259ac459a0b5f557858fe5f5fc1c94092b14d9cbed0d4d7636d6eb55"),
								},
							},
						},
					},
				},
				{
					ReorgBlock: 16054035,
					MovedLogs: map[uint64][]reorgsim.MoveLogs{
						16054035: {
							{
								NewBlock: 16054047,
								TxHashes: []common.Hash{
									common.HexToHash("0x96bf497e7521d389a07d9735ca1518d72c6ceead69b3f6f6fef371a97fb3a398"),
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestServiceEngineENSV2(t *testing.T) {
	err := testutils.RunTestCase(t, "TestServiceEngineENSV2", testCasesENSV2, testServiceEngineENSV2)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testServiceEngineENSV2(t *testing.T, caseNumber int) error {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		superwatcher.PolicyExpensive,
	} {
		testCase := testCasesENSV1[caseNumber-1]
		testCase.Policy = policy
		b, _ := json.Marshal(testCase)
		t.Logf("testServiceEngineENSV1 case %d, testCase %s", caseNumber, b)

		ensStore := datagateway.NewMockDataGatewayENS()

		stateDgw, err := runENS(testCase, ensStore)
		if err != nil {
			lastRecordedBlock, _ := stateDgw.GetLastRecordedBlock(nil)
			return errors.Wrapf(err, "error in servicetest test, lastRecordedBlock %d", lastRecordedBlock)
		}

		// Test if moved logs were properly removed from their parking blocks
		movedHashes, logsPark, logsDst := reorgsim.LogsReorgPaths(testCase.Events)
		debugDB := ensStore.(datagateway.DebugDataGateway)
		for _, txHash := range movedHashes {
			parks := logsPark[txHash]

			if err := findDeletionFromParks(parks, debugDB); err != nil {
				t.Error(err.Error())
			}
		}

		results, err := ensStore.GetENSes(nil)
		if err != nil {
			return errors.Wrap(err, "ensStore.GetENSes failed after servicetest")
		}

		// Test if final results are correct
		for _, ens := range results {
			if ens == nil {
				t.Fatal("nil ens result")
			}

			var reorgIndex int
			for i, event := range testCase.Events {
				if ens.BlockNumber >= event.ReorgBlock {
					reorgIndex = i
				}
			}

			// Don't check unreorged hash
			if reorgIndex != 0 {
				expectedHash := reorgsim.ReorgHash(ens.BlockNumber, reorgIndex)
				ensBlockHash := common.HexToHash(ens.BlockHash)

				if ensBlockHash != expectedHash {
					t.Errorf(
						"unexpected ens blockHash for block %d (reorgIndex %d), expecting %s, got %s",
						ens.BlockNumber, reorgIndex, expectedHash.String(), ens.BlockHash,
					)
				}

				if err := invalidateNullFieldsENS(ens); err != nil {
					t.Error("ens has invalid values", err.Error())
				}
			}

			// Test if moved logs end up with correct values
			txHash := common.HexToHash(ens.TxHash)
			dst, ok := logsDst[txHash]
			if !ok {
				// Logs was not moved
				continue
			}

			if ens.BlockNumber != dst {
				t.Errorf(
					"invalid block number for ENS %s (txHash %s) - expecting %d, got %d",
					ens.ID, ens.TxHash, dst, ens.BlockNumber,
				)
			}
		}
	}

	return nil
}
