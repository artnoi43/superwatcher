package emitter

import (
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

// New returns a new, default superwatcher.WatcherEmitter.
func New(
	conf *config.Config,
	client superwatcher.EthClient,
	stateDataGateway watcherstate.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	filterResultChan chan<- *superwatcher.FilterResult,
	errChan chan<- error,
	debug bool,
) superwatcher.WatcherEmitter {
	if debug {
		logger.Debug(
			"initializing watcherEmitter",
			zap.Any("addresses", addresses), zap.Any("topics", topics),
		)
	}

	return emitter.New(
		conf,
		client,
		stateDataGateway,
		addresses,
		topics,
		syncChan,
		filterResultChan,
		errChan,
		debug,
	)
}
