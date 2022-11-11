package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"

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

	stateDataGateway := watcherstate.NewRedisStateDataGateway(
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

	// Hard-coded topic values for testing
	demoContracts := hardcode.DemoContracts(
		hardcode.Uniswapv3Factory,
		hardcode.ENSRegistrar,
		hardcode.ENSController,
	)

	// Init demo service instances and items with demoContracts
	emitterAddresses, emitterTopics, demoRoutes, demoServices := contractsToServices(demoContracts)
	logger.Debug("init: addresses", zap.Any("emitterAddresses", emitterAddresses))
	logger.Debug("init: topics", zap.Any("emitterTopics", emitterTopics))
	logger.Debug("init: demoRoutes", zap.Any("demoRoutes", demoRoutes))
	logger.Debug("init: demoServices", zap.Any("demoServices", demoServices))

	// It will later wraps uniswapv3PoolEngine and oneInchLimitOrderEngine
	// and like wise needs their FSMs too.
	demoEngine := demoengine.New(
		demoRoutes,
		demoServices,
	)

	watcherEmitter, watcherEngine := initsuperwatcher.New[*types.Header](
		conf,
		ethClient,
		stateDataGateway,
		emitterAddresses,
		[][]common.Hash{emitterTopics},
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

func contractsToServices(
	demoContracts map[string]contracts.BasicContract,
) (
	[]common.Address,
	[]common.Hash,
	map[subengines.SubEngineEnum]map[common.Address][]common.Hash, // demoRoutes
	map[subengines.SubEngineEnum]superwatcher.ServiceEngine, // demoServices
) {
	// Demo sub-engines
	demoRoutes := make(map[subengines.SubEngineEnum]map[common.Address][]common.Hash)
	demoServices := make(map[subengines.SubEngineEnum]superwatcher.ServiceEngine)

	// ENS sub-engine has 2 contracts
	// so we can't init the engine in the for loop below
	var ensRegistrar, ensController contracts.BasicContract

	// Topics and addresses to be used by watcher emitter
	var emitterTopics []common.Hash
	var emitterAddresses []common.Address

	for contractName, demoContract := range demoContracts {
		var contractTopics = make([]common.Hash, len(demoContract.ContractEvents))
		var subEngine subengines.SubEngineEnum

		switch contractName {
		case hardcode.Uniswapv3Factory:
			subEngine = subengines.SubEngineUniswapv3Factory
			demoServices[subEngine] = uniswapv3factoryengine.New(demoContract)

		case hardcode.ENSRegistrar, hardcode.ENSController:
			// demoServices for ENS will be created outside of this for loop
			subEngine = subengines.SubEngineENS
			if contractName == hardcode.ENSRegistrar {
				ensRegistrar = demoContract
			} else {
				ensController = demoContract
			}
		}

		for i, event := range demoContract.ContractEvents {
			contractTopics[i] = event.ID
		}

		if demoRoutes[subEngine] == nil {
			demoRoutes[subEngine] = make(map[common.Address][]common.Hash)
		}
		demoRoutes[subEngine][demoContract.Address] = contractTopics
		emitterAddresses = append(emitterAddresses, demoContract.Address)
	}

	// Initialize ensEngine
	ensEngine := ensengine.New(ensRegistrar, ensController)
	demoServices[subengines.SubEngineENS] = ensEngine

	return emitterAddresses, emitterTopics, demoRoutes, demoServices
}
