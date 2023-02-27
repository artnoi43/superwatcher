package poller

import (
	"fmt"

	"github.com/artnoi43/gsl"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

func collectLogs(
	m map[uint64]*mapLogsResult,
	logs []types.Log,
) (
	uint64,
	error,
) {
	var last uint64

	for i, log := range logs {
		number := log.BlockNumber
		hash := log.BlockHash
		result := m[number]

		switch {
		// nil result is when collectLogs have not yet visited this block
		case result == nil:
			fmt.Println("newBlock", number)
			m[number] = &mapLogsResult{
				Block: superwatcher.Block{
					Number: number,
					Hash:   log.BlockHash,
					Logs:   []*types.Log{&logs[i]},
				},
			}

		default:
			if h := result.Hash; h != hash {
				return number, errors.Wrapf(
					errHashesDiffer, "logs in block %d has different hashes: %s vs %s",
					number, hashStr(h), hashStr(hash),
				)
			}

			result.Logs = append(result.Logs, &logs[i])
		}
	}

	return last, nil
}

func collectHeaders(
	m map[uint64]*mapLogsResult,
	fromBlock uint64,
	toBlock uint64,
	headers map[uint64]superwatcher.BlockHeader,
) (
	uint64,
	error,
) {
	last := fromBlock
	for n := fromBlock; n <= toBlock; n++ {
		header, ok := headers[n]
		if !ok {
			continue
		}

		result, ok := m[n]
		last = n

		if ok {
			if rHash, hHash := result.Hash, header.Hash(); rHash != hHash {
				return last, errors.Wrapf(errHashesDiffer, "block %d resultHash differs from headerHash: %s vs %s",
					n, hashStr(rHash), hashStr(hHash),
				)
			}
		} else if !ok {
			m[n] = &mapLogsResult{
				Block: superwatcher.Block{
					Number: n,
					Header: header,
					Hash:   header.Hash(),
				},
			}
		}
	}

	return last, nil
}

func hashStr(h common.Hash) string {
	return gsl.StringerToLowerString(h)
}

// deleteUnusableResult removes all map keys that are >= lastGood.
// It used when poller.mapLogs encounter a situation where fresh chain data
// on a same block has different block hashes.
// func deleteMapResults[T any](pollResult map[uint64]T, lastGood uint64) {
// 	for k := range pollResult {
// 		if k >= lastGood {
// 			delete(pollResult, k)
// 		}
// 	}
// }
