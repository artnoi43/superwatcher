package mock

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

// TODO: refactor

func TestMockPoller(t *testing.T) {
	reorgBlocks := []uint64{125, 143}
	reorgIndex := 0
	seen := make(map[uint64]int)

	filterRange := uint64(10)
	lastRec := uint64(100)
	endBlock := uint64(172)

	poller := NewPoller(reorgBlocks)

	var fromBlockReorged bool
	var fromBlock, toBlock uint64
	for lastRec <= endBlock {

		if fromBlockReorged {
			// TODO: don't break but figure out appropriate fromBlock, toBlock
			break

		} else {
			fromBlock = lastRec + 1 - filterRange
			toBlock = lastRec + filterRange
		}

		currentReorgBlock := reorgBlocks[reorgIndex]

		t.Log("from", fromBlock, "to", toBlock, "currentReorg", currentReorgBlock)

		for i := fromBlock; i <= toBlock; i++ {
			seen[i]++
		}

		if seen[currentReorgBlock] > 1 {
			if reorgIndex < len(reorgBlocks)-1 {
				reorgIndex++
			}
		}

		result, err := poller.Poll(nil, fromBlock, toBlock)

		if result != nil {
			fromBlockReorged = false

			t.Log(
				"from", fromBlock, "to", toBlock, "currentReorg", currentReorgBlock,
				"lastGood", result.LastGoodBlock,
				"goodBlocks", len(result.GoodBlocks), "reorgedBlocks", len(result.ReorgedBlocks),
			)

			lastRec = result.LastGoodBlock
			continue
		}

		if currentReorgBlock == fromBlock {
			if err != nil {
				t.Log("error from Poll", err.Error())
				if !errors.Is(err, superwatcher.ErrFromBlockReorged) {
					t.Error("expecting errFromBlockReorged, got", err.Error())
				}

				lastRec = fromBlock - 1
				fromBlockReorged = true
			}

			lastRec = result.LastGoodBlock
		}
	}
}
