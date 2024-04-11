package providers

import (
	"fmt"
	"github.com/lissdx/aqua-security/internal/drivers"
	asLogger "github.com/lissdx/aqua-security/internal/pkg/logger"
	"github.com/lissdx/aqua-security/internal/pkg/utils"
	"github.com/lissdx/aqua-security/internal/processes/data_updater_process"
	"github.com/lissdx/aqua-security/pkg/processor"
	"github.com/spf13/viper"
	"strings"
)

func NewProcessor(config *viper.Viper, logger asLogger.TOLogger, store drivers.Store) *processor.Processor {
	processName := utils.NormalizeStringToUpper(config.GetString(`PROCESS_NAME`))

	newProcess := processorFactory(strings.TrimSpace(processName), config, logger, store)
	if newProcess == nil {
		panic(fmt.Sprintf("can't start new process: %s", processName))
	}

	return &newProcess
}

func processorFactory(processName string,
	config *viper.Viper,
	logger asLogger.TOLogger,
	store drivers.Store) processor.Processor {

	switch strings.ToUpper(processName) {
	case "DATA_UPDATER":
		return data_updater_process.NewDataUpdaterProcess(config, logger, store)
	default:
		logger.Error(fmt.Sprintf("Can't find and create process: %s", processName))
		return nil
	}
}
