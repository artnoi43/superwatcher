package watcher

import (
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/watcher/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
)

func (w *watcher) publishLog(l *types.Log) {
	if l != nil {
		if w.debug {
			logger.Debug("publishLog", zap.Any("l", l))
		}
		w.logChan <- l
		return
	}

	logger.Panic("nil log sent to publishLog")
}

func (w *watcher) publishReorg(b *reorg.BlockInfo) {
	if b != nil {
		if w.debug {
			logger.Debug("publishReorg", zap.Any("b", b))
		}
		w.reorgChan <- b
		return
	}

	logger.Panic("nil log sent to publishReorg")
}
