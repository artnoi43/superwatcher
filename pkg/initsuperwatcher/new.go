package initsuperwatcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/domain/usecase/emitter"
	"github.com/artnoi43/superwatcher/internal/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"
)

// New returns the default implementation of emitter.WatcherEmitter and emitter.WatcherEngine.
// These two objects are paired, and this is the preferred way of initializting superwatcher.
func New(
	conf *config.Config,
	ethClient *ethclient.Client,
	stateDataGateway watcherstate.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,

	// TODO: For prod, should we create chans inside this func instead?
	filterResultChan chan *superwatcher.FilterResult,
	errChan chan error,
	serviceEngine superwatcher.ServiceEngine,
	debug bool,
) (
	superwatcher.WatcherEmitter,
	superwatcher.WatcherEngine,
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
