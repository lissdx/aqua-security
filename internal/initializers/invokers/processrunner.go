package invokers

import (
	"context"
	"fmt"
	asLogger "github.com/lissdx/aqua-security/internal/pkg/logger"
	"github.com/lissdx/aqua-security/pkg/processor"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const version = "07042024"

func ProcessRunner(newProcessor *processor.Processor, logger asLogger.TOLogger, config *viper.Viper, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {

			logger.Info(fmt.Sprintf("version: %s Run ProcessRunner service: %s, process: %s", version, config.GetString("SERVICE_NAME"), config.GetString("PROCESS_NAME")))
			go func(process *processor.Processor) {
				(*process).Run()
			}(newProcessor)

			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Stop ProcessRunner")
			func(process *processor.Processor) {
				(*process).Stop()
				logger.Info("ProcessRunner Successfully Stopped")
			}(newProcessor)
			return nil
		},
	})
}
