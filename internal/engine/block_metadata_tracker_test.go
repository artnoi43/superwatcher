package engine

import (
	"math/big"
	"reflect"
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
	tracker := NewTracker(3)
	// GetBlockMetadata should not return nil even if it's empty
	block69 := newBlockInfo(69)
	if met := tracker.GetBlockMetadata("test", block69); met == nil {
		t.Fatal("GetBlockMetadata returns nil")
	} else {
		if met.state != StateNull {
			t.Fatalf("expecting Null state, got %s", met.state.String())
		}
	}

	met69 := tracker.GetBlockMetadata("test", block69)
	met69.state.Fire(EventGotLog)
	if met69.state != StateSeen {
		t.Fatalf("expecting Seen state, got %s", met69.state.String())
	}

	// Overwrite met69 with blank metadata from GetBlockMetadata
	met69 = tracker.GetBlockMetadata("test", block69)
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
	met69 = tracker.GetBlockMetadata("test", block69)
	met69.state.Fire(EventGotLog)
	if met69.state != StateSeen {
		t.Fatalf("expecing met69.state to change to Seen, got %s", met69.state.String())
	}

	// Save back
	tracker.SetBlockMetadata("test", block69, met69)
	// And get it out again - the state should remain Seen
	met69 = tracker.GetBlockMetadata("test", block69)
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

	type foo struct {
		a int
		b string
	}
	type bar struct {
		x int
		y string
	}

	met69.artifacts = []superwatcher.Artifact{
		&foo{
			a: 69, b: "foo69",
		},
		&bar{
			x: 69, y: "bar69",
		},
	}

	// Save metadata with artifacts back to tracker
	tracker.SetBlockMetadata("test", block69, met69)
	met69 = nil

	met69 = tracker.GetBlockMetadata("test", block69)
	if met69.state != StateSeen {
		t.Fatalf("expecting met69.state to remain Seen, got %s", met69.state.String())
	}
	for i, art := range met69.artifacts {
		switch i {
		case 0:
			fooArt, ok := art.(*foo)
			if !ok {
				t.Fatalf("artifact 0 failed type assertion, expecting *foo, got %s", reflect.TypeOf(art).String())
			}
			if fooArt.a != 69 {
				t.Fatalf("fooArt.a is not 69, got %d", fooArt.a)
			}
			if fooArt.b != "foo69" {
				t.Fatalf("fooArt.b is not foo69, got %s", fooArt.b)
			}
		case 1:
			barArt, ok := art.(*bar)
			if !ok {
				t.Fatalf("artifact 1 failed type assertion, expecting *bar, got %s", reflect.TypeOf(art).String())
			}
			if barArt.x != 69 {
				t.Fatalf("barArt.x is not 69, got %d", barArt.x)
			}
			if barArt.y != "bar69" {
				t.Fatalf("barArt.y is not bar69, got %s", barArt.y)
			}
		}
	}
}

func returnArtifacts() []superwatcher.Artifact {
	type dummyArtifact string
	return []superwatcher.Artifact{
		[]dummyArtifact{
			"dummy0",
			"dummy1",
			"dummy2",
		},
	}
}

func TestArtifact(t *testing.T) {
	artifacts := returnArtifacts()
	t.Log("type of artifacts", reflect.TypeOf(artifacts))

	if l := len(artifacts); l != 1 {
		t.Fatalf("expecting len artifacts of 1, got %d", l)
	}

	for _, artifact := range artifacts {
		t.Log("type of artifact", reflect.TypeOf(artifact))
	}
}
