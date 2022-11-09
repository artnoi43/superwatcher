package emitter

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/wangjia184/sortedset"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

// blockTracker name needs revision!
type blockTracker struct {
	set *sortedset.SortedSet
}

func newTracker() *blockTracker {
	return &blockTracker{
		set: sortedset.New(),
	}
}

// AddBlockInfo add a new *BlockInfo to t
func (t *blockTracker) addTrackerBlock(b *superwatcher.BlockInfo) {
	k := strconv.FormatInt(int64(b.Number), 10)
	t.set.AddOrUpdate(k, sortedset.SCORE(b.Number), b)
}

// clearUntil removes items in t from left to right.
func (t *blockTracker) clearUntil(blockNumber uint64) {
	for {
		oldest := t.set.PeekMin()
		if oldest == nil || oldest.Score() > sortedset.SCORE(blockNumber) {
			break
		}
		t.set.PopMin()
	}
}

// GetSavedBlockByBlockNumber returns *BlockInfo from t with key blockNumber
func (t *blockTracker) getTrackerBlockInfo(blockNumber uint64) (*superwatcher.BlockInfo, bool) {
	k := strconv.FormatUint(blockNumber, 10)
	node := t.set.GetByKey(k)
	if node == nil {
		return nil, false
	}
	val, ok := node.Value.(*superwatcher.BlockInfo)
	if !ok {
		logger.Panic(fmt.Sprintf("type assertion failed - expecting *BlockInfo, found %s", reflect.TypeOf(node.Value)))
	}
	return val, true
}

func (t *blockTracker) Len() int {
	return t.set.GetCount()
}
