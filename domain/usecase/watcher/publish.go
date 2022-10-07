package watcher

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func (w *watcher) PublishLog(l *types.Log) error {
	w.logChan <- l
	return nil
}
