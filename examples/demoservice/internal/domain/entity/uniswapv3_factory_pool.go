package entity

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/soyart/superwatcher/examples/demoservice/internal/subengines"
)

type Uniswapv3PoolCreated struct {
	Address      common.Address `json:"lpAddress"`
	Token0       common.Address `json:"token0"`
	Token1       common.Address `json:"token1"`
	Fee          uint64         `json:"fee"`
	BlockCreated uint64         `json:"blockCreated"`
	BlockHash    common.Hash    `json:"blockHash"`
	TxHash       common.Hash    `json:"txHash"`
}

type Uniswapv3FactoryWatcherKey struct {
	lpAddress   string
	blockNumber uint64
}

func (k Uniswapv3FactoryWatcherKey) ForSubEngine() subengines.SubEngineEnum {
	return subengines.SubEngineUniswapv3Factory
}

func (k Uniswapv3FactoryWatcherKey) BlockNumber() uint64 {
	return k.blockNumber
}

func (p *Uniswapv3PoolCreated) ItemKey(opts ...interface{}) any {
	return Uniswapv3FactoryWatcherKey{
		lpAddress:   p.Address.String(),
		blockNumber: p.BlockCreated,
	}
}

func (p *Uniswapv3PoolCreated) DebugString() string {
	return fmt.Sprintf("addr: %s, t0: %s, t1: %s", p.Address.String(), p.Token0.String(), p.Token1.String())
}
