package emittertest

import (
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

type TestConfig struct {
	Param     reorgsim.Param        `json:"param"`
	Events    []reorgsim.ReorgEvent `json:"reorgEvents"`
	FromBlock uint64                `json:"fromBlock"`
	ToBlock   uint64                `json:"toBlock"`
	LogsFiles []string              `json:"logs"`
}

var logsPath = "../../test_logs"
