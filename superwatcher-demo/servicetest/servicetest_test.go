package servicetest

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate/mockwatcherstate"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

type testCase struct {
	startBlock uint64
	reorgBlock uint64
	exitBlock  uint64
	logsFiles  []string
}

type testComponents struct {
	conf          *config.EmitterConfig
	client        superwatcher.EthClient
	serviceEngine superwatcher.ServiceEngine
}

func initTestComponents(
	conf *config.EmitterConfig,
	serviceEngine superwatcher.ServiceEngine,
	logsFullPaths []string,
	start,
	reorgAt,
	exit uint64,
) (
	*testComponents,
	reorgsim.Param, // For logging
) {
	param := reorgsim.Param{
		StartBlock:    start,
		BlockProgress: 5,
		ReorgedBlock:  reorgAt,
		ExitBlock:     exit,
	}

	return &testComponents{
		conf:          conf,
		client:        reorgsim.NewReorgSimFromLogsFiles(param, logsFullPaths, conf.LogLevel),
		serviceEngine: serviceEngine,
	}, param
}

func serviceEngineTestTemplate(components *testComponents, param reorgsim.Param) error {
	// Use nil addresses and topics
	emitter, engine := initsuperwatcher.New(
		components.conf,
		components.client,
		mockwatcherstate.New(components.conf.StartBlock),
		nil,
		nil,
		components.serviceEngine,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := emitter.Loop(ctx); err != nil {
			if errors.Is(err, reorgsim.ErrExitBlockReached) {
				cancel()
				emitter.Shutdown()
				return
			}

			panic("unexpected emitter error: " + err.Error())
		}
	}()

	if err := engine.Loop(ctx); err != nil {
		if errors.Is(err, reorgsim.ErrExitBlockReached) {
			return nil
		}

		return err
	}

	wg.Wait()

	return nil
}
