package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/artnoi43/gsl/soyutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"

	"github.com/artnoi43/superwatcher/superwatcher-demo/config"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/hardcode"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/routerengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/ensengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/uniswapv3factoryengine"
)

func main() {
	conf, err := soyutils.ReadFileYAMLPointer[config.Config]("./superwatcher-demo/config/config.yaml")
	if err != nil {
		panic("failed to read YAML config: " + err.Error())
	}

	chain := conf.Chain
	if chain == "" {
		panic("empty chain")
	}

	ethClient, err := ethclient.Dial(conf.SuperWatcherConfig.NodeURL)
	if err != nil {
		panic("new ethclient failed: " + err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: conf.SuperWatcherConfig.RedisConnAddr,
	})
	if rdb == nil {
		panic("nil redis")
	}

	stateDataGateway, err := watcherstate.NewRedisStateDataGateway(
		"superwatcher-demo"+":"+chain,
		rdb,
	)
	if err != nil {
		panic("new stateDataGateway failed: " + err.Error())
	}

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
	emitterAddresses, emitterTopics, demoRoutes, demoServices := contractsToServices(demoContracts, rdb, conf.SuperWatcherConfig.LogLevel)
	logger.Debug("init: addresses", zap.Any("emitterAddresses", emitterAddresses))
	logger.Debug("init: topics", zap.Any("emitterTopics", emitterTopics))
	logger.Debug("init: demoRoutes", zap.Any("demoRoutes", demoRoutes))
	logger.Debug("init: demoServices", zap.Any("demoServices", demoServices))

	// It will later wraps uniswapv3PoolEngine and oneInchLimitOrderEngine
	// and like wise needs their FSMs too.
	demoEngine := routerengine.New(
		demoRoutes,
		demoServices,
		conf.SuperWatcherConfig.LogLevel,
	)

	watcherEmitter, watcherEngine := initsuperwatcher.New(
		conf.SuperWatcherConfig,
		// Wraps ethClient to make HeaderByNumber returns superwatcher.EmitterBlockHeader
		ethClient,
		// Both data gateways can be implemented separately by different structs,
		// but here in the demo it's just using default implementation.
		superwatcher.GetStateDataGatewayFunc(stateDataGateway.GetLastRecordedBlock), // stateDataGateway alone is fine too
		superwatcher.SetStateDataGatewayFunc(stateDataGateway.SetLastRecordedBlock), // stateDataGateway alone is fine too
		emitterAddresses,
		[][]common.Hash{emitterTopics},
		demoEngine,
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
	rdb *redis.Client,
	logLevel uint8,
) (
	[]common.Address,
	[]common.Hash,
	map[subengines.SubEngineEnum]map[common.Address][]common.Hash, // demoRoutes
	map[subengines.SubEngineEnum]superwatcher.ServiceEngine, // demoServices
) {
	// Demo sub-engines
	demoRoutes := make(map[subengines.SubEngineEnum]map[common.Address][]common.Hash)
	demoServices := make(map[subengines.SubEngineEnum]superwatcher.ServiceEngine)

	dgwENS := datagateway.NewEnsDataGateway(rdb)
	dgwPoolFactory := datagateway.NewDataGatewayPoolFactory(rdb)

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
			demoServices[subEngine] = uniswapv3factoryengine.New(demoContract, dgwPoolFactory, logLevel)

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
	ensEngine := ensengine.New(ensRegistrar, ensController, dgwENS, logLevel)
	demoServices[subengines.SubEngineENS] = ensEngine

	return emitterAddresses, emitterTopics, demoRoutes, demoServices
}
