package entity

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type Uniswapv3PoolCreated struct {
	Address common.Address
	Token0  common.Address
	Token1  common.Address
}

func (p *Uniswapv3PoolCreated) ItemKey(opts ...interface{}) string {
	return fmt.Sprintf("%s-%s-%s", p.Address.String(), p.Token0.String(), p.Token1.String())
}

func (p *Uniswapv3PoolCreated) DebugString() string {
	return fmt.Sprintf("addr: %s, t0: %s, t1: %s", p.Address.String(), p.Token0.String(), p.Token1.String())
}
