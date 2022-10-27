package main

import (
	"context"

	"os/signal"
	"syscall"

	engine "github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter/enums"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
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
	var addresses []common.Address
	var topics [][]common.Hash

	client := engine.NewEngine(
		conf,
		ethClient, addresses,
		topics)

	defer stop()

	client.LogHandler(func(log []*types.Log) (engine.Artifact, error) {
		// Do something about log
		return engine.Artifact{}, nil
	})

	client.ReorgHandler(func(Log []*types.Log, artifact []engine.Artifact) (engine.Artifact, error) {
		return engine.Artifact{}, nil
	})

	if err := client.Loop(ctx); err != nil {

	}
}
