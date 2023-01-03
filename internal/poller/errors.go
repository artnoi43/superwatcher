package poller

import (
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

var (
	errNoHash   = errors.Wrap(superwatcher.ErrProcessReorg, "missing hash for a block")   // Emitter has a missing block hash
	errNoHeader = errors.Wrap(superwatcher.ErrProcessReorg, "missing header for a block") // Emitter has a missing block hash
)
