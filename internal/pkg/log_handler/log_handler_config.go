package log_handler

import (
	asLogger "github.com/lissdx/aqua-security/internal/pkg/logger"
)

type config struct {
	StageName   string
	ProcessName string
	Logger      asLogger.TOLogger
}

// Option interface used for setting optional config properties.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (fn optionFunc) apply(c *config) {
	fn(c)
}

func WithStageName(stageName string) Option {
	return optionFunc(func(cfg *config) {
		cfg.StageName = stageName
	})
}

func WithProcessName(processName string) Option {
	return optionFunc(func(cfg *config) {
		cfg.ProcessName = processName
	})
}

func WithLogger(logger asLogger.TOLogger) Option {
	return optionFunc(func(cfg *config) {
		cfg.Logger = logger
	})
}
