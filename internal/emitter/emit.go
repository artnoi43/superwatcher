package emitter

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
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
	e.debugMsg("emitError called")

	if e.errChan == nil {
		e.debugMsg("emitError", zap.String("debug", "errChan is nil"))
		return
	}

	if err != nil {
		e.debugMsg("blocking in emitError")
		e.errChan <- err
	}
}
