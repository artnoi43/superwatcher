package superwatcher

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthClient defines all Ethereum client methods used in superwatcher.
// HeaderByNumber returns BlockHeader because if it uses the actual *types.Header
// then the mock client in `reorgsim` would have to mock types.Header too,
// which is an overkill for now.
type EthClient interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(context.Context, *big.Int) (BlockHeader, error)
}

// ethClient represents methods superwatcher expects to use from a real *ethclient.Client.
type ethClient interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
}

// ethClientWrapper wraps *ethclient.Client to implement EthClient
// with its HeaderByNumber method signature
type ethClientWrapper struct {
	client ethClient
}

func (w *ethClientWrapper) BlockNumber(ctx context.Context) (uint64, error) {
	return w.client.BlockNumber(ctx) //nolint:wrapcheck
}

func (w *ethClientWrapper) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return w.client.FilterLogs(ctx, q) //nolint:wrapcheck
}

func (w *ethClientWrapper) HeaderByNumber(ctx context.Context, number *big.Int) (BlockHeader, error) {
	h, err := w.client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err // nolint:wrapcheck
	}

	return BlockHeaderWrapper{
		Header: h,
	}, nil
}

func WrapEthClient(client ethClient) EthClient {
	return &ethClientWrapper{client: client}
}
