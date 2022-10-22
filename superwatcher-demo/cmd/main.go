package main

/*
   superwatcher-demo shows how to use/embed superwatcher in service code.

   It is designed to use superwatcher to track the following events:
   1. UniswapV3Factory: 'PoolCreated' event
   2. UniswapV3 Pool: 'Swap" event
   3. 1inch Limit Order: 'OrderFilled' and 'OrderCanceled' events

   All logs from 3 contracts is handled by 1 instance of *engine.watcherEngine,
   where field *engine.watcherEngine.serviceEngine is demoengine.demoEngine.

   *demoengine.demoEngine handles all 3 contracts by wrapping other so-called "sub-engines".
   For example, to handle contract Uniswapv3Factory, demoEngine uses uniswapv3factoryengine.uniswapv3PoolFactoryEngine.
*/

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
	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
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
		"testSuperWatcherClient",
		rdb,
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	blockChan := make(chan *reorg.BlockInfo)
	errChan := make(chan error)
	reorgChan := make(chan *reorg.BlockInfo)

	// Hard-coded values for testing
	contractAddresses, contractABIs, contractsEvents, topics := hardcode.DemoAddressesAndTopics(hardcode.Uniswapv3Factory)

	// Demo sub-engines
	demoUseCases := make(map[common.Address]subengines.UseCase)
	demoServices := make(map[subengines.UseCase]engine.ServiceEngine[subengines.DemoKey, engine.ServiceItem[subengines.DemoKey]])

	// All addresses to be filtered by emitter
	var watcherAddresses []common.Address

	for contractName, contractAddr := range contractAddresses {
		switch contractName {
		case hardcode.Uniswapv3Factory:
			poolFactoryEngine := uniswapv3factoryengine.NewUniswapV3Engine(
				contractABIs[contractAddr],
				contractsEvents[contractAddr],
			)
			demoServices[subengines.UseCaseUniswapv3Factory] = poolFactoryEngine
			demoUseCases[contractAddr] = subengines.UseCaseUniswapv3Factory
		}

		watcherAddresses = append(watcherAddresses, contractAddr)
	}
	poolFactoryFSM, err := demoServices[subengines.UseCaseUniswapv3Factory].ServiceStateTracker()
	if err != nil {
		logger.Panic("error getting poolFactoryFSM from poolFactoryEngine", zap.Error(err))
	}

	// demoEngine only wraps uniswapv3PoolFactoryEngine for now.
	// It will later wraps uniswapv3PoolEngine and oneInchLimitOrderEngine
	// and like wise needs their FSMs too.
	demoEngine := demoengine.New(
		demoUseCases,
		demoServices,
		demoengine.NewDemoFSM(poolFactoryFSM),
	)

	watcherEmitter, watcherEngine := superwatcher.New(
		conf,
		ethClient,
		stateDataGateway,
		watcherAddresses,
		topics,
		nil, // Only use blockChan, fuck logChan
		blockChan,
		reorgChan,
		errChan,
		demoEngine,
		true,
	)

	shutdown := func() {
		ethClient.Close()
		if err := rdb.Close(); err != nil {
			logger.Error(
				"error during graceful shutdown - Redis client not properly closed",
				zap.String("error", err.Error()),
			)
		}
		logger.Info("shutdown called")
	}

	defer stop()
	defer shutdown()

	go func() {
		if err := watcherEmitter.Loop(ctx); err != nil {
			logger.Error("DEMO: emitter error", zap.Error(err))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := watcherEngine.Loop(ctx); err != nil {
			logger.Error("DEMO: engine error", zap.Error(err))
		}
	}()
	wg.Wait()
}
