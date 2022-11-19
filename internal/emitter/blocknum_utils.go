package emitter

import "github.com/artnoi43/gsl/gslutils"

// * firstBlock = min(current block, lastRecordedBlock+1)
// * lastBlock = min(current block, firstBlock + lookbackBlocks)
// * FilterLog(firstBlock - lookbackBlocks, lastBlock)

// fromBlockToBlock returns fromBlock and toBlock for emitter.filterLogs in **normal** circumstances.
// If the chain is reorging, or there is any exception, use something else to compute the numbers.
func fromBlockToBlock(curr, lastRec, lookBack uint64) (fromBlock, toBlock uint64) {
	l := gslutils.Min([]uint64{curr, lastRec + 1})
	r := gslutils.Min([]uint64{curr, l + lookBack})

	// uint can't be negative.
	// e.g. lookBack 300, currentBlock 100 will lead to uint overflow
	if lookBack > lastRec {
		fromBlock = 0
	} else {
		fromBlock = l - lookBack
	}

	toBlock = r

	return fromBlock, toBlock
}
