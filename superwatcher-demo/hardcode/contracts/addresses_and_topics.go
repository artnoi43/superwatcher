package contracts

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	oneinchlimitorder "github.com/artnoi43/superwatcher/superwatcher-demo/hardcode/contracts/1inchlimitorder"
	"github.com/artnoi43/superwatcher/superwatcher-demo/hardcode/contracts/uniswapv3pool"
)

const (
	uniswapv3Pool     = "uniswapv3pool"
	oneInchLimitOrder = "oneInchLimitOrder"
)

var contractABIsMap = map[string]string{
	uniswapv3Pool:     uniswapv3pool.Uniswapv3PoolABI,
	oneInchLimitOrder: oneinchlimitorder.OneInchLimitOrderABI,
}

var contractAddressesMap = map[string]string{
	uniswapv3Pool:     "0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168",
	oneInchLimitOrder: "0x119c71d3bbac22029622cbaec24854d3d32d2828",
}

var contractTopicsMap = map[string][]string{
	uniswapv3Pool:     {"Swap"},
	oneInchLimitOrder: {"OrderFilled", "OrderCanceled"},
}

func interestingTopics(abiStr string, eventKeys ...string) ([]common.Hash, error) {
	contractABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, errors.Wrap(err, "read ABI failed")
	}

	var topics []common.Hash
	for _, eventKey := range eventKeys {
		event, found := contractABI.Events[eventKey]
		if !found {
			return nil, fmt.Errorf("eventKey %s not found", eventKey)
		}
		topics = append(topics, event.ID)
	}

	return topics, nil
}

// Hard-code known topics and addresses
func AddressesAndTopics() ([]common.Address, [][]common.Hash) {
	var addresses []common.Address
	for _, addr := range contractAddressesMap {
		addresses = append(addresses, common.HexToAddress(addr))
	}

	var topics []common.Hash
	for contract, abiStr := range contractABIsMap {
		topicKeys := contractTopicsMap[contract]
		contractTopics, _ := interestingTopics(abiStr, topicKeys...)
		topics = append(topics, contractTopics...)
	}

	return addresses, [][]common.Hash{topics}
}
