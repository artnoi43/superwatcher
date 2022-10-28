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
	"github.com/artnoi43/superwatcher/data/watcherstate"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter"
	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/domain/usecase/superwatcher"
	"github.com/artnoi43/superwatcher/lib/enums"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/demoengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines/uniswapv3factoryengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/hardcode"
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

	filterResultChan := make(chan *emitter.FilterResult)
	errChan := make(chan error)

	// Hard-coded values for testing
	contractAddresses, contractABIs, contractsEvents, topics := hardcode.DemoAddressesAndTopics(hardcode.Uniswapv3Factory)

	// Demo sub-engines
	demoUseCases := make(map[common.Address]subengines.SubEngineEnum)
	demoServices := make(map[subengines.SubEngineEnum]engine.ServiceEngine)

	// All addresses to be filtered by emitter
	var watcherAddresses []common.Address

	for contractName, contractAddr := range contractAddresses {
		switch contractName {
		case hardcode.Uniswapv3Factory:
			poolFactoryEngine := uniswapv3factoryengine.NewUniswapV3Engine(
				contractABIs[contractAddr],
				contractsEvents[contractAddr],
			)
			demoServices[subengines.SubEngineUniswapv3Factory] = poolFactoryEngine
			demoUseCases[contractAddr] = subengines.SubEngineUniswapv3Factory
		}

		watcherAddresses = append(watcherAddresses, contractAddr)
	}

	// demoEngine only wraps uniswapv3PoolFactoryEngine for now.
	// It will later wraps uniswapv3PoolEngine and oneInchLimitOrderEngine
	// and like wise needs their FSMs too.
	demoEngine := demoengine.New(
		demoUseCases,
		demoServices,
	)

	watcherEmitter, watcherEngine := superwatcher.New(
		conf,
		ethClient,
		stateDataGateway,
		watcherAddresses,
		topics,
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
