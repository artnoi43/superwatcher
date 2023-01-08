package poller

import (
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
)

//
// import "github.com/ethereum/go-ethereum/core/types"
//
// func collectLogs(logs []types.Log) map[uint64][]*types.Log {
// 	m := make(map[uint64][]*types.Log)
//     for i, log := range logs {
//     }
// }

func hashStr(h common.Hash) string {
	return gslutils.StringerToLowerString(h)
}

// deleteUnusableResult removes all map keys that are >= lastGood.
// It used when poller.mapLogs encounter a situation where fresh chain data
// on a same block has different block hashes.
func deleteMapResults[T any](pollResult map[uint64]T, lastGood uint64) {
	for k := range pollResult {
		if k >= lastGood {
			delete(pollResult, k)
		}
	}
}
