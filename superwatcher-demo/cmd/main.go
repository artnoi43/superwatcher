package main

import (
	"context"

	"os/signal"
	"syscall"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/lib/enums"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/hardcode"
	engine "github.com/artnoi43/superwatcher/superwatcher-demo/internal"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
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

	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	_, addresses, topics := hardcode.GetABIAddressesAndTopics()

	enginew := engine.NewEngine(
		conf,
		ethClient, addresses,
		topics)
	defer stop()

	enginew.HandleLog(func(a *types.Log) {

	})

	enginew.HandleReorg()

	if err := enginew.Loop(ctx); err != nil {
		logger.Error("main error", zap.String("error", err.Error()))
	}
}
