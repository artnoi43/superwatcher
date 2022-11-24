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
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/ensengine"
)

type testCase struct {
	conf          *config.EmitterConfig
	client        superwatcher.EthClient
	serviceEngine superwatcher.ServiceEngine
}

func newCase(
	conf *config.EmitterConfig,
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

func TestServiceEngineENS(t *testing.T) {
	conf := &config.EmitterConfig{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		Chain:         string(enums.ChainEthereum),
		StartBlock:    15984020,
		FilterRange:   10,
		GoBackRetries: 2,
		LoopInterval:  0,
	}

	logsPath := "../assets/ens"
	logsPathFiles := []string{
		logsPath + "/logs_reorg_test.json",
	}

	ensStore := datagateway.NewMockDataGatewayENS()
	ensEngine := ensengine.NewEnsSubEngineSuite(ensStore).Engine

	reorgedAt := uint64(15984033)
	tc, param := newCase(
		conf,
		ensEngine,
		logsPathFiles,
		conf.StartBlock,
		reorgedAt,
		15984100,
	)

	if err := serviceEngineTestTemplate(t, tc, param); err != nil {
		t.Error("error in test template", err.Error())
	}

	results, err := ensStore.GetENSes(nil)
	if err != nil {
		t.Error("error from ensStore", err.Error())
	}

	for _, result := range results {
		t.Log(result.BlockNumber, result.BlockHash)
	}
}

func serviceEngineTestTemplate(t *testing.T, tc *testCase, param reorgsim.ReorgParam) error {
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
