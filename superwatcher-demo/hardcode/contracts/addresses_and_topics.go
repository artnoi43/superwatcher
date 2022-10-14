package contracts

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
	oneinchlimitorder "github.com/artnoi43/superwatcher/superwatcher-demo/hardcode/contracts/1inchlimitorder"
	"github.com/artnoi43/superwatcher/superwatcher-demo/hardcode/contracts/uniswapv3pool"
)

const (
	uniswapv3Pool         = "uniswapv3pool"
	uniswapv3PoolAddr     = "0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168"
	oneInchLimitOrder     = "oneInchLimitOrder"
	oneInchLimitOrderAddr = "0x119c71d3bbac22029622cbaec24854d3d32d2828"
)

var contractABIsMap = map[string]string{
	uniswapv3Pool:     uniswapv3pool.Uniswapv3PoolABI,
	oneInchLimitOrder: oneinchlimitorder.OneInchLimitOrderABI,
}

var contractAddressesMap = map[string]common.Address{
	uniswapv3Pool:     common.HexToAddress(uniswapv3PoolAddr),
	oneInchLimitOrder: common.HexToAddress(oneInchLimitOrderAddr),
}

var contractTopicsMap = map[common.Address][]string{
	contractAddressesMap[uniswapv3PoolAddr]: {"Swap"},
	contractAddressesMap[oneInchLimitOrder]: {"OrderFilled", "OrderCanceled"},
}

func interestingTopics(abiStr string, eventKeys ...string) (abi.ABI, []common.Hash, error) {
	contractABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return abi.ABI{}, nil, errors.Wrap(err, "read ABI failed")
	}

	var topics []common.Hash
	for _, eventKey := range eventKeys {
		event, found := contractABI.Events[eventKey]
		if !found {
			return abi.ABI{}, nil, errors.Wrapf(ErrNoSuchABIEvent, "eventKey %s not found", eventKey)
		}
		topics = append(topics, event.ID)
	}

	return contractABI, topics, nil
}

// Hard-code known topics and addresses
func GetABIAddressesAndTopics() (map[common.Address][]string, map[common.Hash]abi.ABI, []common.Address, [][]common.Hash) {
	badTopic := common.Hash{}

	var addresses []common.Address
	for _, addr := range contractAddressesMap {
		addresses = append(addresses, addr)
	}

	var topics []common.Hash
	mapTopicsABI := make(map[common.Hash]abi.ABI)
	for contractName, abiStr := range contractABIsMap {
		contractAddr := contractAddressesMap[contractName]
		topicKeys := contractTopicsMap[contractAddr]
		contractABI, contractTopics, err := interestingTopics(abiStr, topicKeys...)
		if err != nil {
			logger.Panic("failed to init ABI, topics, and address", zap.String("error", err.Error()))
		}

		topics = append(topics, contractTopics...)
		for _, topic := range contractTopics {
			if topic.Hex() == badTopic.Hex() {
				continue
			}
			mapTopicsABI[topic] = contractABI
		}
	}

	return contractTopicsMap, mapTopicsABI, addresses, [][]common.Hash{topics}
}
