package reorgsim

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
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

func testFilterLogsReorgV2(t *testing.T, caseNumber int, testConf multiReorgConfig) error {
	if len(testConf.Events) == 0 {
		return errors.New("got 0 ReorgEvent")
	}

	sim, err := NewReorgSimV2FromLogsFiles(testConf.Param, testConf.Events, testConf.LogsFiles, 2)
	if err != nil {
		t.Fatalf("[%d] failed to create new ReorgSimV2: %s", caseNumber, err.Error())
	}

	movedTo := make(map[uint64][]common.Hash)
	forked := make([]bool, len(testConf.Events))
	for _, event := range testConf.Events {
		for _, moves := range event.MovedLogs {
			for _, move := range moves {
				for _, moveHash := range move.TxHashes {
					if gslutils.Contains(movedTo[move.NewBlock], moveHash) {
						continue
					}

					movedTo[move.NewBlock] = append(movedTo[move.NewBlock], moveHash)
				}
			}
		}
	}

	rSim := sim.(*ReorgSimV2)

	callReorgSimV2 := func(baseNumber uint64) (uint64, []types.Log, error) {
		blockNumber, err := rSim.BlockNumber(nil) // Trigger forward block progression
		if err != nil {
			if !errors.Is(err, ErrExitBlockReached) {
				t.Fatal("unexpected error from ReorgSimV2.BlockNumber", err.Error())
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
		simBlockNumber, logs, err := callReorgSimV2(number)
		if err != nil {
			t.Fatal("ReorgSimV2.FilterLogs returned non-nil error", err.Error())
		}

		event := testConf.Events[reorgCount]
		if simBlockNumber > event.ReorgBlock {
			forked[reorgCount] = true

			if gslutils.Contains(forked, false) {
				reorgCount++
			}
		}

		mappedLogs := mapLogsToNumber(logs)

		for blockNumber, blockLogs := range mappedLogs {
			for _, moves := range event.MovedLogs {
				for _, move := range moves {
					if move.NewBlock != blockNumber {
						continue
					}

					blockTxHashes := gslutils.Map(blockLogs, func(l types.Log) (common.Hash, bool) {
						return l.TxHash, true
					})

					for _, moveTxHash := range move.TxHashes {
						if gslutils.Contains(blockTxHashes, moveTxHash) {
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
