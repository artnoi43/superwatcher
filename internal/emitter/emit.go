package emitter

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

func (e *emitter) emitFilterResult(result *superwatcher.FilterResult) {
	if e.filterResultChan == nil {
		e.debugger.Debug(2, "emitFilterResult", zap.String("debug", "filterResultChan is nil"))
		return
	}

	if result != nil {
		// Only log if there's some logs
		if len(result.GoodBlocks)+len(result.ReorgedBlocks) != 0 {
			if e.debug {
				var goodBlocks []uint64
				for _, b := range result.GoodBlocks {
					goodBlocks = append(goodBlocks, b.Number)
				}

				var reorgedBlocks []uint64
				for _, b := range result.ReorgedBlocks {
					reorgedBlocks = append(reorgedBlocks, b.Number)
				}

				e.debugger.Debug(2, "emitFilterResult", zap.Uint64s("goodBlocks", goodBlocks), zap.Uint64s("reorgedBlocks", reorgedBlocks))
			}
		}

		e.filterResultChan <- result
		return
	}

	logger.Panic("nil filterResult got sent to emitFilterREsult")
}

func (e *emitter) emitError(err error) {
	e.debugger.Debug(2, "emitError called")

	if e.errChan == nil {
		e.debugger.Debug(2, "emitError", zap.String("debug", "errChan is nil"))
		return
	}

	if err != nil {
		e.debugger.Debug(2, "blocking in emitError")
		e.errChan <- err
	}
}
