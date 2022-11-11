package reorgsim

import (
	"encoding/json"
	"os"

	"github.com/artnoi43/superwatcher"
	"github.com/ethereum/go-ethereum/core/types"
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

// newReorgSim returns a new reorgSim with hard-coded good and reorged chains.
func newReorgSim(lookBack, lastRecord, reorgedAt uint64) superwatcher.EthClient[block] {
	mappedLogs := initLogs()
	chain, reorgedChain := newBlockChain(mappedLogs, reorgedAt)

	return &reorgSim{
		lastRecord:   lastRecord,
		lookBack:     lookBack,
		chain:        chain,
		reorgedChain: reorgedChain,
		seen:         make(map[uint64]int),
	}
}

func initLogs() map[uint64][]types.Log {
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
