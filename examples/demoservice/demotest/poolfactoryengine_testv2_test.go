package demotest

import (
	"testing"

	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
)

func TestServiceEnginePoolFactoryV2(t *testing.T) {
	logsPath := testLogsPath + "/poolfactory"
	testCases := []servicetest.TestCase{
		// 16054014
		// 16054066
		// 16054117
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

	for _, testCase := range testCases {
		serviceDataGateway := datagateway.NewMockDataGatewayPoolFactory()
		stateDataGateway, err := testServiceEnginePoolFactoryV1(testCase, serviceDataGateway)
		if err != nil {
			lastRecordedBlock, _ := stateDataGateway.GetLastRecordedBlock(nil)
			t.Errorf("lastRecordedBlock: %d error in full servicetest (poolfactory): %s", lastRecordedBlock, err.Error())
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
}
