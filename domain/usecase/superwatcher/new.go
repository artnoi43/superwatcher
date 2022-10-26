package superwatcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter"
	"github.com/artnoi43/superwatcher/domain/usecase/engine"
)

func New(
	conf *config.Config,
	ethClient *ethclient.Client,
	stateDataGateway datagateway.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,

	// TODO: For prod, should we create chans inside this func instead?
	filterResultChan chan *emitter.FilterResult,
	errChan chan error,
	serviceEngine engine.ServiceEngine,
	debug bool,
) (
	emitter.WatcherEmitter,
	engine.WatcherEngine,
) {

	syncChan := make(chan struct{})

	emitter := emitter.New(
		conf,
		ethClient,
		stateDataGateway,
		addresses,
		topics,
		syncChan,
		filterResultChan,
		errChan,
		debug,
	)

	engine := engine.New(
		conf,
		serviceEngine,
		stateDataGateway,
		syncChan,
		filterResultChan,
		errChan,
		debug,
	)

	return emitter, engine
}
