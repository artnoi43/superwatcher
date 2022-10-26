package emitter

import (
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
)

// NewEmitter initializes contract info from config
func New(
	conf *config.Config,
	client ethClient,
	stateDataGateway datagateway.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	filterResultChan chan<- *FilterResult,
	errChan chan<- error,
	debug bool,
) WatcherEmitter {
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
