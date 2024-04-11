package logger

import (
	"fmt"
	"go.uber.org/zap"
)

type ZapAdapter struct {
	logger *zap.Logger
}

func NewZapAdapter(logger *zap.Logger) TOLogger {
	return &ZapAdapter{logger: logger}
}

func (z ZapAdapter) Debug(args ...interface{}) {
	if z.logger.Core().Enabled(zap.DebugLevel) {
		logStr, a := getArgs(args...)
		z.logger.Debug(fmt.Sprintf(logStr, a...))
	}
}

func (z ZapAdapter) Info(args ...interface{}) {
	if z.logger.Core().Enabled(zap.InfoLevel) {
		logStr, a := getArgs(args...)
		z.logger.Info(fmt.Sprintf(logStr, a...))
	}
}

func (z ZapAdapter) Warn(args ...interface{}) {
	if z.logger.Core().Enabled(zap.WarnLevel) {
		logStr, a := getArgs(args...)
		z.logger.Warn(fmt.Sprintf(logStr, a...))
	}
}

func (z ZapAdapter) Error(args ...interface{}) {
	logStr, a := getArgs(args...)
	z.logger.Error(fmt.Sprintf(logStr, a...))
}

func (z ZapAdapter) Panic(args ...interface{}) {
	logStr, a := getArgs(args...)
	z.logger.Panic(fmt.Sprintf(logStr, a...))
}

func (z ZapAdapter) Fatal(args ...interface{}) {
	logStr, a := getArgs(args...)
	z.logger.Fatal(fmt.Sprintf(logStr, a...))
}
