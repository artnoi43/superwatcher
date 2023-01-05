package emitter

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// emitter is the default implementation for superwatcher.Emitter.
type emitter struct {
	sync.RWMutex

	// These fields are used for filtering event logs
	conf *superwatcher.Config

	// client is an implementation of `superwatcher.EthClient`.
	// client.FilterLogs method will be embed into emitterPoller.
	client superwatcher.EthClient

	// stateDataGateway is used to get lastRecordedBlock to determine the next fromBlock
	stateDataGateway superwatcher.GetStateDataGateway

	// poller.Poll filters logs and returns superwatcher.PollResult for emitter to emit
	poller superwatcher.EmitterPoller

	// These fields are gateways via which external components interact with emitter

	pollResultChan chan<- *superwatcher.PollResult // Channel used to send result to consumer
	errChan        chan<- error                    // Channel used to send emitter/emitterPoller errors
	syncChan       <-chan struct{}                 // Channel used to sync with consumer

	// emitter.debug allows us to check if we should calls debugger when debugging in a large for loop.
	// This should save some CPU time.
	debug    bool
	debugger *debugger.Debugger
}

func New(
	conf *superwatcher.Config,
	client superwatcher.EthClient,
	stateDataGateway superwatcher.GetStateDataGateway,
	poller superwatcher.EmitterPoller,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	pollResultChan chan<- *superwatcher.PollResult,
	errChan chan<- error,
) superwatcher.Emitter {
	return &emitter{
		conf:             conf,
		client:           client,
		stateDataGateway: stateDataGateway,
		poller:           poller,
		syncChan:         syncChan,
		pollResultChan:   pollResultChan,
		errChan:          errChan,
		debug:            conf.LogLevel > 0,
		debugger:         debugger.NewDebugger("emitter", conf.LogLevel),
	}
}

func (e *emitter) Poller() superwatcher.EmitterPoller {
	e.Lock()
	defer e.Unlock()

	return e.poller
}

func (e *emitter) SetPoller(poller superwatcher.EmitterPoller) {
	e.Lock()
	defer e.Unlock()

	e.poller = poller
}

// Loop wraps loopEmit to provide graceful shutdown mechanism for emitter.
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
			if err := e.loopEmit(ctx, status); err != nil {
				e.debugger.Debug(1, "loopEmit returned", zap.Any("status", status), zap.Error(err))
				e.emitError(errors.Wrap(err, "error in loopEmit"))
			}
		}
	}
}

// Shutdowns closes `e.pollResultChan` and `e.errChan`.
func (e *emitter) Shutdown() {
	e.debugger.Debug(1, "shutting down emitter - closing channels")
	close(e.pollResultChan)
	close(e.errChan)
}

// SyncsEngine blocks until a signal is sent to `e.syncChan`.
func (e *emitter) SyncsEngine() {
	e.debugger.Debug(1, "waiting for engine sync")
	<-e.syncChan
	e.debugger.Debug(1, "synced with engine")
}
