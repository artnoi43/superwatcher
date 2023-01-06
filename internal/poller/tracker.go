package poller

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/wangjia184/sortedset"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// blockTracker stores the `superwatcher.Block` with blockNumber as key.
// It is used by poller to store `Block` from the last call of `poller.poll`
// to see if a block's hash has changed.
// TODO: blockTracker name needs revision!
// TODO: Remove caller and logging
type blockTracker struct {
	sortedSet *sortedset.SortedSet
	user      string
	debugger  *debugger.Debugger
}

func newTracker(user string, debugLevel uint8) *blockTracker {
	key := fmt.Sprintf("blockTracker for %s", user)

	return &blockTracker{
		sortedSet: sortedset.New(),
		user:      user,
		debugger:  debugger.NewDebugger(key, debugLevel),
	}
}

// addTrackerBlock adds `*Block` |b| to the store using |b.Number| as key
func (t *blockTracker) addTrackerBlock(b *superwatcher.Block) {
	t.debugger.Debug(
		3,
		"adding block",
		zap.Uint64("blockNumber", b.Number),
		zap.String("blockHash", b.String()),
		zap.Int("lenLogs", len(b.Logs)),
	)

	k := strconv.FormatUint(b.Number, 10)
	t.sortedSet.AddOrUpdate(k, sortedset.SCORE(b.Number), b)
}

// getTrackerBlock returns `*Block` from t with key |blockNumber|
func (t *blockTracker) getTrackerBlock(blockNumber uint64) (*superwatcher.Block, bool) {
	k := strconv.FormatUint(blockNumber, 10)
	node := t.sortedSet.GetByKey(k)
	if node == nil {
		return nil, false
	}

	val, ok := node.Value.(*superwatcher.Block)
	if !ok {
		logger.Panic(fmt.Sprintf("type assertion failed - expecting *Block, found %s", reflect.TypeOf(node.Value)))
	}

	return val, true
}

func (t *blockTracker) removeBlock(blockNumber uint64) error {
	k := strconv.FormatUint(blockNumber, 10)
	del := t.sortedSet.Remove(k)
	if del == nil {
		return errors.New("node was not in set")
	}

	return nil
}

// clearUntil removes `*Block` in t from left to right.
func (t *blockTracker) clearUntil(blockNumber uint64) {
	for {
		oldest := t.sortedSet.PeekMin()
		if oldest == nil || oldest.Score() > sortedset.SCORE(blockNumber) {
			break
		}

		t.sortedSet.PopMin()
	}
}

func (t *blockTracker) Len() int {
	return t.sortedSet.GetCount()
}
