package debugger

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
)

type Debugger struct {
	Key   string // Debugger key - should represent something where this debugger is embedded or used.
	Level uint8  // Level maps to different verbosity level, higher being more verbose
}

func NewDebugger(key string, level uint8) *Debugger {
	return &Debugger{
		Key:   key,
		Level: level,
	}
}

func (d *Debugger) Debug(level uint8, msg string, fields ...zap.Field) {
	if d.Level >= level {
		msg = fmt.Sprintf("%s: %s", d.Key, msg)
		logger.Debug(msg, fields...)
	}
}
