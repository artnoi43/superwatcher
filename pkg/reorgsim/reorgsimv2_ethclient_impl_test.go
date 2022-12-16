package reorgsim

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
)

func TestFilterLogsV2(t *testing.T) {
	param := ParamV1{
		BaseParam: BaseParam{
			StartBlock:    defaultStartBlock,
			BlockProgress: 20,
		},
		ReorgEvent: ReorgEvent{
			ReorgedBlock: defaultReorgedAt,
			MovedLogs:    nil,
		},
	}

	sim := NewReorgSimV2FromLogsFiles(param.BaseParam, []ReorgEvent{param.ReorgEvent}, defaultLogsFiles, 4)

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
