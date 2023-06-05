package demotest

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/soyart/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/soyart/superwatcher/testlogs"
)

var (
	testLogsPath = "../../../testlogs"
)

func init() {
	testlogs.SetLogsPath(testLogsPath)
}

// findDeletionFromParks ensure that there's at least a DEL operation from DebugDataGateway on |parks|.
// |parks| are blocks which the log appeared on during reorg events, but were not it final, canonical blocks
func findDeletionFromParks(parks []uint64, db datagateway.DebugDataGateway) error {
	for _, park := range parks {
		var foundDel bool
		for _, writeLog := range db.WriteLogs() {
			method, _, blockNumber, _, err := writeLog.Unmarshal()
			if err != nil {
				return errors.Wrap(err, "bad write logs")
			}

			if method != "DEL" {
				continue
			}

			if blockNumber == park {
				foundDel = true
			}
		}

		if !foundDel {
			return fmt.Errorf("found 0 deletion operations for parking block %d", park)
		}
	}

	return nil
}
