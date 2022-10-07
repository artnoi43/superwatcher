package hardcode

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	_1inchlimitorder "github.com/artnoi43/superwatcher/lib/contracts/1inchlimitorder"
	"github.com/artnoi43/superwatcher/lib/contracts/uniswapv3pool"
)

func interestingTopics(abiStr string, eventKeys ...string) ([][]common.Hash, error) {
	contractABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, errors.Wrap(err, "read ABI failed")
	}

	var topics [][]common.Hash
	for _, eventKey := range eventKeys {
		topic := contractABI.Events[eventKey].ID
		topics = append(topics, []common.Hash{topic})
	}

	return topics, nil
}

// Hard-code known topics and addresses
func AddressesAndTopics() ([]common.Address, [][]common.Hash) {
	addresses := []common.Address{
		common.HexToAddress("0x119c71d3bbac22029622cbaec24854d3d32d2828"), // 1inch Limit Order
		common.HexToAddress("0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168"), // Uniswap v3 DAI/USDC
	}
	oneInchTopics, _ := interestingTopics(_1inchlimitorder.OneInchLimitOrderABI, "OrderCanceled", "OrderFilled")
	uniswapV3Topics, _ := interestingTopics(uniswapv3pool.Uniswapv3PoolABI, "Swap")
	topics := append(oneInchTopics, uniswapV3Topics...)

	return addresses, topics
}
