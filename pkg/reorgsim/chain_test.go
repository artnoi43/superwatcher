package reorgsim

import (
	"fmt"
	"testing"

	"github.com/artnoi43/gsl"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	logsPath                 = "../../testlogs"
	defaultStartBlock uint64 = 15900000
	defaultReorgedAt  uint64 = 15944444
	defaultLogsFiles         = []string{
		logsPath + "/logs_poolfactory.json",
		logsPath + "/logs_lp.json",
	}
)

type moveConfig struct {
	logsFiles []string
	event     ReorgEvent
}

func TestReorgMoveLogs(t *testing.T) {
	tests := []moveConfig{
		{
			logsFiles: []string{logsPath + "/logs_lp_5.json"},
			event: ReorgEvent{
				ReorgBlock: 15966522,
				MovedLogs: map[uint64][]MoveLogs{
					15966522: {
						{
							NewBlock: 15966527,
							TxHashes: []common.Hash{
								common.HexToHash("0x53f6b4200c700208fe7bb8cb806b0ce962a75e7a31d8a523fbc4affdc22ffc44"),
							},
						},
					},
					15966525: {
						{
							NewBlock: 15966527,
							TxHashes: []common.Hash{
								common.HexToHash("0xa46b7e3264f2c32789c4af8f58cb11293ac9a608fb335e9eb6f0fb08be370211"),
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		if err := testReorgMoveLogs(t, tc); err != nil {
			t.Error(err.Error())
		}
	}
}

func testReorgMoveLogs(t *testing.T, conf moveConfig) error {
	logs := InitMappedLogsFromFiles(conf.logsFiles...)
	_, reorgedChain := NewBlockChainReorgMoveLogs(logs, conf.event)

	movedLogs := make(map[common.Hash]bool)

	for moveFrom, moves := range conf.event.MovedLogs {
		for _, move := range moves {
			moveFromBlock := reorgedChain[moveFrom]
			for _, log := range moveFromBlock.logs {
				if gsl.Contains(move.TxHashes, log.TxHash) {
					return fmt.Errorf("log %s was not removed", log.TxHash.String())
				}
			}

			moveToBlock := reorgedChain[move.NewBlock]
			var foundMovedLog bool
			for _, log := range moveToBlock.logs {
				if gsl.Contains(move.TxHashes, log.TxHash) {
					foundMovedLog = true
					movedLogs[log.TxHash] = true
				}
			}

			if !foundMovedLog {
				return fmt.Errorf("moved logs (from %d) not found in moveToBlock %d", moveFrom, move.NewBlock)
			}

			for _, txHash := range move.TxHashes {
				if !movedLogs[txHash] {
					return fmt.Errorf("moved log %s was not in moveToBlock %d", txHash.String(), move.NewBlock)
				}
			}
		}
	}

	return nil
}

func initDefaultChains(reorgedAt uint64) (BlockChain, BlockChain) {
	return newBlockChainReorgSimple(InitMappedLogsFromFiles(defaultLogsFiles...), reorgedAt)
}

func TestNewBlockChainNg(t *testing.T) {
	oldChain, reorgedChain := initDefaultChains(defaultReorgedAt)
	if err := testBlockChain(t, oldChain, reorgedChain); err != nil {
		t.Fatal(err.Error())
	}
}

// Test if NewBlockChain works properly
func TestNewBlockChain(t *testing.T) {
	oldChain, reorgedChain := initDefaultChains(defaultReorgedAt)
	if err := testBlockChain(t, oldChain, reorgedChain); err != nil {
		t.Fatal(err.Error())
	}
}

func testBlockChain(t *testing.T, oldChain, reorgedChain BlockChain) error {
	for blockNumber, reorgedBlock := range reorgedChain {
		oldBlock := oldChain[blockNumber]

		oldLogs := oldBlock.Logs()
		reorgedLogs := reorgedBlock.Logs()

		if lo, lr := len(oldLogs), len(reorgedLogs); lo != lr {
			return fmt.Errorf("len(logs) not match on block %d", blockNumber)
		}

		if !reorgedBlock.toBeForked {
			continue
		}

		if oldBlock.Hash() == reorgedBlock.Hash() {
			return fmt.Errorf("old and reorg block hashes match on block %d:%s", blockNumber, oldBlock.Hash().String())
		}

		if blockNumber < defaultReorgedAt && reorgedBlock.toBeForked {
			return fmt.Errorf("unreorged block %d from reorgedChain tagged with toBeForked", blockNumber)
		}

		if blockNumber > defaultReorgedAt && !reorgedBlock.toBeForked {
			return fmt.Errorf("reorgedBlock %d not tagged with toBeForked", blockNumber)
		}

		for i, reorgedLog := range reorgedLogs {
			oldLog := oldLogs[i]

			// Uncomment to change txHash when reorg too
			// if reorgedLog.TxHash == oldLog.TxHash {
			// 	t.Fatal("old and reorg log txHash match")
			// }

			if reorgedLog.BlockHash == oldLog.BlockHash {
				return fmt.Errorf("old and reorg log blockHash match %d:%s", blockNumber, reorgedLog.BlockHash.String())
			}
		}
	}

	return nil
}
func prontMapLen[T comparable, U any](m map[T][]U, keyString, lenString string) {
	for k, arr := range m {
		fmt.Println(keyString, k, lenString, len(arr))
	}
}

func prontLogs(logs []types.Log) {
	for _, log := range logs {
		fmt.Println("blockNumber", log.BlockNumber, "blockHash", log.BlockHash.String(), "txHash", log.TxHash.String())
	}
}

func prontBlockChain(chain BlockChain) {
	for _, b := range chain {
		fmt.Println(
			"blockNumber", b.blockNumber,
			"blockhash", b.Hash().String(),
			"len(logs)", len(b.logs),
			"forked", b.toBeForked,
		)
		prontLogs(b.logs)
	}
}
