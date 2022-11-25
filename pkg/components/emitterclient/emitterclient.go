package emitterclient

import (
	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitterclient"
)

func New(
	conf *config.EmitterConfig,
	syncChan chan<- struct{},
	filterResultChan <-chan *superwatcher.FilterResult,
	errChan <-chan error,
) superwatcher.EmitterClient {
	return emitterclient.New(
		conf,
		syncChan,
		filterResultChan,
		errChan,
		conf.LogLevel,
	)
}
