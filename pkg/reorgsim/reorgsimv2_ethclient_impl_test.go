package reorgsim

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestFilterLogsV2(t *testing.T) {
	param := ParamV1{
		BaseParam: BaseParam{
			StartBlock:    defaultStartBlock,
			BlockProgress: 20,
		},
		ReorgEvent: ReorgEvent{
			ReorgBlock: defaultReorgedAt,
			MovedLogs:  nil,
		},
	}

	sim, err := NewReorgSimV2FromLogsFiles(param.BaseParam, []ReorgEvent{param.ReorgEvent}, defaultLogsFiles, 4)
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
	logsPath := "../../internal/emitter/assets"
	logsFiles := []string{
		logsPath + "/logs_lp.json",
		logsPath + "/logs_poolfactory.json",
	}
	param := BaseParam{
		StartBlock:    15944410,
		ExitBlock:     15944530,
		BlockProgress: 20,
	}
	event := ReorgEvent{
		ReorgBlock: 15944415,
		MovedLogs:  nil,
	}

	sim, err := NewReorgSimV2FromLogsFiles(
		param,
		[]ReorgEvent{event},
		logsFiles,
		2,
	)

	if err != nil {
		t.Fatal("failed to create new ReorgSimV2", err.Error())
	}

	rSim := sim.(*ReorgSimV2)

	filter := func(base uint64) ([]types.Log, error) {
		return rSim.FilterLogs(nil, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(base) - 2),
			ToBlock:   big.NewInt(int64(base) + 2),
		})
	}

	for number := param.StartBlock; number < param.ExitBlock; number++ {
		logs, err := filter(number)
		if err != nil {
			t.Fatal("error from FilterLogs", err.Error())
		}
		rLogs, err := filter(number)
		if err != nil {
			t.Fatal("error from FilterLogs", err.Error())
		}

		for i, log := range logs {
			if log.BlockNumber >= event.ReorgBlock {
				rLog := rLogs[i]

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
}
