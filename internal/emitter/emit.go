package emitter

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

func (e *emitter) emitFilterResult(result *superwatcher.FilterResult) {
	if result != nil {
		// Only log if there's some logs
		var nilChan bool
		if len(result.GoodBlocks)+len(result.ReorgedBlocks) != 0 {
			var goodBlocks []uint64
			for _, b := range result.GoodBlocks {
				goodBlocks = append(goodBlocks, b.Number)
			}

			var reorgedBlocks []uint64
			for _, b := range result.ReorgedBlocks {
				reorgedBlocks = append(reorgedBlocks, b.Number)
			}

			nilChan = e.filterResultChan == nil
			e.debugger.Debug(
				2, "emitFilterResult",
				zap.Uint64s("goodBlocks", goodBlocks),
				zap.Uint64s("reorgedBlocks", reorgedBlocks),
				zap.Bool("nil resultChan", nilChan),
			)
		}

		if !nilChan {
			e.filterResultChan <- result
		}

		return
	}

	logger.Panic("nil filterResult got sent to emitFilterREsult")
}

func (e *emitter) emitError(err error) {
	e.debugger.Debug(2, "emitError called")

	if err != nil {
		e.debugger.Debug(
			2, "blocking in emitError",
			// Use zap.String here because we don't want to log stack trace here
			zap.String("error to be sent", err.Error()),
		)
		e.errChan <- err
	}
}
