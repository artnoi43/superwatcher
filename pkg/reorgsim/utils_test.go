package reorgsim

import "testing"

func TestLogsFinalDst(t *testing.T) {
	foo := testsReorgSim[2]

	_, logsDest := LogsFinalDst(foo.Events)
	t.Log(logsDest)
}
