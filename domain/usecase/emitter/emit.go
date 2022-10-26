package emitter

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
)

func (e *emitter) emitFilterResult(result *FilterResult) {
	if e.filterResultChan == nil {
		e.debugMsg("publishReorg", zap.String("debug", "filterResultChan is nil"))
		return
	}

	// Publish
	if result != nil {
		e.debugMsg("publishReorg", zap.Any("b", result))
		e.filterResultChan <- result
		return
	}

	logger.Panic("nil log sent to publishFilterResult")
}
