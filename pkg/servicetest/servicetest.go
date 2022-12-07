package servicetest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
	"github.com/artnoi43/superwatcher/pkg/mockwatcherstate"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

// TestCase will be converted into config.EmitterConfig and reorgsim.Param
// to cerate TestComponents
type TestCase struct {
	StartBlock uint64   `json:"startBlock"`
	ReorgBlock uint64   `json:"reorgBlock"`
	ExitBlock  uint64   `json:"exitBlock"`
	LogsFiles  []string `json:"logFiles"`
}

// TestComponents is used by RunServiceTestComponents to instantiate
// superwatcher.WatcherEmitter and superwatcher.WatcherEngine for RunService
type TestComponents struct {
	conf          *config.EmitterConfig
	client        superwatcher.EthClient
	serviceEngine superwatcher.ServiceEngine
}

func InitTestComponents(
	conf *config.EmitterConfig,
	serviceEngine superwatcher.ServiceEngine,
	logsFullPaths []string,
	start uint64,
	reorgAt uint64,
	exit uint64,
) (
	*TestComponents,
	reorgsim.Param, // For logging
) {
	param := reorgsim.Param{
		StartBlock:    start,
		BlockProgress: 5,
		ReorgedBlock:  reorgAt,
		ExitBlock:     exit,
	}

	return &TestComponents{
		conf:          conf,
		client:        reorgsim.NewReorgSimFromLogsFiles(param, logsFullPaths, conf.LogLevel),
		serviceEngine: serviceEngine,
	}, param
}

// RunServiceTestComponents runs the entire service using |components| and |param|.
// It does so by setting up superwatcher.WatcherEmitter and superwatcher.WatcherEngine
// and pass these objects to RunService.
// StateDataGateway is created within this function and will be returned to caller
func RunServiceTestComponents(components *TestComponents, param reorgsim.Param) (superwatcher.StateDataGateway, error) {
	// Use nil addresses and topics
	fakeRedis := mockwatcherstate.New(components.conf.StartBlock)
	emitter, engine := initsuperwatcher.New(
		components.conf,
		components.client,
		fakeRedis,
		fakeRedis,
		nil,
		nil,
		components.serviceEngine,
	)

	return fakeRedis, RunService(emitter, engine)
}

// RunService executes the most basic emitter and engine logic, and returns an error from these components.
func RunService(emitter superwatcher.WatcherEmitter, engine superwatcher.WatcherEngine) error {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	var retErr error
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := emitter.Loop(ctx); err != nil {
			if errors.Is(err, reorgsim.ErrExitBlockReached) {
				cancel()
				emitter.Shutdown()
				return
			}

			retErr = err
		}
	}()

	if err := engine.Loop(ctx); err != nil {
		if errors.Is(err, reorgsim.ErrExitBlockReached) {
			return nil
		}

		retErr = err
	}

	wg.Wait()

	return retErr
}

type DebugEngine struct {
	ReorgedAt          uint64
	EmitterFilterRange uint64

	// Users can inject these extra debugging functions
	HandleFuncGoodLog      func(*types.Log)
	HandleFuncReorgedLog   func(*types.Log)
	HandleFuncEmitterError func(error)

	debugger *debugger.Debugger
}

func (e *DebugEngine) HandleGoodLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	e.debugger.Debug(2, fmt.Sprintf("HandleGoodLogs: got %d logs", len(logs)))
	for _, log := range logs {
		e.debugger.Debug(
			1, "good log info",
			zap.Uint64("blockNumber", log.BlockNumber),
			zap.String("blockHash", gslutils.StringerToLowerString(log.BlockHash)),
		)

		// Calling injected func
		if e.HandleFuncGoodLog != nil {
			e.HandleFuncGoodLog(log)
		}
	}
	return nil, nil
}

func (e *DebugEngine) HandleReorgedLogs(logs []*types.Log, artifacts []superwatcher.Artifact) ([]superwatcher.Artifact, error) {
	e.debugger.Debug(
		1, "GOT REORGED LOG IN SERVICETEST",
		zap.Int("len(artifacts)", len(artifacts)),
		zap.Any("artifacts", artifacts),
	)

	for _, log := range logs {
		e.debugger.Debug(
			1, "reorged log info",
			zap.Uint64("blockNumber", log.BlockNumber),
			zap.String("blockHash", gslutils.StringerToLowerString(log.BlockHash)),
		)

		if e.HandleFuncReorgedLog != nil {
			e.HandleFuncReorgedLog(log)
		}

		// TODO: Polish test checks
		if log.BlockNumber != e.ReorgedAt {
			if log.BlockNumber > e.ReorgedAt+e.EmitterFilterRange {
				return nil, fmt.Errorf(
					"reorgedAt is different from logs passed to HandleReorgedLogs: expecting %d, got %d",
					e.ReorgedAt, log.BlockNumber,
				)
			}
		}
	}

	return nil, nil
}

func (e *DebugEngine) HandleEmitterError(err error) error {
	e.debugger.Debug(1, "got error", zap.Error(err))

	if e.HandleFuncEmitterError != nil {
		e.HandleFuncEmitterError(err)
	}

	return err
}
