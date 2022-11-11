package superwatcher

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockHeader interface {
	Hash() common.Hash
}

// EthClient defines all *ethclient.Client methods used in superwatcher.
// To use a normal *ethclient.Client with superwatcher, wrap it first
// with ethclientwrapper.WrapEthClient.
type EthClient interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)

	// HeaderByNumber returns EmitterBlockHeader so that we can easily mock
	// a client without having to construct *types.Header ourselves.
	// TODO: Perhaps returns the real *types.Header if we have time?
	HeaderByNumber(ctx context.Context, number *big.Int) (BlockHeader, error)
}
