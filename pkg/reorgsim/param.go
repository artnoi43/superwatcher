package reorgsim

import "github.com/pkg/errors"

// BaseParam is the basic parameters for the mock client. Chain reorg parameters are NOT included here.
type BaseParam struct {
	// StartBlock will be used as initial ReorgSim.currentBlock. ReorgSim.currentBlock increases by `BlockProgess`
	// after each call to ReorgSim.BlockNumber.
	StartBlock uint64 `json:"startBlock"`
	// BlockProgress represents how many block numbers should ReorgSim.currentBlock
	// should increase during each call to ReorgSim.BlockNumber.
	BlockProgress uint64 `json:"blockProgress"`
	// ExitBlock is used in ReorgSim.BlockNumber to return ErrExitBlockReached once its currentBlock reaches ExitBlock.
	ExitBlock uint64 `json:"exitBlock"`

	Debug bool `json:"-"`
}

var DefaultParam = BaseParam{
	BlockProgress: 20,
	Debug:         true,
}

// ReorgEvent is parameters for chain reorg events.
type ReorgEvent struct {
	// TODO: add ReorgTrigger - the block which when seen, triggers a reorg from ReorgBlock
	// ReorgTrigger uint64 `json:"reorgTrigger"`

	// ReorgBlock is the pivot block after which ReorgSim should use ReorgSim.reorgedChain.
	ReorgBlock uint64 `json:"reorgBlock"`
	// MovedLogs represents all of the moved logs after a chain reorg event. The map key is the source blockNumber.
	MovedLogs map[uint64][]MoveLogs `json:"movedLogs"`
}

var errInvalidReorgEvents = errors.New("invalid reorg events")
