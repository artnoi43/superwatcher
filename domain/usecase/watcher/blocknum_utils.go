package watcher

import "github.com/artnoi43/gsl/gslutils"

// * firstBlock = min(current block, lastRecordedBlock+1)
// * lastBlock = min(current block, firstBlock + lookback_blocks)
// * FilterLog(firstBlock - lookback_blocks, lastBlock)

func fromBlockToBlock(curr, lastRec, lookBack uint64) (fromBlock, toBlock uint64) {
	l := gslutils.Min([]uint64{curr, lastRec + 1})
	r := gslutils.Min([]uint64{curr, l + lookBack})
	fromBlock = l - lookBack
	toBlock = r
	// uint can't be negative.
	// e.g. lookBack 300, currentBlock 100 will lead to uint overflow
	if lookBack > lastRec {
		fromBlock = 0
	}
	return fromBlock, toBlock
}
