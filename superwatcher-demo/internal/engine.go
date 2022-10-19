package engine

import (
	"context"
	"math/big"
	"os/signal"
	"syscall"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/emitter"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/emitter/reorg"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// import (
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/pkg/errors"
// 	"go.uber.org/zap"

// 	"github.com/artnoi43/superwatcher/lib/logger/debug"
// )

// // ServiceEngine[T] defines what service should implement and inject into engine.
// type ServiceEngine[K itemKey, T ServiceItem[K]] interface {
// 	// ServiceStateTracker returns service-specific finite state machine.
// 	ServiceStateTracker() (ServiceFSM[K], error)

// 	// MapLogToItem maps Ethereum event log into service-specific type T.
// 	MapLogToItem(l *types.Log) (T, error)

// 	// ActionOptions can be implemented to define arbitary options to be passed to ItemAction.
// 	// ActionOptions(T, EngineLogState, ServiceItemState) []interface{}

// 	// ItemAction is called every time a new, fresh log is converted into ServiceItem,
// 	// The state returned represents the service state that will be assigned to the ServiceItem.
// 	// ItemAction(T, EngineLogState, ServiceItemState, ...interface{}) (State, error)

// 	// ReorgOption can be implemented to define arbitary options to be passed to HandleReorg.
// 	// ReorgOptions(T, EngineLogState, ServiceItemState) []interface{}

// 	// HandleReorg is called in *engine.handleReorg.
// 	// HandleReorg(T, EngineLogState, ServiceItemState, ...interface{}) (State, error)

// 	// TODO: work this out
// 	HandleEmitterError(error) error
// }

// type engine[K itemKey, T ServiceItem[K]] struct {
// 	client        watcherClient[T]    // Interfaces with emitter
// 	serviceEngine ServiceEngine[K, T] // Injected service code
// 	// engineFSM     EngineFSM           // Engine internal state machine
// 	debug bool
// }

type ethClient interface {
	BlockNumber(context.Context) (uint64, error)
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)

	// Not sure if needed
	// BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type engine struct {
	emitter   emitter.Emitter
	logChan   chan<- *types.Log
	reorgChan chan<- *reorg.BlockInfo
}
type Engine interface {
	GetLog()
}

func NewEngine(conf *config.Config,
	client ethClient,
	dataGateway datagateway.DataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	logChan chan<- *types.Log,
	reorgChan chan<- *reorg.BlockInfo,
	handleReorg func(e ...interface{})) Engine {

	ctx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	emitter := emitter.NewEmitter(
		conf,
		client,
		dataGateway,
		addresses,
		topics,
		logChan,
		reorgChan)

	emitter.Loop(ctx)

	return &engine{
		emitter:   emitter,
		logChan:   logChan,
		reorgChan: reorgChan,
	}
}

func (e *engine) GetLog() {

	//get lo from emitter
	// watcher reorg
	// e.handleReorg()

	log := e.logChan
	reorgItem := e.reorgChan

}

// func (e *engine[K, T]) handleLog() error {
// 	serviceEngine, serviceFSM, engineFSM, err := e.initStuff("handleLog")
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		newLog := e.client.WatcherCurrentLog()
// 		if err := handleLog(newLog, serviceEngine, serviceFSM, engineFSM, e.debug); err != nil {
// 			return errors.Wrap(err, "handleLog failed in engine.HandleLog")
// 		}
// 	}
// }

// func (e *engine[K, T]) handleBlock() error {
// 	serviceEngine, serviceFSM, engineFSM, err := e.initStuff("handleBlock")
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		newBlock := e.client.WatcherCurrentBlock()
// 		for _, log := range newBlock.Logs {
// 			if err := handleLog(
// 				log,
// 				serviceEngine,
// 				serviceFSM,
// 				engineFSM,
// 				e.debug,
// 			); err != nil {
// 				return errors.Wrap(err, "handleLog failed in handleBlock")
// 			}
// 		}
// 	}
// }

// func (e *engine[K, T]) handleReorg() error {
// 	serviceEngine, serviceFSM, engineFSM, err := e.initStuff("handleBlock")
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		reorgedBlock := e.client.WatcherReorg()
// 		for _, reorgedLog := range reorgedBlock.Logs {
// 			if err := handleReorgedLog(
// 				reorgedLog,
// 				serviceEngine,
// 				serviceFSM,
// 				engineFSM,
// 			); err != nil {
// 				return errors.Wrap(err, "handleReorg failed in handleReorgedLog")
// 			}
// 		}
// 	}
// }

// func (e *engine[K, T]) handleError() error {
// 	for {
// 		err := e.client.WatcherError()
// 		if err != nil {
// 			err = e.serviceEngine.HandleEmitterError(err)
// 			if err != nil {
// 				return errors.Wrap(err, "serviceEngine failed to handle error")
// 			}

// 			// Emitter error handled in service without error
// 			continue
// 		}

// 		e.debugMsg("got nil error from emitter - should not happen")
// 	}
// }

// func (e *engine[K, T]) initStuff(method string) (ServiceEngine[K, T], ServiceFSM[K], EngineFSM, error) {
// 	serviceFSM, err := e.serviceEngine.ServiceStateTracker()
// 	if err != nil {
// 		return nil, nil, nil, errors.Wrapf(err, "failed to init serviceFSM for %s", method)
// 	}

// 	return e.serviceEngine, serviceFSM, e.engineFSM, nil
// }

// func (e *engine[K, T]) debugMsg(msg string, fields ...zap.Field) {
// 	debug.DebugMsg(e.debug, msg, fields...)
// }
