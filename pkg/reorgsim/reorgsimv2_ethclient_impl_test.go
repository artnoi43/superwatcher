package reorgsim

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func TestFilterLogsV2(t *testing.T) {
	param := BaseParam{
		StartBlock:    defaultStartBlock,
		BlockProgress: 20,
	}
	event := ReorgEvent{
		ReorgBlock: defaultReorgedAt,
		MovedLogs:  nil,
	}

	sim, err := NewReorgSimV2FromLogsFiles(param, []ReorgEvent{event}, defaultLogsFiles, 4)
	if err != nil {
		t.Fatal("error creating ReorgSimV2", err.Error())
	}

	ctx := context.Background()
	logs, err := sim.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(69),
		ToBlock:   big.NewInt(70),
	})
	if err != nil {
		t.Errorf("FilterLogs returned error: %s", err.Error())
	}
	if len(logs) != 0 {
		t.Fatalf("expecing 0 logs, got %d", len(logs))
	}

	logs, err = sim.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(10000000),
		ToBlock:   big.NewInt(16000000),
	})
	if err != nil {
		t.Errorf("FilterLogs returned error: %s", err.Error())
	}
	if len(logs) == 0 {
		t.Fatalf("expecting > 0 logs, got 0 log")
	}
}

// Test if FilterLogs returns logs with correct hash (reorged hash).
// No logs are moved in this test (yet).
func TestFilterLogsReorgV2(t *testing.T) {
	for i, test := range testsReorgSimV2 {
		testFilterLogsReorgV2(t, i+1, test)
	}
}

func testFilterLogsReorgV2(t *testing.T, testCase int, testConf multiReorgConfig) error {
	if len(testConf.Events) == 0 {
		return errors.New("got 0 ReorgEvent")
	}

	sim, err := NewReorgSimV2FromLogsFiles(testConf.Param, testConf.Events, testConf.LogsFiles, 2)
	if err != nil {
		t.Fatalf("[%d] failed to create new ReorgSimV2: %s", testCase, err.Error())
	}

	rSim := sim.(*ReorgSimV2)
	filter := func(base uint64) ([]types.Log, error) {
		return rSim.FilterLogs(nil, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(base) - 2),
			ToBlock:   big.NewInt(int64(base) + 2),
		})
	}

	for number := testConf.Param.StartBlock; number < testConf.Param.ExitBlock; number++ {
		logs, err := filter(number)
		if err != nil {
			return errors.Wrapf(err, "[%d] error from FilterLogs", testCase)
		}
		rLogs, err := filter(number)
		if err != nil {
			return errors.Wrapf(err, "[%d] error from FilterLogs", testCase)
		}

		for i, log := range logs {
			if log.BlockNumber >= testConf.Events[len(testConf.Events)-1].ReorgBlock {
				var rLog *types.Log

				// TODO: How to test if lengths not identical?
				if len(rLogs) == len(logs) {
					rLog = &rLogs[i]
				} else {
					t.Log(testCase, log.BlockNumber, len(logs), len(rLogs))
					for _, _rLog := range rLogs {
						if _rLog.TxHash == log.TxHash && _rLog.TxIndex == log.TxIndex {
							rLog = &_rLog
							break
						}
					}
				}

				// If the reorgedLog is missing
				if rLog == nil {
					t.Fatalf(
						"[%d] reorgedLog missing for log (block %d %s) (tx %s index %d)",
						testCase, log.BlockNumber, log.BlockHash.String(), log.TxHash.String(), log.TxIndex,
					)
				}

				if log.BlockHash == rLog.BlockHash {
					// If hashes are the same, but it's reorged hash, then it's ok
					if log.BlockHash == PRandomHash(log.BlockNumber) {
						continue
					}

					t.Fatalf("[number = %d] log and rLog hashes match on %d-%d", number, log.BlockNumber, rLog.BlockNumber)
				}
			}
		}
	}

	return nil
}
