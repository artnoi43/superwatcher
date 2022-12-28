package poller

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
)

func TestTracker(t *testing.T) {
	b := &superwatcher.BlockInfo{
		Number: 69,
		Hash:   common.BigToHash(big.NewInt(69)),
	}

	tracker := newTracker("testTracker", 4)
	tracker.addTrackerBlockInfo(b)

	h100 := common.BigToHash(big.NewInt(100))
	b.Hash = h100

	_b, found := tracker.getTrackerBlockInfo(69)
	if !found {
		t.Fatal("missing tracker blockInfo")
	}

	if _b.Hash != b.Hash || _b.Hash != h100 {
		t.Log("b.Hash", b.Hash, "h100", h100, "_b.Hash", _b.Hash)
	}
}
