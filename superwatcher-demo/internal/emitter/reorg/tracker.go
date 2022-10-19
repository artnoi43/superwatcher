package reorg

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/wangjia184/sortedset"

	"github.com/artnoi43/superwatcher/lib/logger"
)

// Tracker name needs revision!
type Tracker struct {
	set *sortedset.SortedSet
}

func NewTracker() *Tracker {
	return &Tracker{
		set: sortedset.New(),
	}
}

// AddBlockInfo add a new *BlockInfo to t
func (t *Tracker) AddTrackerBlock(b *BlockInfo) {
	k := strconv.FormatInt(int64(b.Number), 10)
	t.set.AddOrUpdate(k, sortedset.SCORE(b.Number), b)
}

// ClearUntil removes items in t from left to right.
func (t *Tracker) ClearUntil(blockNumber uint64) {
	for {
		oldest := t.set.PeekMin()
		if oldest == nil || oldest.Score() > sortedset.SCORE(blockNumber) {
			break
		}
		t.set.PopMin()
	}
}

// GetSavedBlockByBlockNumber returns *BlockInfo from t with key blockNumber
func (t *Tracker) GetTrackerBlockInfo(blockNumber uint64) (*BlockInfo, bool) {
	k := strconv.FormatUint(blockNumber, 10)
	node := t.set.GetByKey(k)
	if node == nil {
		return nil, false
	}
	val, ok := node.Value.(*BlockInfo)
	if !ok {
		logger.Panic(fmt.Sprintf("type assertion failed - expecting *BlockInfo, found %s", reflect.TypeOf(node.Value)))
	}
	return val, true
}

func (t *Tracker) Len() int {
	return t.set.GetCount()
}
