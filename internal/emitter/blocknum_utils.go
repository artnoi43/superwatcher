package emitter

import (
	"github.com/soyart/gsl"
	"github.com/pkg/errors"
)

// fromBlockToBlockNormal returns fromBlock and toBlock for superwatcher.Poller.Poll in **normal** circumstances.
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
	// lastRecordedBlock = 80, filterRange = 10
	// 71  - 90   [normalCase] -> lastRecordedBlock = 90,  lookBack = 10, fwdRange = 90 - 80    = 10
	// 81  - 100  [normalCase] -> lastRecordedBlock = 100, lookBack = 10, fwdRange = 100 - 90   = 10
	// 91  - 110  [normalCase] -> lastRecordedBlock = 110, lookBack = 10, fwdRange = 110 - 100  = 10
	// 101 - 120  [normalCase] -> lastRecordedBlock = 120, lookBack = 10, fwdRange = 120 - 110  = 10

	toBlock = gsl.Min(currentBlock, lastRecordedBlock+filterRange)

	firstNewBlock := lastRecordedBlock + 1
	firstNewBlock = gsl.Min(currentBlock, firstNewBlock)

	if filterRange > firstNewBlock {
		fromBlock = emitterStartBlock
	} else {
		fromBlock = firstNewBlock - filterRange
	}

	return fromBlock, toBlock
}

func fromBlockToBlockIsReorging(
	emitterStartBlock uint64,
	currentBlock uint64,
	lastRecordedBlock uint64,
	filterRange uint64,
	maxRetries uint64,
	prevStatus *emitterStatus,
) (
	uint64,
	uint64,
	error,
) {
	// The lookBack range will grow after each retries, but not the forward range
	// lastRecordedBlock = 80, filterRange = 10, maxGoBackRetries = 3
	// 71  - 90   [normalCase]                -> lastRecordedBlock = 90,  lookBack = 10, fwdRange = 90 - 80   = 10
	// 81  - 100  [normalCase]                -> lastRecordedBlock = 100, lookBack = 10, fwdRange = 100 - 90  = 10
	// 91  - 110  [normalCase] 91 reorged     -> lastRecordedBlock = 110, lookBack = 10, fwdRange = 110 - 100 = 10
	// 81  - 110  # 81 reorged in this loop   -> lastRecordedBlock = 110, lookBack = 15, fwdRange = 110 - 110 = 0
	// 71  - 110  # 71 reorged in this loop   -> lastRecordedBlock = 110, lookBack = 20, fwdRange = 110 - 110 = 0
	// 61  - 110  # none reorged in this loop -> lastRecordedBlock = 110, lookBack = 25, fwdRange = 110 - 110 = 0
	// 101 - 120  [normalCase]                -> lastRecordedBlock = 120, lookBack = 10, fwdRange = 120 - 110 = 10

	if prevStatus.RetriesCount > maxRetries {
		// Why we use '>' instead of '>='
		// 0 status{0, false} -> normalFilter reorg -> status{1, true}
		// 1 status{1, true}  -> reorgFilter  reorg -> status{2, true} // 1st goBack
		// 2 status{2, true}  -> reorgFilter  reorg -> status{3, true} // 2nd goBack
		// 3 statis{3, true}  -> reorgFilter  reorg -> status{4, true} // 3rd goBack
		// 4 s 3 true,  maxReached

		return prevStatus.FromBlock, prevStatus.ToBlock, errors.Wrapf(
			ErrMaxRetriesReached, "%d goBackRetries", prevStatus.RetriesCount,
		)
	}

	// goBack in this case (reorging) is fixed
	goBack := filterRange
	var fromBlock, toBlock uint64
	// goBack is 1500, but (prev) fromBlock is 1050
	if goBack > prevStatus.FromBlock {
		fromBlock = emitterStartBlock
	} else {
		// goBack from the last fromBlock
		fromBlock = prevStatus.FromBlock - goBack
	}

	// toBlock does not go back, so we don't update it,
	// unless currentBlock was shrunk during reorg too.
	toBlock = gsl.Min(prevStatus.ToBlock, currentBlock)

	return fromBlock, toBlock, nil
}

func fromBlockToBlockGoBackFirstRun(
	emitterStartBlock uint64,
	currentBlock uint64,
	lastRecordedBlock uint64,
	filterRange uint64,
	maxRetries uint64,
	prevStatus *emitterStatus,
) (
	fromBlock, toBlock, goBack uint64,
) {
	goBack = filterRange * prevStatus.RetriesCount

	firstNewBlock := lastRecordedBlock + 1
	// Prevent overflow
	if goBack > firstNewBlock {
		fromBlock = emitterStartBlock
	} else {
		fromBlock = firstNewBlock - goBack
	}

	// The range is the same in firstStart
	toBlock = fromBlock + filterRange - 1

	return fromBlock, toBlock, goBack
}
