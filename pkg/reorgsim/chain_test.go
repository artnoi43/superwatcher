package reorgsim

import (
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
)

type moveConfig struct {
	logsFiles []string
	event     ReorgEvent
}

func TestReorgMoveLogs(t *testing.T) {
	tests := []moveConfig{
		{
			logsFiles: []string{"../../internal/emitter/assets/logs_lp_5.json"},
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
				if gslutils.Contains(move.TxHashes, log.TxHash) {
					return fmt.Errorf("log %s was not removed", log.TxHash.String())
				}
			}

			moveToBlock := reorgedChain[move.NewBlock]
			var foundMovedLog bool
			for _, log := range moveToBlock.logs {
				if gslutils.Contains(move.TxHashes, log.TxHash) {
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
