package superwatcher

import "github.com/pkg/errors"

var (
	ErrFetchError       = errors.New("fetch from ethclient failed")            // Ethereum node fetch error
	ErrProcessReorg     = errors.New("error in emitter reorg detection logic") // Bug in reorg detection logic
	ErrFromBlockReorged = errors.New("fromBlock reorged")                      // fromBlock was reorged
)
