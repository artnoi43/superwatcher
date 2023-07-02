package engine

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/wangjia184/sortedset"
	"go.uber.org/zap"

	"github.com/soyart/superwatcher/pkg/logger"
	"github.com/soyart/superwatcher/pkg/logger/debugger"
)

// TODO: remove callerMethod after debugging?
type callerMethod string

const (
	callerGoodLogs    = "callerGoodLogs"
	callerReorgedLogs = "callerReorgedLogs"
)

type metadataTracker interface {
	ClearUntil(blockNumber uint64)
	SetBlockMetadata(callerMethod, *blockMetadata)
	GetBlockMetadata(callerMethod, uint64, string) *blockMetadata
}

// metadataTrackerImpl is an in-memory store for keeping engine internal states.
// It is used to decide whether or not to pass the logs to service engine.
type metadataTrackerImpl struct {
	sync.RWMutex

	// Field `Tracker.sortedSet` maps txHash to blockMetadata.
	// The score is blockNumber. This allow us to use ClearUntil.
	sortedSet *sortedset.SortedSet
	debugger  *debugger.Debugger
}

func newTracker(debugLevel uint8) *metadataTrackerImpl {
	return &metadataTrackerImpl{
		sortedSet: sortedset.New(),
		debugger:  debugger.NewDebugger("metadataTracker", debugLevel),
	}
}

// ClearUntil removes items in t from left to right.
func (t *metadataTrackerImpl) ClearUntil(blockNumber uint64) {
	t.Lock()
	defer t.Unlock()

	t.debugger.Debug(
		2, "clearing engine state tracker",
		zap.Uint64("untilBlock", blockNumber),
	)

	for {
		oldest := t.sortedSet.PeekMin()
		if oldest == nil || oldest.Score() > sortedset.SCORE(blockNumber) {
			break
		}

		t.sortedSet.PopMin()
	}
}

func (t *metadataTrackerImpl) SetBlockMetadata(
	caller callerMethod,
	metadata *blockMetadata,
) {
	t.Lock()
	defer t.Unlock()

	t.debugger.Debug(
		3, "adding blockMetadata",
		zap.String("caller", string(caller)),
		zap.Uint64("blockNumber", metadata.blockNumber),
		zap.String("blockHash", metadata.blockHash),
		zap.Any("metadata artifacts", metadata.artifacts),
	)

	t.sortedSet.AddOrUpdate(metadata.blockHash, sortedset.SCORE(metadata.blockNumber), metadata)
}

func (t *metadataTrackerImpl) GetBlockMetadata(
	caller callerMethod,
	blockNumber uint64,
	blockHash string,
) *blockMetadata {
	t.RLock()
	defer t.RUnlock()

	node := t.sortedSet.GetByKey(blockHash)
	// Avoid panicking when assert type on nil value
	if node == nil {
		if caller == callerReorgedLogs {
			t.debugger.Debug(
				1, "nil metadata for reorged logs",
				zap.Uint64("blockNumber", blockNumber),
				zap.String("blockHash", blockHash),
			)

			panic(fmt.Sprintf("nil metadata for block %d blockHash %s", blockNumber, blockHash))
		}

		return &blockMetadata{
			blockNumber: blockNumber,
			blockHash:   blockHash,
		}
	}

	meta, ok := node.Value.(*blockMetadata)
	if !ok {
		logger.Panic(
			fmt.Sprintf("type assertion failed - expecting *blockMetadata, found %s", reflect.TypeOf(node.Value)),
		)
	}

	return meta
}

func (t *metadataTrackerImpl) Len() int {
	t.RLock()
	defer t.RUnlock()

	return t.sortedSet.GetCount()
}
