package poller

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/batch"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

// headerByNumberBatch is an intermediate type for marshaling and
// unmarshaling rpc.BatchElem for method batch.MethodGetBlockByNumber.
// It implements batch.Interface, so it can be passed to batch.CallBatch.
type headerByNumberBatch struct {
	client string                   // filled by mapLogs from reflect.TypeOf(client).String()
	number uint64                   // filled by mapLogs
	header superwatcher.BlockHeader // filled by Unmarshal
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
