package reorgsim

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
)

func TestFilterLogs(t *testing.T) {
	param := ReorgParam{
		StartBlock: startBlock,
		ReorgedAt:  reorgedAt,
	}
	sim := NewReorgSim(param, defaultLogs)
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
		t.Fatalf("expecting >0 logs, got 0")
	}
}

func TestExitBlock(t *testing.T) {
	exitBlock := reorgedAt + 100
	t.Log("exit block", exitBlock)

	param := ReorgParam{
		StartBlock:    startBlock,
		BlockProgress: 5,
		ReorgedAt:     reorgedAt,
		ExitBlock:     exitBlock,
	}
	sim := NewReorgSim(param, defaultLogs)

	ctx := context.Background()
	for {
		blockNumber, err := sim.BlockNumber(ctx)
		if err != nil {
			t.Log("err", err.Error())
			if blockNumber < exitBlock {
				t.Fatalf("blockNumber < exit: %d < %d", blockNumber, exitBlock)
			}

			if !errors.Is(err, ErrExitBlockReached) {
				t.Fatalf("invalid error returned by reorgsim.BlockNumber: %s", err.Error())
			}

			break
		}

		if blockNumber > exitBlock {
			t.Fatalf("blockNumber < exit: %d < %d", blockNumber, exitBlock)
		}
	}
}
