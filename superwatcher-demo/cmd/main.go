package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/demoengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines/ensengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines/uniswapv3factoryengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/hardcode"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
)

func main() {
	conf, err := config.ConfigYAML("./config/config.yaml")
	if err != nil {
		panic("failed to read YAML config: " + err.Error())
	}

	chain := enums.ChainType(conf.Chain)
	if !chain.IsValid() {
		panic("invalid chain: " + conf.Chain)
	}

	ethClient, err := ethclient.Dial(conf.NodeURL)
	if err != nil {
		panic("new ethclient failed: " + err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: conf.RedisConnAddr,
	})
	if rdb == nil {
		panic("nil redis")
	}

	stateDataGateway := watcherstate.NewWatcherStateRedisClient(
		chain,
		"superwatcher-demo",
		rdb,
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	filterResultChan := make(chan *superwatcher.FilterResult)
	errChan := make(chan error)

	// Demo sub-engines
	demoRoutes := make(map[subengines.SubEngineEnum][]common.Address)
	demoServices := make(map[subengines.SubEngineEnum]superwatcher.ServiceEngine)

	// Hard-coded topic values for testing
	demoContracts := hardcode.DemoContracts(
		hardcode.Uniswapv3Factory,
		hardcode.ENSRegistrar,
		hardcode.ENSController,
	)

	// ENS sub-engine has 2 contracts
	// so we can't init the engine in the for loop below
	var ensRegistrar, ensController contracts.BasicContract

	// Topics and addresses to be used by watcher emitter
	var watcherTopics []common.Hash
	watcherAddresses := make([]common.Address, len(demoContracts))

	for contractName, demoContract := range demoContracts {
		switch contractName {
		case hardcode.Uniswapv3Factory:
			subEngine := subengines.SubEngineUniswapv3Factory
			demoRoutes[subEngine] = []common.Address{demoContract.Address}
			demoServices[subengines.SubEngineUniswapv3Factory] = uniswapv3factoryengine.New(demoContract)

		case hardcode.ENSRegistrar, hardcode.ENSController:
			subEngine := subengines.SubEngineENS
			demoRoutes[subEngine] = append(demoRoutes[subEngine], demoContract.Address)
			if contractName == hardcode.ENSRegistrar {
				ensRegistrar = demoContract
			} else {
				ensController = demoContract
			}
		}

		for _, event := range demoContract.ContractEvents {
			watcherTopics = append(watcherTopics, event.ID)
		}
		watcherAddresses = append(watcherAddresses, demoContract.Address)
	}

	ensEngine := ensengine.New(ensRegistrar, ensController)
	demoServices[subengines.SubEngineENS] = ensEngine

	// It will later wraps uniswapv3PoolEngine and oneInchLimitOrderEngine
	// and like wise needs their FSMs too.
	demoEngine := demoengine.New(
		demoRoutes,
		demoServices,
	)

	watcherEmitter, watcherEngine := initsuperwatcher.New(
		conf,
		ethClient,
		stateDataGateway,
		watcherAddresses,
		[][]common.Hash{watcherTopics},
		filterResultChan, // Only use blockChan, fuck logChan
		errChan,
		demoEngine,
		true,
	)

	// Graceful shutdown
	defer func() {
		// Cancel context to stop both superwatcher emitter and engine
		cancel()

		ethClient.Close()
		if err := rdb.Close(); err != nil {
			logger.Error(
				"error during graceful shutdown - Redis client not properly closed",
				zap.Error(err),
			)
		}

		logger.Info("graceful shutdown successful")
	}()

	go func() {
		if err := watcherEmitter.Loop(ctx); err != nil {
			logger.Error("DEMO: emitter returned an error", zap.Error(err))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := watcherEngine.Loop(ctx); err != nil {
			logger.Error("DEMO: engine returned an error", zap.Error(err))
		}
	}()

	wg.Wait()
}
