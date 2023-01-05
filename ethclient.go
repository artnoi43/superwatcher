package superwatcher

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const MethodGetBlockByNumber = "eth_getBlockByNumber"

// EthClient defines all Ethereum client methods used in superwatcher.
// HeaderByNumber returns BlockHeader because if it uses the actual *types.Header
// then the mock client in `reorgsim` would have to mock types.Header too,
// which is an overkill for now.
type EthClient interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(context.Context, *big.Int) (BlockHeader, error)
	rpcEthClient
}

// ethClient represents methods superwatcher expects to use from a real *ethclient.Client.
type ethClient interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
}

// rpcEthClient is used by poller to get data from client in batch,
// e.g. when getting blocks in batch
type rpcEthClient interface {
	BatchCallContext(context.Context, []rpc.BatchElem) error
}

// ethClientWrapper wraps *ethclient.Client to implement EthClient
// with its HeaderByNumber method signature
type ethClientWrapper struct {
	client    ethClient
	rpcClient rpcEthClient
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

func (w *ethClientWrapper) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return w.rpcClient.BatchCallContext(ctx, b) // nolint:wrapcheck
}

func NewEthClient(ctx context.Context, url string) EthClient {
	rpcClient, err := rpc.DialContext(ctx, url)
	if err != nil {
		panic("failed to create new rpcClient " + err.Error())
	}
	return &ethClientWrapper{
		rpcClient: rpcClient,
		client:    ethclient.NewClient(rpcClient),
	}
}
