package reorgsim

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// InitLogs returns unmarshaled hard-coded logs.
// It is export for use in internal/emitter testing.
func InitLogs(filenames []string) map[uint64][]types.Log {
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
