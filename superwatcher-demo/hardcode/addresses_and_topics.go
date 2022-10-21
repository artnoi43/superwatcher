package hardcode

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts"
	oneinchlimitorder "github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/1inchlimitorder"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/uniswapv3factory"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/uniswapv3pool"
)

const (
	uniswapv3Pool         = "uniswapv3pool"
	uniswapv3PoolAddr     = "0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168"
	uniswapv3Factory      = "uniswapv3factory"
	uniswapv3FactoryAddr  = "0x1f98431c8ad98523631ae4a59f267346ea31f984"
	oneInchLimitOrder     = "oneInchLimitOrder"
	oneInchLimitOrderAddr = "0x119c71d3bbac22029622cbaec24854d3d32d2828"
)

var contractABIsMap = map[string]string{
	uniswapv3Pool:     uniswapv3pool.Uniswapv3PoolABI,
	uniswapv3Factory:  uniswapv3factory.Uniswapv3FactoryABI,
	oneInchLimitOrder: oneinchlimitorder.OneInchLimitOrderABI,
}

var contractAddressesMap = map[string]common.Address{
	uniswapv3Pool:     common.HexToAddress(uniswapv3PoolAddr),
	uniswapv3Factory:  common.HexToAddress(uniswapv3FactoryAddr),
	oneInchLimitOrder: common.HexToAddress(oneInchLimitOrderAddr),
}

var contractTopicsMap = map[common.Address][]string{
	contractAddressesMap[uniswapv3Pool]:     {"Swap"},
	contractAddressesMap[uniswapv3Factory]:  {"PoolCreated"},
	contractAddressesMap[oneInchLimitOrder]: {"OrderFilled", "OrderCanceled"},
}

// DemoAddressesAndTopics returns contract information for all demo contracts.
func DemoAddressesAndTopics() (
	map[common.Address]abi.ABI, // Map contract (addr) to ABI
	map[common.Address][]abi.Event, // Map contract (addr) to interesting events
	[]common.Address, // All interesting contract addresses
	[][]common.Hash, // All interesting event log topics
) {
	var addresses []common.Address
	for _, addr := range contractAddressesMap {
		addresses = append(addresses, addr)
	}

	abiMap := make(map[common.Address]abi.ABI)
	interestingEventsMap := make(map[common.Address][]abi.Event)
	var topics []common.Hash
	for contractName, abiStr := range contractABIsMap {

		contractABI, err := abi.JSON(strings.NewReader(abiStr))
		if err != nil {
			logger.Panic("failed to init contractABI", zap.Error(err))
		}

		contractAddr := contractAddressesMap[contractName]
		abiMap[contractAddr] = contractABI

		topicKeys := contractTopicsMap[contractAddr]

		_, interestingEvents, err := contracts.ContractInfo(contractABI, topicKeys...)
		if err != nil {
			logger.Panic("failed to init ABI, topics, and address", zap.String("error", err.Error()))
		}
		interestingEventsMap[contractAddr] = interestingEvents

		// Collect all topics from all contracts
		for _, event := range interestingEvents {
			topics = append(topics, event.ID)
		}
	}

	return abiMap, interestingEventsMap, addresses, [][]common.Hash{topics}
}
