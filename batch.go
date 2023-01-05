package superwatcher

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type BatchCallable interface {
	Marshal() (rpc.BatchElem, error)
	Unmarshal(rpc.BatchElem) error
}

// BatchCall gets []rpc.BatchElem from each batchCall.Marshal() in |batchCalls|.
// It then used |ctx| and the slice []rpc.BatchElem to call client.BatchCallContext.
// After the client call, it interates through |calls| and call batchCall.Unmarshal.
// This means that batchCalls will have their values updated from Unmarshal after BatchCall returns.
func BatchCall(ctx context.Context, client rpcEthClient, batchCalls []BatchCallable) error {
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
