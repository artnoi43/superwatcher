package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/data/watcherstate"
	"github.com/artnoi43/superwatcher/domain/usecase/watcher"
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

	shutdown := func() {
		ethClient.Close()
		if err := rdb.Close(); err != nil {
			logger.Error(
				"error during graceful shutdown - Redis client not properly closed",
				zap.String("error", err.Error()),
			)
		}
	}

	stateDataGateway := watcherstate.NewWatcherStateRedisClient(
		chain,
		"testSuperWatcherClient",
		rdb,
	)

	logChan := make(chan *types.Log)
	errChan := make(chan error)
	reorgChan := make(chan *struct{})

	go monitorChannels(logChan, errChan, reorgChan)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	defer stop()
	defer shutdown()

	// Hard-coded values for testing
	addresses, topics := hardcode.AddressesAndTopics()
	watcher := watcher.NewWatcher(
		conf,
		ethClient,
		nil,
		stateDataGateway,
		addresses,
		topics,
		logChan,
		errChan,
		reorgChan,
	)

	//	go func() {
	if err := watcher.Loop(ctx); err != nil {
		logger.Error("main error", zap.String("error", err.Error()))
	}
	//	}()

	//	watcherClient := watchergateway.NewWatcherClient[any](
	//		logChan,
	//		errChan,
	//		reorgChan,
	//		nil,
	//	)
}

func monitorChannels(
	logChan <-chan *types.Log,
	errChan <-chan error,
	reorgChan <-chan *struct{},
) {
	for {
		select {
		case l := <-logChan:
			logger.Info(
				"got log",
				zap.String("address", l.Address.String()),
				zap.String("topics", l.Topics[0].String()),
			)
		case err := <-errChan:
			logger.Error(
				"got error",
				zap.String("error", err.Error()),
			)
		case r := <-reorgChan:
			if r != nil {
				logger.Info("got reorg event")
			}
		}
	}
}
