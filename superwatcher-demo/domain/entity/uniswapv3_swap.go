package entity

import "github.com/ethereum/go-ethereum/common"

type UniswapSwap struct {
	FromToken common.Address
	ToToken   common.Address
	AmountIn  uint64
}
