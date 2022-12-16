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
		fileLogs := readJsonLogs(filename)
		logs = append(logs, fileLogs...)
	}

	return logs
}

// InitMappedLogsFromFiles returns unmarshaled hard-coded logs.
// It is export for use in internal/emitter testing.
func InitMappedLogsFromFiles(filenames ...string) map[uint64][]types.Log {
	hardcodedLogs := []types.Log{}
	for _, filename := range filenames {
		logs := readJsonLogs(filename)
		hardcodedLogs = append(hardcodedLogs, logs...)
	}
	mappedLogs := mapLogsToNumber(hardcodedLogs)

	return mappedLogs
}

func mapLogsToNumber(logs []types.Log) map[uint64][]types.Log {
	m := make(map[uint64][]types.Log)
	for _, log := range logs {
		m[log.BlockNumber] = append(m[log.BlockNumber], log)
	}

	return m
}

func mapLogsToTxHash(logs []types.Log) map[common.Hash][]types.Log {
	m := make(map[common.Hash][]types.Log)
	for _, log := range logs {
		m[log.TxHash] = append(m[log.TxHash], log)
	}

	return m
}

func readJsonLogs(filename string) []types.Log {
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
