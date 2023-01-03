package reorgsim

import "github.com/pkg/errors"

// Param is the basic parameters for the mock client. Chain reorg parameters are NOT included here.
type Param struct {
	// StartBlock will be used as initial ReorgSim.currentBlock.
	StartBlock uint64 `json:"startBlock"`
	// BlockProgress is used to increment ReorgSim.currentBlock.
	BlockProgress uint64 `json:"blockProgress"`
	// ExitBlock is checked against ReorgSim.currentBlock for test code to exit at a specify block.
	ExitBlock uint64 `json:"exitBlock"`

	Debug bool `json:"-"`
}

// ReorgEvent is parameters for chain reorg events.
type ReorgEvent struct {
	// ReorgTrigger is the block which when seen, triggers a reorg from ReorgBlock
	ReorgTrigger uint64 `json:"reorgTrigger"`
	// ReorgBlock is the pivot block after which ReorgSim should use ReorgSim.reorgedChain.
	ReorgBlock uint64 `json:"reorgBlock"`
	// MovedLogs represents all of the moved logs after a chain reorg event.
	// The map key is the source block number (the block from which the logs are originally in).
	MovedLogs map[uint64][]MoveLogs `json:"movedLogs"`
}

var (
	errInvalidReorgEvents = errors.New("invalid reorg events")
	DefaultParam          = Param{
		BlockProgress: 20,
		Debug:         true,
	}
)

// validateReorgEvent validates order of |events|,
// and will use event.ReorgBlock as event.ReorgTrigger if the latter is 0,
func validateReorgEvent(events []ReorgEvent) ([]ReorgEvent, error) {
	for i, event := range events {
		// Overwrites ReorgTrigger if is 0
		// or invalidates if reorgBlock > reorgTrigger
		if event.ReorgTrigger == 0 {
			events[i].ReorgTrigger = events[i].ReorgBlock
		}

		for from := range event.MovedLogs {
			if from < event.ReorgBlock {
				return nil, errors.Wrapf(
					errInvalidReorgEvents, "logs moved from %d, which is before reorgBlock %d",
					from, event.ReorgBlock,
				)
			}
		}
	}

	return events, nil
}
