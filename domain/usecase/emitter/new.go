package emitter

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	dataGateway datagateway.DataGateway,
	stateDataGateway datagateway.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	logChan chan<- *types.Log,
	blockChan chan<- *reorg.BlockInfo,
	reorgChan chan<- *reorg.BlockInfo,
	errChan chan<- error,
	debug bool,
) WatcherEmitter {
	logger.Debug("initializing watcher", zap.Any("addresses", addresses), zap.Any("topics", topics))
	return &emitter{
		config:           conf,
		client:           client,
		dataGateway:      dataGateway,
		stateDataGateway: stateDataGateway,
		tracker:          reorg.NewTracker(),
		startBlock:       conf.StartBlock,
		addresses:        addresses,
		topics:           topics,
		logChan:          logChan,
		blockChan:        blockChan,
		errChan:          errChan,
		reorgChan:        reorgChan,
		debug:            debug,
	}
}
