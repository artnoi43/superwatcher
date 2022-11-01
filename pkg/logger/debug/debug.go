package debug

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
)

func DebugMsg(shouldPrint bool, msg string, fields ...zap.Field) {
	if shouldPrint {
		logger.Debug(msg, fields...)
	}
}
