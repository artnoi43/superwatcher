package poller

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artnoi43/superwatcher"
)

// headerByNumberBatch is an intermediate type for getting block headers in batch call.
// It implements superwatcher.BatchCallable, so it can be passed to superwatcher.BatchCallContext.
type headerByNumberBatch struct {
	Number uint64
	Header superwatcher.BlockHeader
}

// Marshal returns BatchElem for calling `eth_getBlockByNumber` with h.Number,
// and h.Result will be of type *types.Header.
func (h *headerByNumberBatch) Marshal() (rpc.BatchElem, error) {
	var header types.Header
	return rpc.BatchElem{
		Method: superwatcher.MethodGetBlockByNumber,
		Args: []interface{}{
			hexutil.EncodeBig(big.NewInt(int64(h.Number))),
			false,
		},
		Result: &header,
		Error:  nil,
	}, nil
}

func (h *headerByNumberBatch) Unmarshal(elem rpc.BatchElem) error {
	switch header := elem.Result.(type) {
	case *types.Header:
		// If *types.Header, wrap with BlockHeaderWrapper to implement superwatcher.BlockHeader
		h.Header = superwatcher.BlockHeaderWrapper{Header: header}

	case superwatcher.BlockHeader:
		// Otherwise if it's already a BlockHeader, just use it
		// e.g. when the client is reorgsim.ReorgSim, which overwrites elem.Result with *reorgsim.Block.
		h.Header = header

	default:
		return fmt.Errorf(
			"unexpected result type for HeaderByNumberBatch: %s",
			reflect.TypeOf(elem.Result).String(),
		)
	}

	return nil
}
