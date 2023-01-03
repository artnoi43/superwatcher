package routerengine

import "github.com/pkg/errors"

var (
	errNoService   = errors.New("failed to get service for subengine")
	errNoSubEngine = errors.New("log address not mapped to subengine") //nolint:unused
)
