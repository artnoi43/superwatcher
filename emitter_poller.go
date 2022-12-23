package superwatcher

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EmitterPoller filters event logs from the blockchain and maps []types.Log into *FilterResult with EmitterPoller.Poll.
// The result of EmitterPoller.poll is later used by Emitter to emit to Engine.
// superwatcher users can ignore this type if they have no need to update log addresses and topics on-the-fly,
// as EmitterPoller is already wrapped by Emitter.
type EmitterPoller interface {
	// Poll polls event logs from fromBlock to toBlock, and process the logs into *FilterResult for Emitter
	Poll(ctx context.Context, fromBlock, toBlock uint64) (*FilterResult, error)
	// SetDoReorg changes EmitterPoller's FilterResult mapping logic inside EmitterPoller.Poll.
	// If set to true, EmitterPoller maps reorged logs in FilterResult,
	// if set to false, EmitterPoller will not map reorged logs.
	SetDoReorg(bool)
	// DoReorg returns if EmitterPoller is currently processing chain reorg inside EmitterPoller.Poll
	DoReorg() bool
	// Addresses reads EmitterPoller's current event log addresses
	Addresses() []common.Address
	// Topics reads EmitterPoller's current event log topics
	Topics() [][]common.Hash
	// SetAddresses changes EmitterPoller's event log addresses on-the-fly
	SetAddresses([]common.Address)
	// SetTopics changes EmitterPoller's event log topics on-the-fly
	SetTopics([][]common.Hash)
}

// FilterFunc is what used by EmitterPoller to filter event logs. It mirrors EthClient.FilterLogs method signature.
type FilterFunc func(context.Context, ethereum.FilterQuery) ([]types.Log, error)
