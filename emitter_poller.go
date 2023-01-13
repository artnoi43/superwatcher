package superwatcher

import (
	"context"
	"fmt"
)

// EmitterPoller fetches event logs and other blockchain data and maps it into *PollerResult.
// The result of EmitterPoller.poll is later used by Emitter to emit to Engine.
type EmitterPoller interface {
	// Poll polls event logs from fromBlock to toBlock, and process the logs into *PollerResult for Emitter
	Poll(ctx context.Context, fromBlock, toBlock uint64) (*PollerResult, error)
	// Policy gets current Policy
	Policy() Policy
	// SetPolicy sets new Policy (NOTE: changing Policy mid-run not tested)
	SetPolicy(Policy) error

	// EmitterPoller also implements Controller
	Controller
}

// Policy (enum) specifies how EmitterPoller considers which blocks
// to include in its _tracking list_.
// For every block in this _tracking list_, the poller compares
// the saved block hash with newly polled one.
type Policy uint8

const (
	// PolicyFast makes poller only process and track blocks with interesting logs.
	// Hashes from blocks without logs are not processed, unless they were reorged
	// and had their logs removed, in which case the poller gets their headers _once_
	// to check their newer hashes, and remove the empty block from tracking list.
	PolicyFast Policy = iota

	// PolicyNormal makes poller only process and track blocks with interesting logs,
	// but if the poller detects that a block has its logs removed, it will process
	// and start tracking that block until the block goes out of poller scope.
	// The difference between PolicyFast and PolicyNormal
	// that PolicyNormal will keep tracking the reorged empty blocks.
	PolicyNormal

	// PolicyExpensive makes poller process and track all blocks' headers,
	// regardless of whether the blocks have interesting logs or not, or Config.DoHeader value.
	PolicyExpensive

	// PolicyFastBlock behaves like PolicyExpensive, but instead of fetching
	// event logs and headers, the poller fetches event logs and blocks.
	PolicyExpensiveBlock
)

func (level Policy) String() string {
	switch level {
	case PolicyFast:
		return "FAST"
	case PolicyNormal:
		return "NORMAL"
	case PolicyExpensive:
		return "EXPENSIVE"
	case PolicyExpensiveBlock:
		return "EXPENSIVE_BLOCK"
	}

	return fmt.Sprintf("UNKNOWN LEVEL %d", level)
}
