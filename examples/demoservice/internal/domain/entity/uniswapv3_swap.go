package entity

import "github.com/ethereum/go-ethereum/common"

type Uniswapv3Swap struct {
	Address common.Address
	Token0  common.Address
	Token1  common.Address
}
