package watcher

import (
	"github.com/artnoi43/superwatcher/domain/usecase/watcher/reorg"
	"github.com/ethereum/go-ethereum/core/types"
)

func (w *watcher) publishLog(l *types.Log) {
	w.logChan <- l
}

func (w *watcher) publishReorg(b *reorg.BlockInfo) {
	w.reorgChan <- b
}
