package providers

import (
	"fmt"
	strUtil "github.com/lissdx/aqua-security/internal/pkg/utils"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const GoEnvEnvName = "GO_ENV"
const ProcessNameEnvName = "PROCESS_NAME"
const ProcessConfSubFolderEnvName = "PROCESS_CONF_SUB_FOLDER"

func NewConfig() *viper.Viper {

	processName := strings.TrimSpace(strings.ToLower(os.Getenv(ProcessNameEnvName)))
	if strUtil.IsEmptyString(processName) {
		panic("unknown process name, please set PROCESS_NAME env")
	}

	var subFolderName = strings.TrimSpace(os.Getenv(ProcessConfSubFolderEnvName))
	if len(subFolderName) == 0 {
		subFolderName = processName
	}

	// set default GO_ENV if not existed
	env := strings.TrimSpace(strings.ToLower(os.Getenv(GoEnvEnvName)))
	if strUtil.IsEmptyString(env) {
		env = "development"
		err := os.Setenv(GoEnvEnvName, env)
		if err != nil {
			panic(err)
		}
	}

	config := viper.New()
	config.AutomaticEnv()
	config.SetConfigName(env)
	config.AddConfigPath(fmt.Sprintf("./configs/%s", subFolderName))

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// for local usage
	// actually on pod we have an ability to LOG_LEVEL from env
	if config.IsSet("LOG_LEVEL") {
		logLevel := config.GetString("LOG_LEVEL")
		err = os.Setenv("LOG_LEVEL", logLevel)
		if err != nil {
			panic(err)
		}
	}

	return config
}
