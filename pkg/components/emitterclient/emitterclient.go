package emitterclient

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitterclient"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

func New(
	conf *config.Config,
	syncChan chan<- struct{},
	filterResultChan <-chan *superwatcher.FilterResult,
	errChan <-chan error,
	debug bool,
) superwatcher.EmitterClient {
	if debug {
		logger.Debug("initializing emitterClient")
	}

	return emitterclient.New(
		conf,
		syncChan,
		filterResultChan,
		errChan,
		debug,
	)
}
