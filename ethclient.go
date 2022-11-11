package superwatcher

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EmitterBlockHeader interface {
	Hash() common.Hash
}

// EthClient defines all *ethclient.Client methods used in superwatcher
type EthClient[H EmitterBlockHeader] interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (H, error)
}
