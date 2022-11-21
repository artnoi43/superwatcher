package servicetest

import (
	"context"
	"sync"
	"testing"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate/mockwatcherstate"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/ensengine"
)

func TestServiceENS(t *testing.T) {
	conf := &config.Config{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		Chain:           string(enums.ChainEthereum),
		StartBlock:      69,
		LookBackBlocks:  10,
		LookBackRetries: 2,
		LoopInterval:    0,
	}

	logsPath := "../assets/ens"
	param := reorgsim.ReorgParam{}
	fakeEthClient := reorgsim.NewReorgSim(
		param,
		[]string{
			logsPath + "/logs_multi_names.json",
			logsPath + "/logs_single_name.json",
		},
	)

	fakeRedis := mockwatcherstate.New(conf.StartBlock)
	ensEngine := ensengine.NewEnsSubEngineSuite().Engine

	// Use nil addresses and topics
	emitter, engine := initsuperwatcher.New(conf, fakeEthClient, fakeRedis, nil, nil, ensEngine, true)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := emitter.Loop(ctx); err != nil {
			cancel()
		}
	}()

	engine.Loop(ctx)
	wg.Wait()
}
