package emitter

import (
	"github.com/pkg/errors"
)

var (
	errNoNewBlock       = errors.New("no new block")
	errFromBlockReorged = errors.New("fromBlock reorged")
	errFetchError       = errors.New("fetch from ethclient failed")

	errFetchHeader = errors.Wrap(errFetchError, "failed to get block headers")
	errFetchLogs   = errors.Wrap(errFetchError, "failed to filter logs")
	errNoHash      = errors.Wrap(errFetchError, "missing hash for a block")

	ErrMaxRetriesReached = errors.New("emitter has reached max goBackRetries")
	ErrEmitterShutdown   = errors.New("emitter was told to shutdown - Loop context done")
)

// wrapErrBlockNumber wraps err with |sentinelError| (if not nil) and blockNumber
// Usage example: `wrapBlockNumber(69, err, errFetchLogs)`
func wrapErrBlockNumber(blockNumber uint64, err error, sentinelError error) error {
	if sentinelError != nil {
		err = errors.Wrap(err, sentinelError.Error())
	}

	return errors.Wrapf(err, "blockNumber %d", blockNumber)
}
