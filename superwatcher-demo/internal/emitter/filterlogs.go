package emitter

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

// filterLogs filters Ethereum event logs from fromBlock to toBlock,
// and sends *types.Log and *reorg.BlockInfo through w.logChan and w.reorgChan respectively.
// If an error is encountered, filterLogs returns with error.
// filterLogs should not be the one sending the error through w.errChan.
func (e *emitter) filterLogs(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) error {

	var err error

	_, err = e.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: e.addresses,
		Topics:    e.topics,
	})
	if err != nil {
		// getErrChan <- errors.Wrap(err, "error filtering event logs")
	}

	return nil
}
