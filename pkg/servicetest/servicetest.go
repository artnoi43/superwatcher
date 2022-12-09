package servicetest

import (
	"context"
	"errors"
	"sync"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
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
		client:        reorgsim.NewReorgSimFromLogsFiles(param, logsFullPaths, conf.LogLevel, nil),
		serviceEngine: serviceEngine,
	}, param
}

// RunServiceTestComponents runs the entire service using |components| and |param|.
// It does so by setting up superwatcher.WatcherEmitter and superwatcher.WatcherEngine
// and pass these objects to RunService.
// StateDataGateway is created within this function and will be returned to caller
func RunServiceTestComponents(components *TestComponents) (superwatcher.StateDataGateway, error) {
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
