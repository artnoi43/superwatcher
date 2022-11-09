package engine

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/wangjia184/sortedset"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/logger/debug"
)

type MetadataTracker interface {
	ClearUntil(blockNumber uint64)
	SetBlockMetadata(*superwatcher.BlockInfo, *blockMetadata)
	GetBlockMetadata(*superwatcher.BlockInfo) *blockMetadata
	SetBlockState(*superwatcher.BlockInfo, EngineBlockState)
	GetBlockState(*superwatcher.BlockInfo) EngineBlockState
}

// Tracker name needs revision!
type Tracker struct {
	sync.RWMutex

	// field "set" maps blockNumber to blockMetadata
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
// TODO: Currently broken
func (t *Tracker) ClearUntil(blockNumber uint64) {
	t.Lock()
	defer t.Unlock()

	debug.DebugMsg(t.debug, "clearing engine state tracker", zap.Uint64("until", blockNumber))

	for {
		oldest := t.set.PeekMin()
		if oldest == nil || oldest.Score() > sortedset.SCORE(blockNumber) {
			break
		}

		t.set.PopMin()
	}
}

func (t *Tracker) SetBlockMetadata(b *superwatcher.BlockInfo, metadata *blockMetadata) {
	t.Lock()
	defer t.Unlock()

	t.set.AddOrUpdate(b.BlockNumberString(), sortedset.SCORE(b.Number), metadata)
}

func (t *Tracker) GetBlockMetadata(b *superwatcher.BlockInfo) *blockMetadata {
	t.RLock()
	defer t.RUnlock()

	node := t.set.GetByKey(b.BlockNumberString())
	// Avoid panicking when assert type on nil value
	if node == nil {
		return &blockMetadata{blockNumber: b.Number}
	}

	metadata, ok := node.Value.(*blockMetadata)
	if !ok {
		logger.Panic(
			fmt.Sprintf("type assertion failed - expecting EngineLogState, found %s", reflect.TypeOf(node.Value)),
		)
	}

	return metadata
}

func (t *Tracker) SetBlockState(b *superwatcher.BlockInfo, state EngineBlockState) {
	metadata := t.GetBlockMetadata(b)
	if metadata == nil {
		// Create new metadata if null
		metadata = &blockMetadata{}
	}

	// Overwrite metadata.state
	metadata.state = state
	t.SetBlockMetadata(b, metadata)
}

func (t *Tracker) GetBlockState(b *superwatcher.BlockInfo) EngineBlockState {
	metadata := t.GetBlockMetadata(b)
	if metadata == nil {
		return EngineBlockStateNull
	}

	return metadata.state
}

func (t *Tracker) Len() int {
	t.RLock()
	defer t.RUnlock()

	return t.set.GetCount()
}
