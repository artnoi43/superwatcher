package emitter

import (
	"github.com/pkg/errors"
)

var (
	ErrEmitterShutdown   = errors.New("emitter was told to shutdown - Loop context done") // emitter received a shutdown signal
	ErrMaxRetriesReached = errors.New("emitter has reached max goBackRetries")            // emitter.conf.GoBackRetries has been reached
	errNoNewBlock        = errors.New("no new block")                                     // No new block after the last recorded block
)
