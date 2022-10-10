package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/data/watcherstate"
	"github.com/artnoi43/superwatcher/domain/usecase/watcher"
	"github.com/artnoi43/superwatcher/domain/usecase/watcher/reorg"
	"github.com/artnoi43/superwatcher/domain/usecase/watchergateway"
	"github.com/artnoi43/superwatcher/lib/enums"
	"github.com/artnoi43/superwatcher/lib/logger"
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

	logChan := make(chan *types.Log)
	errChan := make(chan error)
	reorgChan := make(chan *reorg.BlockInfo)

	// Hard-coded values for testing
	addresses, topics := hardcode.AddressesAndTopics()
	watcher := watcher.NewWatcherDebug(
		conf,
		ethClient,
		nil, // No DataGateway yet
		stateDataGateway,
		addresses,
		topics,
		logChan,
		errChan,
		reorgChan,
	)

	shutdown := func() {
		ethClient.Close()
		if err := rdb.Close(); err != nil {
			logger.Error(
				"error during graceful shutdown - Redis client not properly closed",
				zap.String("error", err.Error()),
			)
		}
	}

	defer stop()
	defer shutdown()

	go func() {
		if err := watcher.Loop(ctx); err != nil {
			logger.Error("main error", zap.String("error", err.Error()))
		}
	}()

	watcherClient := watchergateway.NewWatcherClientDebug[any](
		logChan,
		errChan,
		reorgChan,
		nil,
	)

	var wg sync.WaitGroup
	wg.Add(3)
	go loopHandleWatcherClientLog(watcherClient, &wg)
	go loopHandleWatcherClientErr(watcherClient, &wg)
	go loopHandleWatcherClientReorg(watcherClient, &wg)
	wg.Wait()
}

func loopHandleWatcherClientLog[T any](wc watchergateway.WatcherClient[T], wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("DEMO: start loopHandleWatcherLog")

	for {
		l := wc.WatcherCurrentLog()
		if l == nil {
			logger.Panic("DEMO: got nil log")
		}

		logger.Info("DEMO: got logs", zap.String("address", l.Address.String()), zap.Any("topics", l.Topics))
	}
}

func loopHandleWatcherClientErr[T any](wc watchergateway.WatcherClient[T], wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("DEMO: start loopHandleWatcherLog")

	for {
		err := wc.WatcherError()
		if err == nil {
			logger.Panic("DEMO: got nil error")
		}

		logger.Info("DEMO: got error", zap.String("error", err.Error()))
	}
}

func loopHandleWatcherClientReorg[T any](wc watchergateway.WatcherClient[T], wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("DEMO: start loopHandleWatcherReorg")
	for {
		reorgedBlock := wc.WatcherReorg()
		if reorgedBlock == nil {
			logger.Panic("DEMO: got nil reorged block")
		}

		logger.Info("DEMO: got reorged blocks", zap.Any("blockNumber", reorgedBlock), zap.String("blockHash", reorgedBlock.Hash.String()))
	}
}

// func monitorChannels(
// 	logChan <-chan *types.Log,
// 	errChan <-chan error,
// 	reorgChan <-chan *struct{},
// ) {
// 	for {
// 		select {
// 		case l := <-logChan:
// 			logger.Info(
// 				"got log",
// 				zap.String("address", l.Address.String()),
// 				zap.String("topics", l.Topics[0].String()),
// 			)
// 		case err := <-errChan:
// 			logger.Error(
// 				"got error",
// 				zap.String("error", err.Error()),
// 			)
// 		case r := <-reorgChan:
// 			if r != nil {
// 				logger.Info("got reorg event")
// 			}
// 		}
// 	}
// }
