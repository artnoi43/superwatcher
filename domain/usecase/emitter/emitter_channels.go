package emitter

import (
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
)

func (e *emitter) publishLog(l *types.Log) {
	if l != nil {
		if e.debug {
			logger.Debug("publishLog", zap.Any("l", l))
		}
		e.logChan <- l
		return
	}

	logger.Panic("nil log sent to publishLog")
}

func (e *emitter) publishReorg(b *reorg.BlockInfo) {
	if b != nil {
		if e.debug {
			logger.Debug("publishReorg", zap.Any("b", b))
		}
		e.reorgChan <- b
		return
	}

	logger.Panic("nil log sent to publishReorg")
}
