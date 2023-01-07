package poller

import (
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

var (
	// errHashesDiffer will be returned when polled chain data entities on the same block
	// has different block hashes, e.g. when logs filtered on the same block in the same call,
	// have different block hashes, or when block header block hash differs from recent log block hash.
	// When the emitter sees this error, it will change status.isReorging to true, and
	// the emitter will make the poller re-poll the block range (with extra look back blocks)
	errHashesDiffer = errors.Wrap(superwatcher.ErrChainIsReorging, "blockHashes differ")

	// We cannot change poller.pollLevel on the fly for now - will lead to reorg bug.
	errDowngradeLevel = errors.Wrap(superwatcher.ErrUserError, "pollLevel cannot be downgraded")
)
