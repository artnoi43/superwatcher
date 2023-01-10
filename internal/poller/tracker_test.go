package poller

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
)

// TestUpdateTrackerValues tests if the tracker's values would change
// if we update the *Block from outside of tracker.
// This is important, because of how we handle Block.LogsMigrated.
func TestUpdateTrackerValues(t *testing.T) {
	b := &superwatcher.Block{
		Number: 69,
		Hash:   common.BigToHash(big.NewInt(69)),
		Logs:   nil,
	}

	tracker := newTracker("testTracker", 4)
	tracker.addTrackerBlock(b)

	h100 := common.BigToHash(big.NewInt(100))
	b.Hash = h100

	trackerBlock, found := tracker.getTrackerBlock(69)
	if !found {
		t.Fatal("missing tracker block")
	}

	if trackerBlock.Hash != b.Hash || trackerBlock.Hash != h100 {
		t.Log("b.Hash", b.Hash, "h100", h100, "trackerBlock.Hash", trackerBlock.Hash)
		t.Error("value in tracker not updated")
	}

	trackerBlock.Logs = []*types.Log{{}, {}}

	if len(b.Logs) == 0 || len(b.Logs) != len(trackerBlock.Logs) {
		t.Log("trackerBlock.Logs", trackerBlock.Logs, "b.Logs", b.Logs)
		t.Error("value in tracker not updated")
	}

	copied := *trackerBlock
	trackerBlock.Logs = nil
	trackerBlock.Hash = common.Hash{}

	if copied.Hash == trackerBlock.Hash || len(copied.Logs) == len(trackerBlock.Logs) {
		t.Error("unexpected copied value")
	}

	if err := tracker.removeBlock(69); err != nil {
		t.Error("cannot removed block 69")
	}

	_, ok := tracker.getTrackerBlock(69)
	if ok {
		t.Error("removed but found")
	}
}
