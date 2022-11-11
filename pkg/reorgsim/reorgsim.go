package reorgsim

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
)

type reorgSim struct {
	lastRecord   uint64
	lookBack     uint64
	chain        blockChain
	reorgedChain blockChain
}

func newReorgSim(lookBack, lastRecord, reorgedAt uint64) *reorgSim {
	mappedLogs := initLogs()
	chain, reorgedChain := newBlockChain(mappedLogs, reorgedAt)

	return &reorgSim{
		lastRecord:   lastRecord,
		lookBack:     lookBack,
		chain:        chain,
		reorgedChain: reorgedChain,
	}
}

func initLogs() map[uint64][]types.Log {
	poolFactoryLogs := readJsonLogs("./poolfactory_logs.json")
	lpLogs := readJsonLogs("./lp_logs.json")

	hardcodedLogs := append(poolFactoryLogs, lpLogs...)
	mappedLogs := mapHardCodedLogs(hardcodedLogs)

	return mappedLogs
}

func mapHardCodedLogs(logs []types.Log) map[uint64][]types.Log {
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
