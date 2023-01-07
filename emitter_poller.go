package superwatcher

import (
	"context"
	"fmt"
)

// EmitterPoller filters event logs from the blockchain and maps []types.Log into *PollResult.
// The result of EmitterPoller.poll is later used by Emitter to emit to Engine.
// superwatcher users can ignore this type if they have no need to update log addresses and topics on-the-fly,
// as EmitterPoller is already wrapped by Emitter.
type EmitterPoller interface {
	// Poll polls event logs from fromBlock to toBlock, and process the logs into *PollResult for Emitter
	Poll(ctx context.Context, fromBlock, toBlock uint64) (*PollResult, error)
	// PollLevel gets current PollLevel
	PollLevel() PollLevel
	// SetPollLevel sets new PollLevel (NOTE: changing PollLevel mid-run not tested)
	SetPollLevel(PollLevel) error
	// EmitterPoller also implements Controller
	Controller
}

// PollLevel (enum) specifies how EmitterPoller considers which blocks to include in its _tracking list_,
// For every block in this _tracking list_, the poller compares the saved block hash with newly polled one.
type PollLevel uint8

const (
	// PollLevelFast makes poller only process and track blocks with interesting logs. Hashes from blocks
	// without logs are discarded, unless they were reorged and had their logs removed, in which case
	// the poller gets their headers _once_ to check their newer hashes, and remove the empty block from tracking list.
	PollLevelFast PollLevel = iota

	// PollLevelNormal makes poller only process and track blocks with interesting logs,
	// but if the poller detects that a block has its logs removed, it will process and track that block
	// until the block goes out of poller scope. The difference between PollLevelFast and PollLevelNormal
	// is that PollLevelNormal will keep tracking the reorged empty blocks.
	PollLevelNormal

	// PollLevelExpensive makes poller process and track all blocks' headers,
	// regardless of whether the blocks have interesting logs or not, or Config.DoHeader value.
	PollLevelExpensive
)

func (level PollLevel) String() string {
	switch level {
	case PollLevelFast:
		return "FAST"
	case PollLevelNormal:
		return "NORMAL"
	case PollLevelExpensive:
		return "EXPENSIVE"
	}

	return fmt.Sprintf("UNKNOWN LEVEL %d", level)
}
