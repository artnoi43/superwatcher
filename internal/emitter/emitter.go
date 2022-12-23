package emitter

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// poller is a struct used to filter event logs and detecting chain reorg,
// i.e. it produces superwatcher.FilterResult for the emitter.
type poller struct {
	topics    [][]common.Hash
	addresses []common.Address

	filterRange uint64
	filterFunc  func(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	tracker     *blockInfoTracker

	debugger *debugger.Debugger
}

// emitter is the default implementation for superwatcher.WatcherEmitter.
type emitter struct {
	// These fields are used for filtering event logs
	conf *config.EmitterConfig

	// client is an implementation of `superwatcher.EthClient`.
	// client.FilterLogs method will be embed into emitterPoller.
	client superwatcher.EthClient

	// stateDataGateway is used to get lastRecordedBlock to determine the next fromBlock
	stateDataGateway superwatcher.GetStateDataGateway

	// poller filters logs and returns superwatcher.FilterResult for emitter to emit
	poller *poller

	// These fields are gateways via which external components interact with emitter

	filterResultChan chan<- *superwatcher.FilterResult // Channel used to send result to consumer
	errChan          chan<- error                      // Channel used to send emitter/emitterPoller errors
	syncChan         <-chan struct{}                   // Channel used to sync with consumer

	// emitter.debug allows us to check if we should calls debugger when debugging in a large for loop.
	// This should save some CPU time.
	debug    bool
	debugger *debugger.Debugger
}

// New returns a new `superwatcher.WatcherEmitter`
func New(
	conf *config.EmitterConfig,
	client superwatcher.EthClient,
	stateDataGateway superwatcher.GetStateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	filterResultChan chan<- *superwatcher.FilterResult,
	errChan chan<- error,
) superwatcher.WatcherEmitter {
	return &emitter{
		poller: &poller{
			topics:      topics,
			addresses:   addresses,
			filterRange: conf.FilterRange,
			filterFunc:  client.FilterLogs,
			tracker:     newTracker("emitter", conf.LogLevel),
			debugger:    debugger.NewDebugger("emitterPoller", conf.LogLevel),
		},
		conf:             conf,
		client:           client,
		stateDataGateway: stateDataGateway,
		syncChan:         syncChan,
		filterResultChan: filterResultChan,
		errChan:          errChan,
		debug:            conf.LogLevel > 0,
		debugger:         debugger.NewDebugger("emitter", conf.LogLevel),
	}
}

// Loop wraps loopFilterLogs to provide graceful shutdown mechanism for emitter.
// When |ctx| is canceled elsewhere, Loop calls *emitter.shutdown and returns value of ctx.Err()
func (e *emitter) Loop(ctx context.Context) error {
	status := new(emitterStatus)

	for {
		select {
		case <-ctx.Done():
			e.debugger.Debug(1, "shutting down emitter", zap.Any("emitterStatus", status))
			e.Shutdown()
			return errors.Wrap(ctx.Err(), ErrEmitterShutdown.Error())

		default:
			if err := e.loopFilterLogs(ctx, status); err != nil {
				e.debugger.Debug(1, "loopFilterLogs returned", zap.Any("status", status), zap.Error(err))
				e.emitError(errors.Wrap(err, "error in loopFilterLogs"))
			}
		}
	}
}

// Shutdowns closes `e.filterResultChan` and `e.errChan`.
func (e *emitter) Shutdown() {
	close(e.filterResultChan)
	close(e.errChan)
}

// SyncsWithEngine blocks until a signal is sent to `e.syncChan`.
func (e *emitter) SyncsWithEngine() {
	e.debugger.Debug(1, "waiting for engine sync")
	<-e.syncChan
	e.debugger.Debug(1, "synced with engine")
}
