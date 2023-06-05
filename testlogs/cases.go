package testlogs

import "github.com/soyart/superwatcher/pkg/reorgsim"

type TestConfig struct {
	Param     reorgsim.Param        `json:"param"`
	Events    []reorgsim.ReorgEvent `json:"reorgEvents"`
	FromBlock uint64                `json:"fromBlock"`
	ToBlock   uint64                `json:"toBlock"`
	LogsFiles []string              `json:"logs"`
}

var LogsPath string

// SetLogsPath prepend all TestConfig.LogsFiles string with |pathToTestLogs|.
// It should be called in testing init(), and the string pathToTestLogs
// should be relative to the tester package.
func SetLogsPath(pathToTestLogs string) {
	LogsPath = pathToTestLogs
	for _, tc := range append(TestCasesV1, TestCasesV2...) {
		for j, logFile := range tc.LogsFiles {
			old := logFile
			tc.LogsFiles[j] = pathToTestLogs + old
		}
	}
}
