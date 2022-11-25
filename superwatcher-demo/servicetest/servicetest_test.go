package servicetest

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
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
	startBlock uint64
	reorgBlock uint64
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
	reorgsim.ReorgParam, // For logging
) {
	param := reorgsim.ReorgParam{
		StartBlock:    start,
		BlockProgress: 5,
		ReorgedAt:     reorgAt,
		ExitBlock:     exit,
	}

	return &testComponents{
		conf:          conf,
		client:        reorgsim.NewReorgSim(param, logsFullPaths),
		serviceEngine: serviceEngine,
	}, param
}

func TestServiceEngineENS(t *testing.T) {
	logsPath := "../assets/ens"
	testCases := []testCase{
		{
			startBlock: 15984020,
			reorgBlock: 15984033,
			logsFiles: []string{
				logsPath + "/logs_reorg_test.json",
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("testCase for ENS: %+v", testCase)
		// We'll later use |ensStore| to check for saved results
		ensStore := datagateway.NewMockDataGatewayENS()

		err := testServiceEngineENS(testCase.startBlock, testCase.reorgBlock, testCase.logsFiles, ensStore)
		if err != nil {
			t.Error("error in full servicetest:", err.Error())
		}

		results, err := ensStore.GetENSes(nil)
		if err != nil {
			t.Error("error from ensStore:", err.Error())
		}

		for _, result := range results {
			if result.BlockNumber >= testCase.reorgBlock {
				t.Log("checking block", result.BlockNumber)
				// Since reorged block uses hash from deterministic PRandomHash,
				// we can check for equality this way
				expectedHash := reorgsim.PRandomHash(result.BlockNumber).String()
				t.Logf("%+v", result)

				if result.BlockHash != gslutils.ToLower(expectedHash) {
					t.Fatal("unexpected blockHash")
				}
				if result.ID == "" {
					t.Fatal("empty ENS ID -- should not happen")
				}
				if result.Name == "" {
					t.Fatal("empty ENS Name -- should not happen")
				}
			}
		}

	}
}

func testServiceEngineENS(startBlock, reorgedAt uint64, logsFiles []string, ensStore datagateway.DataGatewayENS) error {
	conf := &config.EmitterConfig{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		Chain:         string(enums.ChainEthereum),
		StartBlock:    startBlock,
		FilterRange:   10,
		GoBackRetries: 2,
		LoopInterval:  0,
	}

	ensEngine := ensengine.NewEnsSubEngineSuite(ensStore).Engine

	components, param := initTestComponents(
		conf,
		ensEngine,
		logsFiles,
		conf.StartBlock,
		reorgedAt,
		reorgedAt+conf.FilterRange*conf.GoBackRetries,
	)

	return serviceEngineTestTemplate(components, param)
}

func serviceEngineTestTemplate(components *testComponents, param reorgsim.ReorgParam) error {
	// Use nil addresses and topics
	emitter, engine := initsuperwatcher.New(
		components.conf,
		components.client,
		mockwatcherstate.New(components.conf.StartBlock),
		nil,
		nil,
		components.serviceEngine,
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
