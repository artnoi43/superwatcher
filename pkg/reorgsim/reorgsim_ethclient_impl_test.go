package reorgsim

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestFilterLogs(t *testing.T) {
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
	sim := NewReorgSimFromLogsFiles(param, defaultLogsFiles, 2)
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

func TestFilterLogsReorg(t *testing.T) {
	reorgedAt := uint64(15944415)
	logsPath := "../../internal/emitter/assets"
	logsFiles := []string{
		logsPath + "/logs_lp.json",
		logsPath + "/logs_poolfactory.json",
	}

	param := ParamV1{
		BaseParam: BaseParam{
			StartBlock:    reorgedAt,
			BlockProgress: 20,
		},
		ReorgEvent: ReorgEvent{
			ReorgBlock: reorgedAt,
			MovedLogs:  nil,
		},
	}

	rSim := NewReorgSimFromLogsFiles(param, logsFiles, 1).(*ReorgSim)

	block := rSim.Chain()[reorgedAt]
	rBlock := rSim.ReorgedChain()[reorgedAt]
	if !rBlock.toBeForked {
		t.Fatal("rBlock.toBeForked = false")
	}

	if block.Hash() == rBlock.Hash() {
		t.Fatal("block and rBlock hashes match")
	}

	filter := func() ([]types.Log, error) {
		return rSim.FilterLogs(nil, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(reorgedAt) - 2),
			ToBlock:   big.NewInt(int64(reorgedAt) + 2),
		})
	}

	logs, _ := filter()
	filter()
	filter()
	rLogs, _ := filter()

	for i, log := range logs {
		if log.BlockNumber == reorgedAt {
			rLog := rLogs[i]

			if log.BlockHash == rLog.BlockHash {
				t.Fatalf("log and rLog hashes match on %d-%d", log.BlockNumber, rLog.BlockNumber)
			}
		}
	}
}

func TestExitBlock(t *testing.T) {
	exitBlock := defaultReorgedAt + 100
	t.Log("exit block", exitBlock)

	param := ParamV1{
		BaseParam: BaseParam{
			StartBlock:    defaultStartBlock,
			BlockProgress: 5,
			ExitBlock:     exitBlock,
		},
		ReorgEvent: ReorgEvent{
			ReorgBlock: defaultReorgedAt,
			MovedLogs:  nil,
		},
	}

	sim := NewReorgSimFromLogsFiles(param, defaultLogsFiles, 1)

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
