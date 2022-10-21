package superwatcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/domain/usecase/engine"
)

func New[K engine.ItemKey, T engine.ServiceItem[K]](
	conf *config.Config,
	ethClient *ethclient.Client,
	dataGateway datagateway.DataGateway,
	stateDataGateway datagateway.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	// TODO: For prod, should we create chans inside this func instead?
	logChan chan *types.Log,
	blockChan chan *reorg.BlockInfo,
	reorgChan chan *reorg.BlockInfo,
	errChan chan error,
	serviceEngine engine.ServiceEngine[K, T],
	debug bool,
) (
	emitter.WatcherEmitter,
	engine.WatcherEngine,
) {
	emitter := emitter.New(
		conf,
		ethClient,
		dataGateway,
		stateDataGateway,
		addresses,
		topics,
		logChan,
		blockChan,
		reorgChan,
		errChan,
		debug,
	)

	engine := engine.New(
		serviceEngine,
		logChan,
		blockChan,
		reorgChan,
		errChan,
		debug,
	)

	return emitter, engine
}
