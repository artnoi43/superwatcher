package servicetest

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate/mockwatcherstate"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

type testCase struct {
	conf          *config.Config
	client        superwatcher.EthClient
	serviceEngine superwatcher.ServiceEngine
}

func newCase(
	conf *config.Config,
	serviceEngine superwatcher.ServiceEngine,
	logsFullPaths []string,
	start,
	reorgAt,
	exit uint64,
) (
	*testCase,
	reorgsim.ReorgParam, // For logging
) {
	param := reorgsim.ReorgParam{
		StartBlock:    start,
		BlockProgress: 5,
		ReorgedAt:     reorgAt,
		ExitBlock:     exit,
	}

	return &testCase{
		conf:          conf,
		client:        reorgsim.NewReorgSim(param, logsFullPaths),
		serviceEngine: serviceEngine,
	}, param
}

func TestServiceEngine(t *testing.T) {
	conf := &config.Config{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		Chain:           string(enums.ChainEthereum),
		StartBlock:      15944390,
		LookBackBlocks:  10,
		LookBackRetries: 2,
		LoopInterval:    0,
	}

	logsPath := "../../internal/emitter/assets"
	logsPathFiles := []string{
		logsPath + "/logs_lp.json",
		logsPath + "/logs_poolfactory.json",
	}

	reorgedAt := uint64(15944415)
	tc, param := newCase(
		conf,
		&engine{
			reorgedAt:       reorgedAt,
			emitterLookBack: conf.LookBackBlocks,
		},
		logsPathFiles,
		conf.StartBlock,
		reorgedAt,
		reorgedAt+(conf.LookBackBlocks*conf.LookBackRetries),
	)

	if err := testServiceEngine(t, tc, param); err != nil {
		t.Error(err.Error())
	}
}

func testServiceEngine(t *testing.T, tc *testCase, param reorgsim.ReorgParam) error {
	// Use nil addresses and topics
	t.Logf("testCase param: %+v", param)
	emitter, engine := initsuperwatcher.New(
		tc.conf,
		tc.client,
		mockwatcherstate.New(tc.conf.StartBlock),
		nil,
		nil,
		tc.serviceEngine,
		true,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := emitter.Loop(ctx); err != nil {
			cancel()
			emitter.Shutdown()
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
