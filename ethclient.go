package superwatcher

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockHeader interface {
	Hash() common.Hash
}

// EthClient defines all *ethclient.Client methods used in superwatcher.
// HeaderByNumber returns BlockHeader because we don't want to mock the whole types.Header struct in reorgsim
type EthClient interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(context.Context, *big.Int) (BlockHeader, error)
}

type ethClientWrapper struct {
	client *ethclient.Client
}

func (w *ethClientWrapper) BlockNumber(ctx context.Context) (uint64, error) {
	return w.client.BlockNumber(ctx) //nolint:wrapcheck
}

func (w *ethClientWrapper) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return w.client.FilterLogs(ctx, q) //nolint:wrapcheck
}

func (w *ethClientWrapper) HeaderByNumber(ctx context.Context, number *big.Int) (BlockHeader, error) {
	return w.client.HeaderByNumber(ctx, number) //nolint:wrapcheck
}

func WrapEthClient(ethClient *ethclient.Client) EthClient {
	return &ethClientWrapper{client: ethClient}
}
