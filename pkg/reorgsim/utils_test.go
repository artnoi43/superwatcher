package reorgsim

import "testing"

func TestLogsReorgPath(t *testing.T) {
	foo := testsReorgSim[2]

	_, logsPark, logsDest := LogsReorgPaths(foo.Events)
	t.Log(logsPark)
	t.Log(logsDest)
}
