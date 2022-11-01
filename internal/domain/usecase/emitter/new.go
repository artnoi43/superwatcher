package emitter

import (
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"
)

// NewEmitter initializes contract info from config
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
	logger.Debug("initializing watcher", zap.Any("addresses", addresses), zap.Any("topics", topics))
	return &emitter{
		config:           conf,
		client:           client,
		stateDataGateway: stateDataGateway,
		tracker:          reorg.NewTracker(),
		startBlock:       conf.StartBlock,
		addresses:        addresses,
		topics:           topics,
		syncChan:         syncChan,
		filterResultChan: filterResultChan,
		errChan:          errChan,
		debug:            debug,
	}
}
