package engine

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/domain/usecase/emitterclient"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

type WatcherEngine interface {
	Loop(context.Context) error
}

// engine is the actual implementation for WatcherEngine.
// engine uses emitterclient.Client to talk/sync to emitter.
type engine struct {
	client             emitterclient.Client         // Interfaces with emitter
	serviceEngine      ServiceEngine                // Injected service code
	stateDataGateway   datagateway.StateDataGateway // Saves lastRecordedBlock to Redis
	engineStateTracker EngineStateTracker           // Engine internal state machine
	debug              bool
}

func newWatcherEngine(
	client emitterclient.Client,
	serviceEngine ServiceEngine,
	statDataGateway datagateway.StateDataGateway,
	debug bool,
) WatcherEngine {
	return &engine{
		client:             client,
		serviceEngine:      serviceEngine,
		stateDataGateway:   statDataGateway,
		engineStateTracker: NewTracker(debug),
		debug:              debug,
	}
}

func (e *engine) Loop(ctx context.Context) error {
	go func() {
		if err := e.run(ctx); err != nil {
			e.debugMsg("*engine.run exited", zap.Error(err))
		}

		defer e.shutdown()
	}()

	return e.handleError()
}

// shutdown is not exported, and the user of the engine should not attempt to call it.
func (e *engine) shutdown() {
	e.client.Shutdown()
}

func (e *engine) run(ctx context.Context) error {
	e.debugMsg("*engine.run started")
	serviceEngine, serviceFSM, engineFSM, err := e.initStuff("handleBlock")
	if err != nil {
		return err
	}

	emitterLookBackBlock := e.client.WatcherConfig().LookBackBlocks

	for {
		result := e.client.WatcherResult()
		// emitter channels are closed if the result is nil
		if result == nil {
			e.debugMsg("*engine.run: got nil filterResult, emitter was probably shutdown")
			return nil
		}

		e.debugMsg("*engine.run: got new filterResult", zap.Int("goodBlocks", len(result.GoodBlocks)), zap.Int("reorgedBlocks", len(result.ReorgedBlocks)))

		// Handle reorged logs
		for _, reorgedBlock := range result.ReorgedBlocks {
			for _, reorgedLog := range reorgedBlock.Logs {
				if err := handleReorgedLog(reorgedLog, serviceEngine, serviceFSM, engineFSM, e.debug); err != nil {
					return errors.Wrap(err, "*engine.run: handleReorgedLog error")
				}
			}
		}

		// Handle fresh, good blocks
		for _, goodBlock := range result.GoodBlocks {
			for _, goodLog := range goodBlock.Logs {
				if err := handleLog(goodLog, serviceEngine, serviceFSM, engineFSM, e.debug); err != nil {
					return errors.Wrap(err, "*engine.run: handleLog error")
				}
			}
		}

		// Save lastRecordedBlock
		if err := e.stateDataGateway.SetLastRecordedBlock(ctx, result.LastGoodBlock); err != nil {
			return errors.Wrap(err, "*engine.run: failed to save lastRecordedBlock to redis")
		}
		e.debugMsg("set lastRecordedBlock", zap.Uint64("blockNumber", result.LastGoodBlock))

		// TODO: Until what block number should we clear?
		e.engineStateTracker.ClearUntil(result.LastGoodBlock - emitterLookBackBlock)
		// Signal emitter to progress
		e.client.WatcherEmitterSync()
	}
}

func (e *engine) handleError() error {
	e.debugMsg("*engine.handleError started")
	for {
		err := e.client.WatcherError()
		if err != nil {
			err = e.serviceEngine.HandleEmitterError(err)
			if err != nil {
				return errors.Wrap(err, "serviceEngine failed to handle error")
			}

			// Emitter error handled in service without error
			continue
		}

		e.debugMsg("got nil error from emitter - should not happen unless errChan was closed")
		break
	}
	return nil
}

func (e *engine) initStuff(method string) (ServiceEngine, ServiceStateTracker, EngineStateTracker, error) {
	serviceFSM, err := e.serviceEngine.ServiceStateTracker()
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to init serviceFSM for %s", method)
	}

	return e.serviceEngine, serviceFSM, e.engineStateTracker, nil
}

func (e *engine) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(e.debug, msg, fields...)
}
