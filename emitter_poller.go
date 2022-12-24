package superwatcher

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// EmitterPoller filters event logs from the blockchain and maps []types.Log into *FilterResult with EmitterPoller.Poll.
// The result of EmitterPoller.poll is later used by Emitter to emit to Engine.
// superwatcher users can ignore this type if they have no need to update log addresses and topics on-the-fly,
// as EmitterPoller is already wrapped by Emitter.
type EmitterPoller interface {
	// Poll polls event logs from fromBlock to toBlock, and process the logs into *FilterResult for Emitter
	Poll(ctx context.Context, fromBlock, toBlock uint64) (*FilterResult, error)

	Controller
}

// FilterFunc is what used by EmitterPoller to filter event logs. It mirrors EthClient.FilterLogs method signature.
type FilterFunc func(context.Context, ethereum.FilterQuery) ([]types.Log, error)
