package log_handler

import (
	"fmt"
	strUtil "github.com/lissdx/aqua-security/internal/pkg/utils"
	"github.com/lissdx/yapgo/pkg/pipeline"
)

func ErrorHandlerFnFactory(options ...Option) pipeline.ErrorProcessFn {
	cfg := config{}
	for _, opt := range options {
		opt.apply(&cfg)
	}

	if cfg.Logger == nil {
		panic("cannot create error handle function. cfg.Logger is mandatory")
	}

	if strUtil.IsEmptyString(cfg.StageName) || strUtil.IsEmptyString(cfg.ProcessName) {
		panic("cannot create error handle function. ProcessName and StageName are mandatory")
	}

	return func(err error) {
		errorStr := fmt.Sprintf("process: %v, stgName: %s | error: %s", cfg.ProcessName, cfg.StageName, err.Error())
		cfg.Logger.Error(errorStr)
	}
}
