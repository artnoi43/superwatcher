package superwatcher

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthClient defines all *ethclient.Client methods used in superwatcher
type EthClient interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}
