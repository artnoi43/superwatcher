package ethclientwrapper

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/artnoi43/superwatcher"
)

// ethClientWrapper wraps method HeaderByNumber to implement superwatcher.EthClient.
// This is done so that we can mock the blocks without having to mock *types.Header.
type ethClientWrapper struct {
	client *ethclient.Client
}

func WrapEthClient(client *ethclient.Client) superwatcher.EthClient {
	return &ethClientWrapper{
		client: client,
	}
}

func (w *ethClientWrapper) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return w.client.FilterLogs(ctx, query)
}

func (w *ethClientWrapper) BlockNumber(ctx context.Context) (uint64, error) {
	return w.client.BlockNumber(ctx)
}

func (w *ethClientWrapper) HeaderByNumber(ctx context.Context, number *big.Int) (superwatcher.BlockHeader, error) {
	return w.client.HeaderByNumber(ctx, number)
}
