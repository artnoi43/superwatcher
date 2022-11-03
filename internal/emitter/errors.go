package emitter

import "github.com/pkg/errors"

var (
	errFromBlockReorged = errors.New("fromBlock reorged")
	errFetchError       = errors.New("fetch from node failed")

	errFetchHeader = errors.Wrap(errFetchError, "failed to get block header from ethclient")
	errFetchLogs   = errors.Wrap(errFetchError, "failed to filter log from ethclient")

	ErrEmitterShutdown = errors.New("emitter was told to shutdown - Loop context done")
)

// Usage: wrapBlockNumber(err, errFetchLogs, 69)
func wrapErrBlockNumber(originalErr error, ourErr error, blockNumber uint64) error {
	err := errors.Wrap(originalErr, ourErr.Error())
	return errors.Wrapf(err, "blockNumber %d", blockNumber)
}
