package ethclientwrapper

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ethClientWrapper struct {
	client *ethclient.Client
}

func (w *ethClientWrapper) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return w.client.FilterLogs(ctx, query)
}

func (w *ethClientWrapper) BlockNumber(ctx context.Context) (uint64, error) {
	return w.client.BlockNumber(ctx)
}
