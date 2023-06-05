package batch

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/soyart/superwatcher"
)

// Interface is an intermediate type for passing to CallBatch.
type Interface interface {
	Marshal() (rpc.BatchElem, error)
	Unmarshal(rpc.BatchElem) error
}

const MethodGetBlockByNumber = "eth_getBlockByNumber"

// CallBatch gets []rpc.BatchElem from each batchCall.Marshal() in |batchCalls|.
// It then used |ctx| and the slice []rpc.BatchElem to call client.BatchCallContext.
// After the client call, it interates through |calls| calling Unmarshal.
// This means that batchCalls will have their values updated from Unmarshal after CallBatch returns.
func CallBatch(ctx context.Context, client superwatcher.EthClientRPC, batchCalls []Interface) error {
	batchElems := make([]rpc.BatchElem, len(batchCalls))
	for i, call := range batchCalls {
		elem, err := call.Marshal()
		if err != nil {
			return errors.Wrapf(err, "failed to marshal calls[%d]", i)
		}

		batchElems[i] = elem
	}

	if err := client.BatchCallContext(ctx, batchElems); err != nil {
		return errors.Wrap(err, "BatchCallContext failed")
	}

	for i, call := range batchCalls {
		if err := call.Unmarshal(batchElems[i]); err != nil {
			return errors.Wrapf(err, "failed to marshal calls[%d]", i)
		}
	}

	return nil
}
