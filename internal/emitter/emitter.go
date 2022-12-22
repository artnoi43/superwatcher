package emitter

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// emitter implements Watcher, and other than Config,
// other fields of this structure are defined as ifaces,
// to facil mock testing.
type emitter struct {
	// These fields are used for filtering event logs
	conf *config.EmitterConfig

	// Ethereum client for filtering logs
	client superwatcher.EthClient

	// Addresses passed to client.FilterLogs
	addresses []common.Address
	// Topics passed to client.FilterLogs
	topics [][]common.Hash
	// Redis-store for tracking last recorded block
	stateDataGateway superwatcher.GetStateDataGateway
	// Tracker tracks known block hashes to detect chain reorgs
	tracker *blockInfoTracker

	// These fields are gateways via which
	// external services interact with emitter

	filterResultChan chan<- *superwatcher.FilterResult
	errChan          chan<- error
	syncChan         <-chan struct{}

	// emitter.debug allows us to check if we should calls debugger
	// when debugging in a large for loop. This should save some CPU time.
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
	debugger := debugger.NewDebugger("emitter", conf.LogLevel)
	debugger.Debug(1, "initializing emitter")

	return &emitter{
		conf:             conf,
		client:           client,
		stateDataGateway: stateDataGateway,
		tracker:          newTracker("emitter", conf.LogLevel),
		addresses:        addresses,
		topics:           topics,
		syncChan:         syncChan,
		filterResultChan: filterResultChan,
		errChan:          errChan,
		debug:            conf.LogLevel > 0,
		debugger:         debugger,
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
