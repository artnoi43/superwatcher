package poller

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/batch"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

// headerByNumberBatch is an intermediate type for marshaling and
// unmarshaling rpc.BatchElem for method batch.MethodGetBlockByNumber.
// It implements batch.Interface, so it can be passed to batch.CallBatch.
type headerByNumberBatch struct {
	client string                   // filled by getHeadersByNumbers from reflect.TypeOf(client).String()
	number uint64                   // filled by getHeadersByNumbers
	header superwatcher.BlockHeader // filled by headerByNumberBatch.Unmarshal
}

// Marshal returns BatchElem for calling RPC method `eth_getBlockByNumber` (`batch.MethodGetBlockByNumber`)
// with h.Number, and h.Result will be of type *types.Header.
func (h *headerByNumberBatch) Marshal() (rpc.BatchElem, error) {
	// For mock testing with reorgsim
	if h.client == "*reorgsim.ReorgSim" {
		return rpc.BatchElem{
			Method: batch.MethodGetBlockByNumber,
			Args: []interface{}{
				hexutil.EncodeBig(big.NewInt(int64(h.number))),
				false,
			},
			Result: &reorgsim.Block{},
			Error:  nil,
		}, nil
	}

	return rpc.BatchElem{
		Method: batch.MethodGetBlockByNumber,
		Args: []interface{}{
			hexutil.EncodeBig(big.NewInt(int64(h.number))),
			false,
		},
		Result: &types.Header{},
		Error:  nil,
	}, nil
}

func (h *headerByNumberBatch) Unmarshal(elem rpc.BatchElem) error {
	switch header := elem.Result.(type) {
	case *types.Header:
		// If *types.Header, wrap with BlockHeaderWrapper to implement superwatcher.BlockHeader
		h.header = superwatcher.BlockHeaderWrapper{Header: header}

	case superwatcher.BlockHeader:
		// Otherwise if it's already a BlockHeader, just use it
		// e.g. when the client is reorgsim.ReorgSim, which overwrites elem.Result with *reorgsim.Block.
		h.header = header

	default:
		return fmt.Errorf(
			"unexpected result type for HeaderByNumberBatch: %s",
			reflect.TypeOf(elem.Result).String(),
		)
	}

	return nil
}

func getHeadersByNumbers(
	ctx context.Context,
	client superwatcher.EthClientRPC,
	numbers []uint64,
) (
	map[uint64]superwatcher.BlockHeader,
	error,
) {
	typeOfClient := reflect.TypeOf(client).String()
	elems := make([]batch.Interface, len(numbers))
	for i := range numbers {
		elems[i] = &headerByNumberBatch{
			number: numbers[i],
			client: typeOfClient,
		}
	}

	// Call results will be stored in elems
	if err := batch.CallBatch(ctx, client, elems); err != nil {
		return nil, errors.Wrap(err, "failed to batch get block headers")
	}

	results := make(map[uint64]superwatcher.BlockHeader)
	for _, elem := range elems {
		getHeaderCall := elem.(*headerByNumberBatch)
		results[getHeaderCall.number] = getHeaderCall.header
	}

	return nil, nil
}
