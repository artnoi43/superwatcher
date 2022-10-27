package engine

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/wangjia184/sortedset"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

type MetadataTracker interface {
	ClearUntil(blockNumber uint64)
	SetBlockMetadata(*lib.BlockInfo, *blockMetadata)
	GetBlockMetadata(*lib.BlockInfo) *blockMetadata
	SetBlockState(*lib.BlockInfo, EngineBlockState)
	GetBlockState(*lib.BlockInfo) EngineBlockState
}

// Tracker name needs revision!
type Tracker struct {
	sync.Mutex

	// set maps blockNumber to blockMetadata
	set   *sortedset.SortedSet
	debug bool
}

func NewTracker(debug bool) *Tracker {
	return &Tracker{
		set:   sortedset.New(),
		debug: debug,
	}
}

// ClearUntil removes items in t from left to right.
func (t *Tracker) ClearUntil(blockNumber uint64) {
	debug.DebugMsg(t.debug, "clearing engine state tracker", zap.Uint64("until", blockNumber))

	t.Lock()
	defer t.Unlock()

	for {
		oldest := t.set.PeekMin()
		if oldest == nil || oldest.Score() > sortedset.SCORE(blockNumber) {
			break
		}
		t.set.PopMin()
	}
}

func (t *Tracker) SetBlockMetadata(b *lib.BlockInfo, metadata *blockMetadata) {
	t.Lock()
	defer t.Unlock()

	t.set.AddOrUpdate(b.BlockNumberString(), sortedset.SCORE(b.Number), metadata)
}

func (t *Tracker) GetBlockMetadata(b *lib.BlockInfo) *blockMetadata {
	t.Lock()
	defer t.Unlock()

	node := t.set.GetByKey(b.BlockNumberString())
	// Avoid panicking when assert type on nil value
	if node == nil {
		return nil
	}

	metadata, ok := node.Value.(*blockMetadata)
	if !ok {
		logger.Panic(
			fmt.Sprintf("type assertion failed - expecting EngineLogState, found %s", reflect.TypeOf(node.Value)),
		)
	}

	return metadata
}

func (t *Tracker) SetBlockState(b *lib.BlockInfo, state EngineBlockState) {
	metadata := t.GetBlockMetadata(b)
	if metadata == nil {
		// Create new metadata if null
		metadata = &blockMetadata{}
	}

	// Overwrite metadata.state
	metadata.state = state
	t.SetBlockMetadata(b, metadata)
}

func (t *Tracker) GetBlockState(b *lib.BlockInfo) EngineBlockState {
	metadata := t.GetBlockMetadata(b)
	if metadata == nil {
		return EngineBlockStateNull
	}

	return metadata.state
}

func (t *Tracker) Len() int {
	t.Lock()
	defer t.Unlock()

	return t.set.GetCount()
}
