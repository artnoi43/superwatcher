package emitter

import "github.com/artnoi43/gsl/gslutils"

// * firstBlock = min(current block, lastRecordedBlock+1)
// * lastBlock = min(current block, firstBlock + filterRange)
// * FilterLog(firstBlock - filterRange, lastBlock)

// fromBlockToBlock returns fromBlock and toBlock for emitter.filterLogs in **normal** circumstances.
// If the chain is reorging, or if there is any exception, use something else to compute the numbers.
func fromBlockToBlock(curr, lastRecorded, filterRange uint64) (fromBlock, toBlock uint64) {
	l := gslutils.Min([]uint64{curr, lastRecorded + 1})
	r := gslutils.Min([]uint64{curr, l + filterRange})

	// uint can't be negative.
	// e.g. filterRange 300, currentBlock 100 will lead to uint overflow
	if filterRange > lastRecorded {
		fromBlock = 0
	} else {
		fromBlock = l - filterRange
	}

	toBlock = r

	return fromBlock, toBlock
}
