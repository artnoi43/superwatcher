package engine

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
)

func newBlockInfo(number uint64) *superwatcher.BlockInfo {
	return &superwatcher.BlockInfo{
		Number: number,
		Hash:   common.BigToHash(big.NewInt(int64(number))),
	}
}

func TestMetadataTracker(t *testing.T) {
	tracker := NewTracker(true)

	// GetBlockMetadata should not return nil even if it's empty
	block69 := newBlockInfo(69)
	if met := tracker.GetBlockMetadata(block69); met == nil {
		t.Fatal("GetBlockMetadata returns nil")
	} else {
		if met.state != StateNull {
			t.Fatalf("expecting Null state, got %s", met.state.String())
		}
	}

	met69 := tracker.GetBlockMetadata(block69)
	met69.state.Fire(EventGotLog)
	if met69.state != StateSeen {
		t.Fatalf("expecting Seen state, got %s", met69.state.String())
	}

	// Overwrite met69 with blank metadata from GetBlockMetadata
	met69 = tracker.GetBlockMetadata(block69)
	if met69.state != StateNull {
		t.Fatalf("expecting Null state (did not save back yet), got %s", met69.state.String())
	}

	// Copy state reference out and fire on it - and `metadata.state` should change too
	state := met69.state
	state.Fire(EventGotLog)

	if met69.state != StateNull {
		t.Fatalf("expecing met69.state to remain Null, got %s", met69.state.String())
	}

	// Overwrite met69 with blank metadata from GetBlockMetadata
	met69 = tracker.GetBlockMetadata(block69)
	met69.state.Fire(EventGotLog)
	if met69.state != StateSeen {
		t.Fatalf("expecing met69.state to change to Seen, got %s", met69.state.String())
	}

	// Save back
	tracker.SetBlockMetadata(block69, met69)
	// And get it out again - the state should remain Seen
	met69 = tracker.GetBlockMetadata(block69)
	if met69.state != StateSeen {
		t.Fatalf("expecing met69.state to change to Seen, got %s", met69.state.String())
	}

	// State should remain Seen
	met69.state.Fire(EventGotLog)
	met69.state.Fire(EventGotLog)
	met69.state.Fire(EventGotLog)
	if met69.state != StateSeen {
		t.Fatalf("expecing met69.state to remain Seen, got %s", met69.state.String())
	}
}
