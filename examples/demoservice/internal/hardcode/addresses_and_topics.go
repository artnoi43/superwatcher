package hardcode

import (
	"github.com/artnoi43/w3utils"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/lib/contracts/ens/enscontroller"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/lib/contracts/ens/ensregistrar"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/lib/contracts/oneinchlimitorder"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/lib/contracts/uniswapv3factory"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/lib/contracts/uniswapv3pool"
)

// These are the hard-coded keys
const (
	Uniswapv3Pool         = "uniswapv3pool"
	Uniswapv3PoolAddr     = "0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168"
	Uniswapv3Factory      = "uniswapv3factory"
	Uniswapv3FactoryAddr  = "0x1f98431c8ad98523631ae4a59f267346ea31f984"
	OneInchLimitOrder     = "oneInchLimitOrder"
	OneInchLimitOrderAddr = "0x119c71d3bbac22029622cbaec24854d3d32d2828"
	ENSRegistrar          = "ens"
	ENSRegistrarAddr      = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
	ENSController         = "enscontroller"
	ENSControllerAddr     = "0x283Af0B28c62C092C9727F1Ee09c02CA627EB7F5"
)

var contractABIsMap = map[string]string{
	Uniswapv3Pool:     uniswapv3pool.Uniswapv3PoolABI,
	Uniswapv3Factory:  uniswapv3factory.Uniswapv3FactoryABI,
	OneInchLimitOrder: oneinchlimitorder.OneInchLimitOrderABI,
	ENSRegistrar:      ensregistrar.EnsRegistrarABI,
	ENSController:     enscontroller.EnsControllerABI,
}

var contractAddressesMap = map[string]string{
	Uniswapv3Pool:     Uniswapv3FactoryAddr,
	Uniswapv3Factory:  Uniswapv3FactoryAddr,
	OneInchLimitOrder: OneInchLimitOrderAddr,
	ENSRegistrar:      ENSRegistrarAddr,
	ENSController:     ENSControllerAddr,
}

var contractTopicsMap = map[string][]string{
	Uniswapv3Pool:     {"Swap"},
	Uniswapv3Factory:  {"PoolCreated"},
	OneInchLimitOrder: {"OrderCreated", "OrderCanceled", "OrderFilled"},
	ENSRegistrar:      {"NameRegistered", "Transfer", "NewOwner", "NewTTL"},
	ENSController:     {"NameRegistered"},
}

// DemoAddressesAndTopics returns contract information for all demo contracts.
func DemoContracts(contractNames ...string) map[string]w3utils.Contract {
	basicContracts := make(map[string]w3utils.Contract)

	for _, contractName := range contractNames {
		contractABI, ok := contractABIsMap[contractName]
		if !ok {
			logger.Debug("hardcoded contract not found", zap.String("contractName", contractName))
			continue
		}

		topics := contractTopicsMap[contractName]
		address := contractAddressesMap[contractName]
		basicContracts[contractName] = w3utils.NewContract(
			contractName,
			contractABI,
			address,
			topics...,
		)
	}

	return basicContracts
}
