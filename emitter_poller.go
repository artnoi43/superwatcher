package superwatcher

import (
	"context"
)

// EmitterPoller filters event logs from the blockchain and maps []types.Log into *PollResult.
// The result of EmitterPoller.poll is later used by Emitter to emit to Engine.
// superwatcher users can ignore this type if they have no need to update log addresses and topics on-the-fly,
// as EmitterPoller is already wrapped by Emitter.
type EmitterPoller interface {
	// Poll polls event logs from fromBlock to toBlock, and process the logs into *PollResult for Emitter
	Poll(ctx context.Context, fromBlock, toBlock uint64) (*PollResult, error)
	// EmitterPoller also implements Controller
	Controller
}
