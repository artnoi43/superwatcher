package reorgsim

import (
	"encoding/json"
	"os"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func InitLogsFromFiles(filenames ...string) []types.Log {
	var logs []types.Log
	for _, filename := range filenames {
		fileLogs := readLogsJSON(filename)
		logs = append(logs, fileLogs...)
	}

	return logs
}

// InitMappedLogsFromFiles returns unmarshaled hard-coded logs.
// It is export for use in internal/emitter testing.
func InitMappedLogsFromFiles(filenames ...string) map[uint64][]types.Log {
	hardcodedLogs := []types.Log{}
	for _, filename := range filenames {
		logs := readLogsJSON(filename)
		hardcodedLogs = append(hardcodedLogs, logs...)
	}
	mappedLogs := mapLogsToNumber(hardcodedLogs)

	return mappedLogs
}

// LogsFinalDst iterates through |events| to see which blockNumber is the final destination for a log.
// The return value is a map of log's TX hash to its destination block number, that is, the most
// current ReorgEvent.
func LogsFinalDst(events []ReorgEvent) ([]common.Hash, map[common.Hash]uint64) {
	// Collect MovedLogs info
	type trackMove struct {
		from uint64
		to   uint64
	}

	trackLogs := make(map[int]map[common.Hash]*trackMove)

	for eventIndex, event := range events {
		if trackLogs[eventIndex] == nil {
			trackLogs[eventIndex] = make(map[common.Hash]*trackMove)
		}

		for movedFromBlock, moves := range event.MovedLogs {
			for _, move := range moves {
				for _, txHash := range move.TxHashes {
					trackLogs[eventIndex][txHash] = &trackMove{
						from: movedFromBlock,
						to:   move.NewBlock,
					}
				}
			}
		}
	}

	lenEvents := len(events)
	logsDest := make(map[common.Hash]uint64)
	var logsHashes []common.Hash
	for i := 0; i < lenEvents; i++ {
		moved := trackLogs[i]

		for txHash, move := range moved {
			logsDest[txHash] = move.to

			if gslutils.Contains(logsHashes, txHash) {
				continue
			}

			logsHashes = append(logsHashes, txHash)
		}
	}

	return logsHashes, logsDest
}

func mapLogsToNumber(logs []types.Log) map[uint64][]types.Log {
	m := make(map[uint64][]types.Log)
	for _, log := range logs {
		m[log.BlockNumber] = append(m[log.BlockNumber], log)
	}

	return m
}

func mapLogsToTxHash(logs []types.Log) map[common.Hash][]types.Log { //nolint:unused
	m := make(map[common.Hash][]types.Log)
	for _, log := range logs {
		m[log.TxHash] = append(m[log.TxHash], log)
	}

	return m
}

func readLogsJSON(filename string) []types.Log {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	var logs []types.Log
	if err := json.Unmarshal(b, &logs); err != nil {
		panic(err.Error())
	}

	return logs
}

// appendFilterLogs appends logs from |src| to |dst| with |addresses| and |topics|.
func appendFilterLogs(src, dst *[]types.Log, addresses []common.Address, topics [][]common.Hash) {
	for _, log := range *src {
		if addresses == nil && topics == nil {
			*dst = append(*dst, log)
			continue
		}

		if addresses != nil {
			if topics != nil {
				if gslutils.Contains(topics[0], log.Topics[0]) {
					*dst = append(*dst, log)
					continue
				}
			}

			if gslutils.Contains(addresses, log.Address) {
				*dst = append(*dst, log)
				continue
			}
		}
	}
}

func copyBlockChain(chain BlockChain) BlockChain {
	copied := make(BlockChain)
	for k, v := range chain {
		copied[k] = v
	}

	return copied
}
