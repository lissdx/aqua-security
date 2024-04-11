package config

import (
	strUtil "github.com/lissdx/aqua-security/internal/pkg/utils"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const portEnvName = "PORT"
const hostEnvName = "HOST"
const GoEnvEnvName = "GO_ENV"
const LogLevelEnvName = "LOG_LEVEL"

func NewConfig() *viper.Viper {

	env := strings.TrimSpace(strings.ToLower(os.Getenv(GoEnvEnvName)))

	if strUtil.IsEmptyString(env) {
		env = "development_local"
		err := os.Setenv(GoEnvEnvName, env)
		if err != nil {
			panic(err)
		}
	}

	config := viper.New()
	config.AutomaticEnv()
	config.SetConfigName(env)
	config.AddConfigPath("./configs/monitoring-process")

	err := config.ReadInConfig()

	if err != nil {
		panic(err)
	}

	// for local usage
	// actually on pod we have an ability to LOG_LEVEL from env
	if config.IsSet(LogLevelEnvName) {
		logLevel := config.GetString(LogLevelEnvName)
		err = os.Setenv(LogLevelEnvName, logLevel)
		if err != nil {
			panic(err)
		}
	}

	if config.IsSet(portEnvName) &&
		strUtil.IsEmptyString(os.Getenv(portEnvName)) {
		err = os.Setenv(portEnvName, config.GetString(portEnvName))
		if err != nil {
			panic(err)
		}
	}

	if config.IsSet(hostEnvName) &&
		strUtil.IsEmptyString(os.Getenv(hostEnvName)) {
		err = os.Setenv(hostEnvName, config.GetString(hostEnvName))
		if err != nil {
			panic(err)
		}
	}

	return config
}
