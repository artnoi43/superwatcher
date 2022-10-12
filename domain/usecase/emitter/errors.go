package emitter

import "github.com/pkg/errors"

var errFromBlockReorged = errors.New("filterLogs: fromBlock reorged")
