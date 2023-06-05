package reorgsim

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/soyart/gsl"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func TestFilterLogs(t *testing.T) {
	param := Param{
		StartBlock:    defaultStartBlock,
		BlockProgress: 20,
	}
	event := ReorgEvent{
		ReorgBlock: defaultReorgedAt,
		MovedLogs:  nil,
	}

	sim, err := NewReorgSimFromLogsFiles(param, []ReorgEvent{event}, defaultLogsFiles, "TestFilterLogs", 4)
	if err != nil {
		t.Fatal("error creating ReorgSim", err.Error())
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

func TestFilterLogsReorg(t *testing.T) {
	for i, test := range testsReorgSim {
		testFilterLogsReorg(t, i+1, test)
	}
}

func testFilterLogsReorg(t *testing.T, caseNumber int, testConf multiReorgConfig) error {
	if len(testConf.Events) == 0 {
		return errors.New("got 0 ReorgEvent")
	}

	sim, err := NewReorgSimFromLogsFiles(testConf.Param, testConf.Events, testConf.LogsFiles, testConf.Name, 2)
	if err != nil {
		t.Fatalf("[%d] failed to create new ReorgSim: %s", caseNumber, err.Error())
	}

	movedTo := make(map[uint64][]common.Hash)
	forked := make([]bool, len(testConf.Events))
	for _, event := range testConf.Events {
		for _, moves := range event.MovedLogs {
			for _, move := range moves {
				for _, moveHash := range move.TxHashes {
					if gsl.Contains(movedTo[move.NewBlock], moveHash) {
						continue
					}

					movedTo[move.NewBlock] = append(movedTo[move.NewBlock], moveHash)
				}
			}
		}
	}

	rSim := sim.(*ReorgSim)

	callReorgSim := func(baseNumber uint64) (uint64, []types.Log, error) {
		blockNumber, err := rSim.BlockNumber(nil) // Trigger forward block progression
		if err != nil {
			if !errors.Is(err, ErrExitBlockReached) {
				t.Fatal("unexpected error from ReorgSim.BlockNumber", err.Error())
			}
		}

		logs, err := rSim.FilterLogs(nil, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(baseNumber) - 2),
			ToBlock:   big.NewInt(int64(baseNumber) + 2),
		})

		return blockNumber, logs, err
	}

	var reorgCount int
	for number := testConf.Param.StartBlock; number < testConf.Param.ExitBlock; number++ {
		simBlockNumber, logs, err := callReorgSim(number)
		if err != nil {
			t.Fatal("ReorgSim.FilterLogs returned non-nil error", err.Error())
		}

		event := testConf.Events[reorgCount]
		if simBlockNumber > event.ReorgBlock {
			forked[reorgCount] = true

			if gsl.Contains(forked, false) {
				reorgCount++
			}
		}

		mappedLogs := MapLogsToNumber(logs)

		for blockNumber, blockLogs := range mappedLogs {
			for _, moves := range event.MovedLogs {
				for _, move := range moves {
					if move.NewBlock != blockNumber {
						continue
					}

					blockTxHashes := gsl.Map(blockLogs, func(l types.Log) (common.Hash, bool) {
						return l.TxHash, true
					})

					for _, moveTxHash := range move.TxHashes {
						if gsl.Contains(blockTxHashes, moveTxHash) {
							continue
						}

						t.Log("blockTxHashes", blockTxHashes)
						t.Errorf("moveTxHash (%d) %s missing %s from blockTxHashes", blockNumber, moveTxHash.String(), time.Now().String())
					}
				}
			}
		}
	}

	return nil
}
