package log_handler

import (
	"fmt"
	strUtil "github.com/lissdx/aqua-security/internal/pkg/utils"
)

func DebugHandlerFnFactory(options ...Option) func(string) {
	cfg := config{}
	for _, opt := range options {
		opt.apply(&cfg)
	}

	if cfg.Logger == nil {
		panic("cannot create debug handle function. cfg.Logger is mandatory")
	}

	if strUtil.IsEmptyString(cfg.StageName) || strUtil.IsEmptyString(cfg.ProcessName) {
		panic("cannot create debug handle function. ProcessName and StageName are mandatory")
	}

	return func(dbg string) {
		debugStr := fmt.Sprintf("process: %v, stgName: %s | msg: %s", cfg.ProcessName, cfg.StageName, dbg)
		cfg.Logger.Debug(debugStr)
	}
}
