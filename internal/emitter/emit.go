package emitter

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"
)

func (e *emitter) emitFilterResult(result *superwatcher.FilterResult) {
	if e.filterResultChan == nil {
		e.debugMsg("emitFilterResult", zap.String("debug", "filterResultChan is nil"))
		return
	}

	if result != nil {
		e.debugMsg("emitFilterResult", zap.Any("b", result))
		e.filterResultChan <- result
		return
	}

	logger.Panic("nil filterResult got sent to emitFilterREsult")
}

func (e *emitter) emitError(err error) {
	if e.errChan == nil {
		e.debugMsg("emitError", zap.String("debug", "errChan is nil"))
	}

	if err != nil {
		e.errChan <- err
	}

	logger.Panic("nil error got sent to emitError")
}
