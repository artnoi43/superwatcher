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
	testCasesPoolFactoryV2 = []servicetest.TestCase{
		// 16054014
		// 16054066
		// 16054117
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
					ReorgTrigger: 16054016,
					ReorgBlock:   16054014,
					MovedLogs: map[uint64][]reorgsim.MoveLogs{
						16054014: {
							{
								NewBlock: 16054066,
								TxHashes: []common.Hash{
									common.HexToHash("0x9dab332b40f5c6689f9ae13d5b7bd0f1f7dbda2bd7d5ca8045d0670d1d5abe5c"),
								},
							},
						},
					},
				},
				{
					ReorgTrigger: 16054070,
					ReorgBlock:   16054066,
					MovedLogs: map[uint64][]reorgsim.MoveLogs{
						16054066: {
							{
								NewBlock: 16054117,
								TxHashes: []common.Hash{
									common.HexToHash("0x9dab332b40f5c6689f9ae13d5b7bd0f1f7dbda2bd7d5ca8045d0670d1d5abe5c"),
									common.HexToHash("0x4616b90e52ecebf8405179f6505eaa54eeb5182a5f88c9f0873dc8b33d77b6bd"),
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestServiceEnginePoolFactoryV2(t *testing.T) {
	err := testutils.RunTestCase(t, "TestServiceEnginePoolFactoryV2", testCasesPoolFactoryV2, testServiceEnginePoolFactoryV2)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testServiceEnginePoolFactoryV2(t *testing.T, caseNumber int) error {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		// superwatcher.PolicyExpensive,
	} {
		testCase := testCasesPoolFactoryV2[caseNumber-1]
		testCase.Policy = policy

		b, _ := json.Marshal(testCase)
		t.Logf("testServiceEnginePoolFactoryV1 case %d, policy %s, testCase %s", caseNumber, policy.String(), b)

		serviceDataGateway := datagateway.NewMockDataGatewayPoolFactory()
		stateDataGateway, err := runPoolFactory(testCase, serviceDataGateway)
		if err != nil {
			lastRecordedBlock, _ := stateDataGateway.GetLastRecordedBlock(nil)
			return errors.Wrapf(err, "error in full servicetest (poolfactory), lastRecordedBlock %d", lastRecordedBlock)
		}
		// Test if moved logs were properly removed from their parking blocks
		movedHashes, logsPark, logsDst := reorgsim.LogsReorgPaths(testCase.Events)
		debugDB := serviceDataGateway.(datagateway.DebugDataGateway)
		for _, txHash := range movedHashes {
			parks := logsPark[txHash]

			if err := findDeletionFromParks(parks, debugDB); err != nil {
				t.Error(err.Error())
			}
		}

		pools, err := serviceDataGateway.GetPools(nil)
		if err != nil {
			t.Errorf("GetPools failed after service returned: %s", err.Error())
		}

		// Test if final results are correct
		for _, pool := range pools {
			if pool == nil {
				t.Fatal("nil ens result")
			}

			var reorgIndex int
			for i, event := range testCase.Events {
				if pool.BlockCreated >= event.ReorgBlock {
					reorgIndex = i
				}
			}

			// Don't check unreorged hash
			if reorgIndex != 0 {
				expectedHash := reorgsim.ReorgHash(pool.BlockCreated, reorgIndex)

				if pool.BlockHash != expectedHash {
					t.Errorf(
						"unexpected ens blockHash for block %d (reorgIndex %d), expecting %s, got %s",
						pool.BlockCreated, reorgIndex, expectedHash.String(), pool.BlockHash,
					)
				}
			}

			// Test if moved logs end up with correct values
			dst, ok := logsDst[pool.TxHash]
			if !ok {
				// Logs was not moved
				continue
			}

			if pool.BlockCreated != dst {
				t.Errorf(
					"invalid block number for pool %s (txHash %s) - expecting %d, got %d",
					pool.Address.String(), pool.TxHash, dst, pool.BlockCreated,
				)
			}
		}
	}

	return nil
}
