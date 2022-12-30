package reorgsim

import "testing"

func TestLogsReorgPath(t *testing.T) {
	events := testsReorgSim[2].Events

	var allFroms []uint64
	seenFroms := make(map[uint64]bool)
	for _, event := range events {
		for from := range event.MovedLogs {
			allFroms = append(allFroms, from)
			seenFroms[from] = false
		}
	}

	_, logsParks, logsDest := LogsReorgPaths(events)

	for _, parks := range logsParks {
		for _, park := range parks {
			seenFroms[park] = true
		}
	}

	for from := range seenFroms {
		seen, ok := seenFroms[from]
		if !ok {
			t.Fatalf("bad test code - missing key %d in map", from)
		}

		if !seen {
			t.Fatal("logsParks do not include allFroms")
			t.Log(logsParks)
			t.Log(logsDest)
		}
	}
}
