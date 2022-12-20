package emitter

import (
	"github.com/pkg/errors"
)

var (
	ErrEmitterShutdown   = errors.New("emitter was told to shutdown - Loop context done") // emitter received a shutdown signal
	ErrMaxRetriesReached = errors.New("emitter has reached max goBackRetries")            // emitter.conf.GoBackRetries has been reached

	errNoNewBlock       = errors.New("no new block")                               // No new block after the last recorded block
	errFromBlockReorged = errors.New("fromBlock reorged")                          // fromBlock was reorged
	errFetchError       = errors.New("fetch from ethclient failed")                // Ethereum node fetch error
	errProcessReorg     = errors.New("error in emitter reorg detection logic")     // Bug in reorg detection logic
	errNoHash           = errors.Wrap(errProcessReorg, "missing hash for a block") // Emitter has a missing block hash
)
