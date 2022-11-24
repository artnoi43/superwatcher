package emitter

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// emitter implements Watcher, and other than Config,
// other fields of this structure are defined as ifaces,
// to facil mock testing.
type emitter struct {
	// These fields are used for filtering event logs
	conf      *config.EmitterConfig
	client    superwatcher.EthClient
	tracker   *blockInfoTracker
	addresses []common.Address
	topics    [][]common.Hash

	// Redis-store for tracking last recorded block
	stateDataGateway watcherstate.StateDataGateway

	// These fields are gateways via which
	// external services interact with emitter
	filterResultChan chan<- *superwatcher.FilterResult
	errChan          chan<- error
	syncChan         <-chan struct{}

	debug    bool
	debugger *debugger.Debugger
}

// NewEmitter initializes contract info from config
func New(
	conf *config.EmitterConfig,
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
		logger.Debug("initializing watcher", zap.Any("addresses", addresses), zap.Any("topics", topics))
	}

	return &emitter{
		conf:             conf,
		client:           client,
		stateDataGateway: stateDataGateway,
		tracker:          newTracker("emitter", debug),
		addresses:        addresses,
		topics:           topics,
		syncChan:         syncChan,
		filterResultChan: filterResultChan,
		errChan:          errChan,
		debug:            debug,
		debugger: &debugger.Debugger{
			Key:         "emitter",
			ShouldDebug: debug,
		},
	}
}

// Loop wraps loopFilterLogs to provide graceful shutdown mechanism for emitter.
// When ctx is camcled else where, Loop calls *emitter.shutdown and returns ctx.Err()
func (e *emitter) Loop(ctx context.Context) error {
	status := new(filterLogStatus)

	for {
		// NOTE: this is not clean, but a workaround to prevent infinity loop
		select {
		case <-ctx.Done():
			e.debugger.Debug("shutting down emitter", zap.Any("emitterStatus", status))
			e.Shutdown()
			return errors.Wrap(ctx.Err(), ErrEmitterShutdown.Error())

		default:
			if err := e.loopFilterLogs(ctx, status); err != nil {
				e.debugger.Debug("loopFilterLogs returned", zap.Any("status", status), zap.Error(err))
				e.emitError(errors.Wrap(err, "error in loopFilterLogs"))
			}
		}
	}
}

func (e *emitter) Shutdown() {
	close(e.filterResultChan)
	close(e.errChan)
}

func (e *emitter) SyncsWithEngine() {
	e.debugger.Debug("waiting for engine sync")

	<-e.syncChan

	e.debugger.Debug("synced with engine")
}
