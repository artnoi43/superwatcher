package poller

import (
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

var errHashesDiffer = errors.Wrap(superwatcher.ErrFetchError, "blockHashes differ")

// errNoHash       = errors.Wrap(superwatcher.ErrProcessReorg, "missing hash for a block")
// errNoHeader     = errors.Wrap(superwatcher.ErrProcessReorg, "missing header for a block")
