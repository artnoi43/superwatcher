package engine

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/wangjia184/sortedset"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

type EngineStateTracker interface {
	SetEngineState(engineLogStateKey, EngineLogState)
	GetEngineState(engineLogStateKey) EngineLogState
	ClearUntil(blockNumber uint64)
}

// Tracker name needs revision!
type Tracker struct {
	sync.Mutex

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

// AddBlockInfo add a new *BlockInfo to t
func (t *Tracker) SetEngineState(key engineLogStateKey, state EngineLogState) {
	t.Lock()
	defer t.Unlock()

	t.set.AddOrUpdate(key.String(), sortedset.SCORE(key.blockNumber), state)
}

// GetSavedBlockByBlockNumber returns *BlockInfo from t with key blockNumber
func (t *Tracker) GetEngineState(key engineLogStateKey) EngineLogState {
	t.Lock()
	defer t.Unlock()

	node := t.set.GetByKey(key.String())
	if node == nil {
		return EngineLogStateNull
	}

	logger.Debug("post nil")
	state, ok := node.Value.(EngineLogState)
	if !ok {
		logger.Panic(fmt.Sprintf("type assertion failed - expecting EngineLogState, found %s", reflect.TypeOf(node.Value)))
	}

	logger.Debug("ret state")
	return state
}

func (t *Tracker) Len() int {
	t.Lock()
	defer t.Unlock()

	return t.set.GetCount()
}
