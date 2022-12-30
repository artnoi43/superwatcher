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

func assertState(t *testing.T, expected, actual blockState) {
	if expected != actual {
		t.Fatalf("expecting %s state, got %s", expected.String(), actual.String())
	}
}

// TODO: Rewrite
func TestMetadataTracker(t *testing.T) {
	tracker := newTracker(3)
	trackerKey := callerMethod("testTracker")

	// GetBlockMetadata should not return nil even if it's empty
	block69 := newBlockInfo(69)
	if met := tracker.GetBlockMetadata(trackerKey, block69.Number, block69.String()); met == nil {
		t.Fatal("GetBlockMetadata returns nil")
	} else {
		assertState(t, stateNull, met.state)
	}

	met69 := tracker.GetBlockMetadata(trackerKey, block69.Number, block69.String())

	met69.state.Fire(eventSeeBlock)
	assertState(t, stateSeen, met69.state)

	met69.state.Fire(eventHandle)
	assertState(t, stateHandled, met69.state)

	met69.state.Fire(eventSeeBlock)
	assertState(t, stateHandled, met69.state)

	// Overwrite met69 with blank metadata from GetBlockMetadata
	met69 = tracker.GetBlockMetadata(trackerKey, block69.Number, block69.String())
	assertState(t, stateNull, met69.state)

	// Copy state reference out and fire on it - and `metadata.state` should change too
	state := &met69.state
	state.Fire(eventSeeBlock)
	assertState(t, stateSeen, met69.state)
	state.Fire(eventHandle)
	assertState(t, stateHandled, met69.state)

	// Overwrite met69 with blank metadata from GetBlockMetadata
	met69 = tracker.GetBlockMetadata(trackerKey, block69.Number, block69.String())
	met69.state.Fire(eventSeeBlock)
	met69.state.Fire(eventHandle)
	assertState(t, stateHandled, met69.state)

	met69.state.Fire(eventSeeReorg)
	assertState(t, stateReorged, met69.state)

	tracker.SetBlockMetadata(trackerKey, met69)
	met69 = tracker.GetBlockMetadata(trackerKey, block69.Number, block69.String())
	assertState(t, stateReorged, met69.state)

	// Save metadata with artifacts back to tracker
	tracker.SetBlockMetadata(trackerKey, met69)
	// Delete met69 data
	met69 = nil
	// Re-get the saved metadata, and it should remain stateReorged
	met69 = tracker.GetBlockMetadata(trackerKey, block69.Number, block69.String())
	assertState(t, stateReorged, met69.state)

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
