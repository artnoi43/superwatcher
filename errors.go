package superwatcher

import "github.com/pkg/errors"

var (
	// Ethereum node fetch error
	ErrFetchError = errors.New("fetch from ethclient failed")

	// Chain is reorging - will not cause a return from emitter
	// e.g. when logs filtered from the same block has different hashes
	ErrChainIsReorging  = errors.New("chain is reorging and data is not usable for now")
	ErrFromBlockReorged = errors.Wrap(ErrChainIsReorging, "fromBlock reorged")

	// Bug from my own part
	ErrSuperwatcherBug = errors.New("superwatcher bug")
	ErrProcessReorg    = errors.Wrap(ErrSuperwatcherBug, "error in emitter reorg detection logic") // Bug in reorg detection logic

	// User violates some rules/policies, e.g. downgrading poller PollLevel
	ErrUserError = errors.New("user error")
)
