package debugger

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
)

type Debugger struct {
	Key         string
	ShouldDebug bool
}

func (d *Debugger) Debug(msg string, fields ...zap.Field) {
	msg = fmt.Sprintf("%s: %s", d.Key, msg)
	if d.ShouldDebug {
		Debug(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}
