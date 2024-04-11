package providers

import (
	"fmt"
	asLogger "github.com/lissdx/aqua-security/internal/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

func NewLogger(config *viper.Viper) asLogger.TOLogger {
	var loggerConfig zap.Config

	if strings.HasPrefix(config.GetString("GO_ENV"), "development") {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	strLogLevel := strings.TrimSpace(config.GetString("LOG_LEVEL"))
	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(strLogLevel)); err != nil {
		return nil
	}

	loggerConfig = zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	}

	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.ErrorOutputPaths = []string{"stderr"}
	res, err := loggerConfig.Build()

	if err != nil {
		panic(fmt.Sprintf("can't create logger for serice: %s", config.GetString("SERVICE_NAME")))
	}

	return asLogger.NewZapAdapter(res)
}
