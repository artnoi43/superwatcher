package demotest

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/servicetest"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
)

func TestServiceEngineENSV2(t *testing.T) {
	logsPath := testLogsPath + "/ens"
	testCases := []servicetest.TestCase{
		{
			LogsFiles: []string{
				logsPath + "/logs_servicetest_16054000_16054100.json",
			},
			DataGatewayFirstRun: false, // Normal run
			Param: reorgsim.Param{
				StartBlock:    16054000,
				BlockProgress: 20,
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

	for _, testCase := range testCases {
		ensStore := datagateway.NewMockDataGatewayENS()
		fakeRedis, err := runTestServiceEngineENS(testCase, ensStore)
		if err != nil {
			t.Error(err.Error())
		}

		// Test if moved logs were properly removed
		movedHashes, logsPark, _ := reorgsim.LogsReorgPaths(testCase.Events)
		debugDB := ensStore.(datagateway.DebugDataGateway)
		for _, txHash := range movedHashes {
			parks := logsPark[txHash]

			for _, park := range parks {
				var foundDel bool
				for _, writeLog := range debugDB.WriteLogs() {
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
				}
			}
		}

		lastRec, err := fakeRedis.GetLastRecordedBlock(nil)
		t.Log(lastRec)
	}
}
