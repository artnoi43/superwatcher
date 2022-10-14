package entity

import "github.com/ethereum/go-ethereum/common"

type UniswapSwap struct {
	Sender  common.Address
	Amount0 uint64
	Amount1 uint64
}
