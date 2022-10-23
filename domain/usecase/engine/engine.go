package engine

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

var ErrChanClosed = errors.New("emitterClient channel closed")

type WatcherEngine interface {
	Loop(context.Context) error
}

type engine[K ItemKey, T ServiceItem[K]] struct {
	client        EmitterClient[T]    // Interfaces with emitter
	serviceEngine ServiceEngine[K, T] // Injected service code
	engineFSM     EngineFSM           // Engine internal state machine
	debug         bool
}

func newWatcherEngine[K ItemKey, T ServiceItem[K]](
	client EmitterClient[T],
	serviceEngine ServiceEngine[K, T],
	debug bool,
) WatcherEngine {
	return &engine[K, T]{
		client:        client,
		serviceEngine: serviceEngine,
		engineFSM:     NewEngineFSM(),
		debug:         debug,
	}
}

// Loop is subject to great changes.
// As of this writing, it's not even 50% close to the final version I have in mind.
func (e *engine[K, T]) Loop(ctx context.Context) error {
	go func() {
		if err := e.handleFilterResult(); err != nil {
			logger.Error("handleFilterResult error", zap.Error(err))
		}
	}()

	return e.handleError()
}

func (e *engine[K, T]) handleFilterResult() error {
	e.debugMsg("*engine.handleBlock started")
	serviceEngine, serviceFSM, engineFSM, err := e.initStuff("handleBlock")
	if err != nil {
		return err
	}

	for {
		result := e.client.WatcherResult()
		e.debugMsg("handleFilterResult: got new filterResult", zap.Int("len goodBlocks", len(result.GoodBlocks)), zap.Int("len reorgedBlocks", len(result.ReorgedBlocks)))

		// Handle fresh, good blocks
		for _, goodBlock := range result.GoodBlocks {
			for _, goodLog := range goodBlock.Logs {
				if err := handleLog(goodLog, serviceEngine, serviceFSM, engineFSM, e.debug); err != nil {
					return errors.Wrap(err, "")
				}
			}
		}

		// Handle reorged logs
		for _, reorgedBlock := range result.ReorgedBlocks {
			for _, reorgedLog := range reorgedBlock.Logs {
				if err := handleReorgedLog(reorgedLog, serviceEngine, serviceFSM, engineFSM, e.debug); err != nil {
					return errors.Wrap(err, "")
				}
			}
		}
	}
}

func (e *engine[K, T]) handleError() error {
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

		e.debugMsg("got nil error from emitter - should not happen")
	}
}

func (e *engine[K, T]) initStuff(method string) (ServiceEngine[K, T], ServiceFSM[K], EngineFSM, error) {
	serviceFSM, err := e.serviceEngine.ServiceStateTracker()
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to init serviceFSM for %s", method)
	}

	return e.serviceEngine, serviceFSM, e.engineFSM, nil
}

func (e *engine[K, T]) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(e.debug, msg, fields...)
}
