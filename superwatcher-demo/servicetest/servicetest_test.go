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

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/ensengine"
)

type testCase struct {
	conf          *config.Config
	param         reorgsim.ReorgParam
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
) *testCase {
	param := reorgsim.ReorgParam{
		StartBlock:    start,
		BlockProgress: 5,
		ReorgedAt:     reorgAt,
		ExitBlock:     exit,
	}

	return &testCase{
		conf:          conf,
		param:         param,
		client:        reorgsim.NewReorgSim(param, logsFullPaths),
		serviceEngine: serviceEngine,
	}
}

func TestServiceEngineENS(t *testing.T) {
	conf := &config.Config{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		Chain:           string(enums.ChainEthereum),
		StartBlock:      15847800,
		LookBackBlocks:  10,
		LookBackRetries: 2,
		LoopInterval:    0,
	}

	logsPath := "../assets/ens"
	ensLogs := []string{
		logsPath + "/logs_multi_names.json",
		logsPath + "/logs_single_name.json",
	}

	tc := newCase(
		conf,
		ensengine.NewEnsSubEngineSuite().Engine,
		ensLogs,
		conf.StartBlock,
		15847894,
		15847950,
	)

	testServiceEngine(t, tc)
}

func testServiceEngine(t *testing.T, tc *testCase) error {
	// Use nil addresses and topics
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
