package emitter

import (
	"github.com/artnoi43/gsl/gslutils"
)

// * firstBlock = min(current block, lastRecordedBlock+1)
// * lastBlock = min(current block, firstBlock + filterRange)
// * FilterLog(firstBlock - filterRange, lastBlock)

// fromBlockToBlockNormal returns fromBlock and toBlock for emitter.filterLogs in **normal** circumstances.
// If the chain is reorging, or if there is any exception, use something else to compute the numbers.
func fromBlockToBlockNormal(
	emitterStartBlock uint64,
	currentBlock uint64,
	lastRecordedBlock uint64,
	filterRange uint64,
) (
	fromBlock uint64,
	toBlock uint64,
) {
	toBlock = gslutils.Min(currentBlock, lastRecordedBlock+filterRange)

	firstNewBlock := lastRecordedBlock + 1
	firstNewBlock = gslutils.Min(currentBlock, firstNewBlock)

	// Type uint64 can't be negative.
	// TODO: Check if this is right
	// e.g. firstNewBlock for next loop is 50, but the range is 100

	if filterRange > firstNewBlock {
		fromBlock = 0
	} else {
		fromBlock = firstNewBlock - filterRange
	}

	return fromBlock, toBlock
}
