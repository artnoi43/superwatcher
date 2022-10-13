package emitter

import (
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/logger"
)

func (e *emitter) publishLog(l *types.Log) {
	if e.logChan == nil {
		e.debugMsg("publishLog", zap.String("debug", "logChan is nil"))
		return
	}

	// Publish
	if l != nil {
		e.debugMsg("publishLog", zap.Any("l", l))
		e.logChan <- l
		return
	}

	logger.Panic("nil log sent to publishLog")
}

func (e *emitter) publishBlock(b *reorg.BlockInfo) {
	if e.blockChan == nil {
		e.debugMsg("publishBlock", zap.String("debug", "blockChan is nil"))
		return
	}

	// Publish
	if b != nil {
		e.debugMsg("publishBlock", zap.Uint64("blockNumber", b.Number), zap.Int("logs", len(b.Logs)))
		e.blockChan <- b
		return
	}

	logger.Panic("nil log sent to publishBlock")
}

func (e *emitter) publishReorg(b *reorg.BlockInfo) {
	if e.reorgChan == nil {
		e.debugMsg("publishReorg", zap.String("debug", "reorgChan is nil"))
		return
	}

	// Publish
	if b != nil {
		e.debugMsg("publishReorg", zap.Any("b", b))
		e.reorgChan <- b
		return
	}

	logger.Panic("nil log sent to publishReorg")
}
