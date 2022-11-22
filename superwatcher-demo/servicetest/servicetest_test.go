package servicetest

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate/mockwatcherstate"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
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
) *testCase {
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
	}
}

func TestFoo(t *testing.T) {
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

	tc := newCase(
		conf,
		// ensengine.NewEnsSubEngineSuite().Engine,
		&engine{},
		logsPathFiles,
		conf.StartBlock,
		15944415,
		15944555,
	)

	sim := tc.client

	filter := func() ([]types.Log, error) {
		return sim.FilterLogs(nil, ethereum.FilterQuery{
			FromBlock: big.NewInt(15944415),
			ToBlock:   big.NewInt(15944415),
		})
	}

	logs, err := filter()
	if err != nil {
		t.Error(err.Error())
	}

	_logs, err := filter()
	if err != nil {
		t.Error(err.Error())
	}

	for i, l := range logs {
		_l := _logs[i]

		if _l.BlockHash == l.BlockHash {
			t.Fatalf("not reorged")
		}
	}
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

	tc := newCase(
		conf,
		// ensengine.NewEnsSubEngineSuite().Engine,
		&engine{},
		logsPathFiles,
		conf.StartBlock,
		15944415,
		15944555,
	)

	if err := testServiceEngine(t, tc); err != nil {
		t.Error(err.Error())
	}
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
