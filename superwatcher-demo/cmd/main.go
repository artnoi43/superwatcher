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
	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/domain/usecase/superwatcher"
	"github.com/artnoi43/superwatcher/lib/enums"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/uniswapv3factoryengine"
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

	// Demo sub-services
	demoServices := make(
		map[usecase.UseCase]engine.ServiceEngine[usecase.DemoKey, engine.ServiceItem[usecase.DemoKey]],
	)

	// All addresses to be filtered by emitter
	var watcherAddresses []common.Address

	for contractName, contractAddr := range contractAddresses {
		switch contractName {
		case hardcode.Uniswapv3Factory:
			poolFactoryEngine := uniswapv3factoryengine.NewUniswapV3Engine(
				contractABIs[contractAddr],
				contractsEvents[contractAddr],
			)
			demoServices[usecase.UseCaseUniswapv3Factory] = poolFactoryEngine
		}

		watcherAddresses = append(watcherAddresses, contractAddr)
	}

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
		demoServices[usecase.UseCaseUniswapv3Factory], // Now only test run poolFactory sub-service
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
