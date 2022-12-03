package emitter

import (
	"github.com/pkg/errors"
)

var (
	ErrEmitterShutdown   = errors.New("emitter was told to shutdown - Loop context done")
	ErrMaxRetriesReached = errors.New("emitter has reached max goBackRetries")

	errNoNewBlock       = errors.New("no new block")
	errFromBlockReorged = errors.New("fromBlock reorged")
	errFetchError       = errors.New("fetch from ethclient failed")
	errProcessReorg     = errors.New("error in emitter reorg detection logic")

	errFetchLogs = errors.Wrap(errFetchError, "failed to filter logs")
	errNoHash    = errors.Wrap(errProcessReorg, "missing hash for a block")
)
