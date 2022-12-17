package reorgsim

//
// import (
// 	"context"
// 	"math/big"
// 	"testing"
//
// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/core/types"
// )
//
// func TestFilterLogsV2(t *testing.T) {
// 	param := ParamV1{
// 		BaseParam: BaseParam{
// 			StartBlock:    defaultStartBlock,
// 			BlockProgress: 20,
// 		},
// 		ReorgEvent: ReorgEvent{
// 			ReorgBlock: defaultReorgedAt,
// 			MovedLogs:  nil,
// 		},
// 	}
//
// 	sim := NewReorgSimV2FromLogsFiles(param.BaseParam, []ReorgEvent{param.ReorgEvent}, defaultLogsFiles, 4)
//
// 	ctx := context.Background()
// 	logs, err := sim.FilterLogs(ctx, ethereum.FilterQuery{
// 		FromBlock: big.NewInt(69),
// 		ToBlock:   big.NewInt(70),
// 	})
// 	if err != nil {
// 		t.Errorf("FilterLogs returned error: %s", err.Error())
// 	}
// 	if len(logs) != 0 {
// 		t.Fatalf("expecing 0 logs, got %d", len(logs))
// 	}
//
// 	logs, err = sim.FilterLogs(ctx, ethereum.FilterQuery{
// 		FromBlock: big.NewInt(10000000),
// 		ToBlock:   big.NewInt(16000000),
// 	})
// 	if err != nil {
// 		t.Errorf("FilterLogs returned error: %s", err.Error())
// 	}
// 	if len(logs) == 0 {
// 		t.Fatalf("expecting > 0 logs, got 0 log")
// 	}
// }
//
// func TestFilterLogsReorgV2TestFilterLogsReorgV2(t *testing.T) {
// 	reorgedAt := uint64(15944415)
// 	logsPath := "../../internal/emitter/assets"
// 	logsFiles := []string{
// 		logsPath + "/logs_lp.json",
// 		logsPath + "/logs_poolfactory.json",
// 	}
//
// 	param := ParamV1{
// 		BaseParam: BaseParam{
// 			StartBlock:    reorgedAt,
// 			BlockProgress: 20,
// 		},
// 		ReorgEvent: ReorgEvent{
// 			ReorgBlock: reorgedAt,
// 			MovedLogs:  nil,
// 		},
// 	}
//
// 	rSim := NewReorgSimV2FromLogsFiles(
// 		param.BaseParam,
// 		[]ReorgEvent{param.ReorgEvent},
// 		logsFiles,
// 		2,
// 	).(*ReorgSimV2)
//
// 	block := rSim.Chain()[reorgedAt]
// 	rBlock := rSim.ReorgedChain(0)[reorgedAt]
// 	if !rBlock.toBeForked {
// 		t.Fatal("rBlock.toBeForked = false")
// 	}
//
// 	if block.Hash() == rBlock.Hash() {
// 		t.Fatal("block and rBlock hashes match")
// 	}
//
// 	filter := func() ([]types.Log, error) {
// 		return rSim.FilterLogs(nil, ethereum.FilterQuery{
// 			FromBlock: big.NewInt(int64(reorgedAt) - 2),
// 			ToBlock:   big.NewInt(int64(reorgedAt) + 2),
// 		})
// 	}
//
// 	logs, _ := filter()
// 	filter()
// 	filter()
// 	rLogs, _ := filter()
//
// 	for i, log := range logs {
// 		t.Log("foo", log.BlockNumber)
// 		if log.BlockNumber == reorgedAt {
// 			rLog := rLogs[i]
//
// 			if log.BlockHash == rLog.BlockHash {
// 				t.Fatalf("log and rLog hashes match on %d-%d", log.BlockNumber, rLog.BlockNumber)
// 			}
// 		}
// 	}
// }
