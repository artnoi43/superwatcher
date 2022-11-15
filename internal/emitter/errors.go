package emitter

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	errFromBlockReorged = errors.New("fromBlock reorged")
	errFetchError       = errors.New("fetch from ethclient failed")

	errFetchHeader = errors.Wrap(errFetchError, "failed to get block headers")
	errFetchLogs   = errors.Wrap(errFetchError, "failed to filter logs")

	ErrEmitterShutdown = errors.New("emitter was told to shutdown - Loop context done")
)

// wrapErrBlockNumber wraps err with |sentinelError| (if not nil) and blockNumber
// Usage example: `wrapBlockNumber(69, err, errFetchLogs)`
func wrapErrBlockNumber(blockNumber uint64, err error, sentinelError error) error {
	if sentinelError != nil {
		err = errors.Wrap(err, sentinelError.Error())
	}

	return errors.Wrapf(err, "blockNumber %d", blockNumber)
}

func Errorf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}