package hardcode

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/ens"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/oneinchlimitorder"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/uniswapv3factory"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/uniswapv3pool"
)

// These are the hard-coded keys
const (
	Uniswapv3Pool         = "uniswapv3pool"
	uniswapv3PoolAddr     = "0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168"
	Uniswapv3Factory      = "uniswapv3factory"
	uniswapv3FactoryAddr  = "0x1f98431c8ad98523631ae4a59f267346ea31f984"
	OneInchLimitOrder     = "oneInchLimitOrder"
	oneInchLimitOrderAddr = "0x119c71d3bbac22029622cbaec24854d3d32d2828"
	ENS                   = "ens"
	ENSAddr               = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
)

var contractABIsMap = map[string]string{
	Uniswapv3Pool:     uniswapv3pool.Uniswapv3PoolABI,
	Uniswapv3Factory:  uniswapv3factory.Uniswapv3FactoryABI,
	OneInchLimitOrder: oneinchlimitorder.OneInchLimitOrderABI,
	ENS:               ens.EnsABI,
}

var contractAddressesMap = map[string]common.Address{
	Uniswapv3Pool:     common.HexToAddress(uniswapv3PoolAddr),
	Uniswapv3Factory:  common.HexToAddress(uniswapv3FactoryAddr),
	OneInchLimitOrder: common.HexToAddress(oneInchLimitOrderAddr),
	ENS:               common.HexToAddress(ENSAddr),
}

var contractTopicsMap = map[common.Address][]string{
	contractAddressesMap[Uniswapv3Pool]:     {"Swap"},
	contractAddressesMap[Uniswapv3Factory]:  {"PoolCreated"},
	contractAddressesMap[OneInchLimitOrder]: {"OrderCreated", "OrderCanceled", "OrderFilled"},
	contractAddressesMap[ENS]:               {"NameRegistered", "Transfer", "NewOwner", "NewTTL"},
}

// DemoAddressesAndTopics returns contract information for all demo contracts.
func DemoAddressesAndTopics(contractNames ...string) (
	map[string]common.Address, // Map contract (name) to contract addresses
	map[common.Address]abi.ABI, // Map contract (addr) to ABI
	map[common.Address][]abi.Event, // Map contract (addr) to interesting events
	[][]common.Hash, // All interesting event log topics
) {
	addresses := make(map[string]common.Address)
	for key, addr := range contractAddressesMap {
		if contracts.Contains(contractNames, key) {
			addresses[key] = addr
		}
	}

	abiMap := make(map[common.Address]abi.ABI)
	interestingEventsMap := make(map[common.Address][]abi.Event)
	var topics []common.Hash
	for contractName, abiStr := range contractABIsMap {
		if !contracts.Contains(contractNames, contractName) {
			continue
		}

		contractABI, err := abi.JSON(strings.NewReader(abiStr))
		if err != nil {
			logger.Panic("failed to init contractABI", zap.Error(err))
		}

		contractAddr := contractAddressesMap[contractName]
		abiMap[contractAddr] = contractABI

		topicKeys := contractTopicsMap[contractAddr]

		_, interestingEvents, err := contracts.ContractInfo(contractABI, topicKeys...)
		if err != nil {
			logger.Panic("failed to init ABI, topics, and address", zap.Error(err))
		}
		interestingEventsMap[contractAddr] = interestingEvents

		// Collect all topics from all contracts
		for _, event := range interestingEvents {
			topics = append(topics, event.ID)
		}
	}

	return addresses, abiMap, interestingEventsMap, [][]common.Hash{topics}
}
