package reorgsim

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
)

// reorgSim implements superwatcher.EthClient[block],
// and will be used in place of the normal Ethereum client
// to test the default emitter implementation.
type reorgSim struct {
	lastRecord   uint64
	lookBack     uint64
	chain        blockChain
	reorgedChain blockChain

	seen map[uint64]int
}

// NewReorgSim returns a new reorgSim with hard-coded good and reorged chains.
func NewReorgSim(lookBack, lastRecord, reorgedAt uint64) superwatcher.EthClient {
	mappedLogs := InitLogs()
	chain, reorgedChain := NewBlockChain(mappedLogs, reorgedAt)

	return &reorgSim{
		lastRecord:   lastRecord,
		lookBack:     lookBack,
		chain:        chain,
		reorgedChain: reorgedChain,
		seen:         make(map[uint64]int),
	}
}

func InitLogs() map[uint64][]types.Log {
	poolFactoryLogs := readJsonLogs("./poolfactory_logs.json")
	lpLogs := readJsonLogs("./lp_logs.json")

	hardcodedLogs := append(poolFactoryLogs, lpLogs...)
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
