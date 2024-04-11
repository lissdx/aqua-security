package data_updater_process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lissdx/aqua-security/internal/drivers"
	"github.com/lissdx/aqua-security/internal/pkg/consumer"
	"github.com/lissdx/aqua-security/internal/pkg/consumer/fswatchdog"
	dbModel "github.com/lissdx/aqua-security/internal/pkg/gen/db"
	"github.com/lissdx/aqua-security/internal/pkg/log_handler"
	asLogger "github.com/lissdx/aqua-security/internal/pkg/logger"
	strUtil "github.com/lissdx/aqua-security/internal/pkg/utils"
	"github.com/lissdx/aqua-security/internal/validator"
	"github.com/lissdx/aqua-security/pkg/processor"
	"github.com/lissdx/yapgo/pkg/pipeline"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultSystemFilepathDelimiter = "/"
const systemDependedFilepathDelimiterParamName = "SYSTEM_DEPENDED_FILEPATH_DELIMITER"
const outputDirPathParamName = "OUTPUT_DIR_PATH"

const (
	inputConnectionsPref = "input/connections"
	inputResourcesPref   = "input/resources"
	inputScansPref       = "input/scans"
)

const (
	outputImageFilePref      = "image"
	outputRepositoryFilePref = "repository"
)
const (
	inputRepositoryKey = "input.repository"
	inputImageKey      = "input.image"
	inputScanKey       = "input.scan"
	inputConnectionKey = "input.connection"
	unknownKey         = "unknown"
)

type DataUpdaterProcess struct {
	dataSource              consumer.Consumer
	logger                  asLogger.TOLogger
	isActive                bool
	maxLogCnt               int
	processName             processor.Name
	mu                      sync.Mutex
	systemFilepathDelimiter string
	validator               *validator.ValidatorFactory
	store                   drivers.Store
	outputDirPath           string
}

func (dap *DataUpdaterProcess) Run() {

	dap.logger.Info("%v process starting... Run()", dap.Name())

	done := make(chan interface{}) // create control channel
	defer close(done)
	dap.isActive = true
	dap.dataSource.Run()

	eventProcessPipeLine := pipeline.New()
	var dataStream <-chan interface{} = dap.dataSource.ConsumerStream()

	stg1Name := "dataSourceEventHandler"
	eventProcessPipeLine.AddStage(dap.dataSourceEventHandler(stg1Name), dap.errorHandler(stg1Name))
	stg2Name := "filerFilterByPath"
	eventProcessPipeLine.AddFilterStage(dap.filerFilterByPath(stg2Name))
	stg3Name := "dbWriter"
	eventProcessPipeLine.AddStage(dap.dbWriter(stg3Name), dap.errorHandler(stg3Name))
	stg4Name := "filerNotScanFiles"
	eventProcessPipeLine.AddFilterStage(dap.filerNotScanFiles(stg4Name))
	stg5Name := "triggeredReport"
	eventProcessPipeLine.AddStage(dap.triggeredReport(stg5Name), dap.errorHandler(stg5Name))

	doneProcess := eventProcessPipeLine.RunPlug(done, dataStream)

	<-doneProcess
	dap.logger.Info("%s process end Run()", string(dap.Name()))
}

// Stg. part ----------------------------------------------

// dataSourceEventHandler - first stage.
// In this case just debug notification about the event handling
func (dap *DataUpdaterProcess) dataSourceEventHandler(stgName string) pipeline.ProcessFn {

	var debugHandler = dap.debugHandler(stgName)

	return func(inObj interface{}) (outObj interface{}, err error) {
		watchDogEvent, ok := inObj.(fsnotify.Event)

		if !ok {
			return inObj, fmt.Errorf("%+v: not a fsnotify.Event", inObj)
		}

		debugHandler(fmt.Sprintf("received event: %+v", inObj))
		return watchDogEvent, nil
	}
}

// dbWriter update DB
func (dap *DataUpdaterProcess) dbWriter(stgName string) pipeline.ProcessFn {

	var debugHandler = dap.debugHandler(stgName)

	return func(inObj interface{}) (outObj interface{}, err error) {
		debugHandler(fmt.Sprintf("received event: %+v", inObj))
		watchDogEvent := inObj.(fsnotify.Event)

		// read updated/created file
		jFile, err := os.Open(watchDogEvent.Name)
		defer jFile.Close()
		if err != nil {
			return inObj, err
		}

		decoder := json.NewDecoder(jFile)

		if _, err = decoder.Token(); err != nil {
			return inObj, fmt.Errorf("error on decoder.Token() open delimiter: %w", err)
		}

		// Start parse the file
		i := 0
		ignored := 0
		for decoder.More() {
			var unstructuredMap map[string]interface{}
			if err = decoder.Decode(&unstructuredMap); err != nil {
				return inObj, fmt.Errorf("error on decoder.More(). line: %d, error %w", i, err)
			}

			// get the key
			// we will use the key for validation and DB update
			key := dap.getKey(watchDogEvent.Name, unstructuredMap)
			if key == unknownKey {
				dap.logger.Error(fmt.Errorf("cannot match a valid key for object (just ignore it): %+v. file: %s", unstructuredMap, watchDogEvent.Name).Error())
				ignored++
				continue
			}

			// get JSON validator
			valFn, valError := (*dap.validator).GetValidatorFn(key)
			if valError != nil {
				dap.logger.Error(fmt.Errorf("GetValidatorFn for event (just ignore it): %+v. error: %w", unstructuredMap, valError).Error())
				ignored++
				continue
			}

			// validate the object
			d, _ := json.Marshal(unstructuredMap)
			isValid, valError := valFn(string(d))
			if valError != nil || !isValid {
				if valError == nil {
					valError = fmt.Errorf("validation error")
				}
				dap.logger.Error(fmt.Errorf("schema validation error: %w\n object (just ignore it): %s", valError, string(d)).Error())
				ignored++
				continue
			}

			debugHandler(fmt.Sprintf("object is valid: %+v try to insert it into DB", unstructuredMap))

			// insert into DB
			iErr := dap.insert(key, unstructuredMap)
			if iErr != nil {
				dap.logger.Error(fmt.Errorf("error on insert: %w", iErr).Error())
				ignored++
				continue
			}
			i++
		}

		dap.logger.Info("Update DB report: inserted/updated: %d, ignored: %d", i, ignored)
		// Read closing delimiter. `]` or `}`
		if _, err = decoder.Token(); err != nil {
			return inObj, fmt.Errorf("error on decoder.Token() closing delimiter: %w", err)
		}

		return inObj, nil
	}
}

// triggeredReport this stage is triggered by scans[*].json file
// it generates a report and puts it into output folder
func (dap *DataUpdaterProcess) triggeredReport(stgName string) pipeline.ProcessFn {

	var debugHandler = dap.debugHandler(stgName)

	return func(inObj interface{}) (outObj interface{}, err error) {

		debugHandler("Create Scan Report")
		currenTimestamp := time.Now().UTC().Unix()
		// Create repository report
		repositoryScanReport, err := dap.store.GetUnreportedRepositoryList(context.Background())
		if err != nil {
			return inObj, fmt.Errorf("on create repository scan report %w", err)
		}

		data, _ := json.MarshalIndent(repositoryScanReport, "", "  ")
		outRepositoryFileName := fmt.Sprintf("%s/%s_%d.json", dap.outputDirPath, outputRepositoryFilePref, currenTimestamp)
		orf, err := os.Create(outRepositoryFileName)
		if err != nil {
			return inObj, fmt.Errorf("on create report file: %s %w", outRepositoryFileName, err)
		}
		defer orf.Close()

		_, err = orf.Write(data)
		if err != nil {
			return inObj, fmt.Errorf("on write report file %s %w", outRepositoryFileName, err)
		}

		// Create image report
		imageScanReport, err := dap.store.GetUnreportedImageList(context.Background())
		if err != nil {
			return inObj, fmt.Errorf("on create image scan report %w", err)
		}

		data, _ = json.MarshalIndent(imageScanReport, "", "  ")
		outImageFileName := fmt.Sprintf("%s/%s_%d.json", dap.outputDirPath, outputImageFilePref, currenTimestamp)
		oif, err := os.Create(outImageFileName)
		if err != nil {
			return inObj, fmt.Errorf("on create report file: %s %w", outImageFileName, err)
		}
		defer oif.Close()

		_, err = oif.Write(data)
		if err != nil {
			return inObj, fmt.Errorf("on write report file %s %w", outImageFileName, err)
		}

		// Update reported scans
		// set is_reported == true
		scanIdList := make([]int32, 0)
		for _, v := range repositoryScanReport {
			scanIdList = append(scanIdList, v.ScanID)
		}
		for _, v := range imageScanReport {
			scanIdList = append(scanIdList, v.ScanID)
		}

		_, err = dap.store.SetBulkUpdateScanReportTrue(context.Background(), scanIdList)
		if err != nil {
			return inObj, fmt.Errorf("on bulk scan report %w", err)
		}

		debugHandler("scan report successfully generated")
		debugHandler(fmt.Sprintf("repository report: %s", outRepositoryFileName))
		debugHandler(fmt.Sprintf("image report: %s", outImageFileName))
		return inObj, nil
	}
}

// Filter. part ----------------------------------------------
// filerFilterByPath - fitter all events than not related to our flow
// the file path should be:
// input/[connections*.json|resources*.json|scans*.json]
func (dap *DataUpdaterProcess) filerFilterByPath(stgName string) pipeline.FilterFn {

	// TODO: move to config
	allowedPrefixList := []string{
		inputConnectionsPref,
		inputResourcesPref,
		inputScansPref,
	}
	// TODO: move to config
	allowedFileExtension := ".json"

	var debugHandler = dap.debugHandler(stgName)
	var errorHandler = dap.errorHandler(stgName)

	return func(inObj interface{}) (outObj interface{}, result bool) {

		watchDogEvent := inObj.(fsnotify.Event)
		filePath := watchDogEvent.Name

		if !strings.HasSuffix(filePath, allowedFileExtension) {
			errorHandler(fmt.Errorf("file extention is allowed. file path: %s", filePath))
			debugHandler(fmt.Sprintf("event is filtered. reason: file extention is allowed. file path: %s", filePath))
			return inObj, false
		}

		splitPath := strings.Split(filePath, dap.systemFilepathDelimiter)
		if len(splitPath) < 2 {
			errorHandler(fmt.Errorf("file path is not allowed (short). file path: %s", filePath))
			debugHandler(fmt.Sprintf("event is filtered (short). Reasonfile path is not allowed. file path: %s", filePath))
			return inObj, false
		}

		isAllowed := false
		subStr := strings.Join(splitPath[(len(splitPath)-2):], dap.systemFilepathDelimiter)
		for i := 0; i < len(allowedPrefixList) && !isAllowed; i++ {
			isAllowed = strings.HasPrefix(subStr, allowedPrefixList[i])
		}

		if !isAllowed {
			errorHandler(fmt.Errorf("file path is not allowed. file path: %s", filePath))
			debugHandler(fmt.Sprintf("event is filtered. Reasonfile path is not allowed. file path: %s", filePath))
			return inObj, false
		}

		debugHandler(fmt.Sprintf("status: OK, event: %s", watchDogEvent))

		return watchDogEvent, true
	}
}

// filerNotScanFiles is filtering input/[connections*.json|resources*.json]
// but input/scans*.json is allowed
func (dap *DataUpdaterProcess) filerNotScanFiles(stgName string) pipeline.FilterFn {

	var debugHandler = dap.debugHandler(stgName)

	return func(inObj interface{}) (outObj interface{}, result bool) {

		watchDogEvent := inObj.(fsnotify.Event)
		filePath := watchDogEvent.Name

		splitPath := strings.Split(filePath, dap.systemFilepathDelimiter)
		isAllowed := false
		subStr := strings.Join(splitPath[(len(splitPath)-2):], dap.systemFilepathDelimiter)
		isAllowed = strings.HasPrefix(subStr, inputScansPref)

		if !isAllowed {
			debugHandler(fmt.Sprintf("event is filtered. Reason is not a scan-file. file path: %s", filePath))
			return inObj, false
		}

		debugHandler(fmt.Sprintf("... is scan file, let's triger the REPORT: file: %s ", filePath))

		return watchDogEvent, true
	}
}

// Process part --------------------------------------

func (dap *DataUpdaterProcess) Stop() {
	defer dap.mu.Unlock()
	dap.mu.Lock()
	if dap.isActive {
		dap.isActive = false
		dap.dataSource.Stop()
		dap.logger.Info("%s process STOP", string(dap.processName))
	} else {
		dap.logger.Info("%s process already not active", string(dap.processName))
	}
}

func (dap *DataUpdaterProcess) Name() processor.Name {
	return dap.processName
}

// Error and Debug handlers part --------------------------------------
func (dap *DataUpdaterProcess) errorHandler(stgName string) pipeline.ErrorProcessFn {
	opt := []log_handler.Option{log_handler.WithProcessName(string(dap.Name())),
		log_handler.WithStageName(stgName),
		log_handler.WithLogger(dap.logger),
	}

	return log_handler.ErrorHandlerFnFactory(opt...)
}

func (dap *DataUpdaterProcess) debugHandler(stgName string) func(string) {
	opt := []log_handler.Option{log_handler.WithProcessName(string(dap.Name())),
		log_handler.WithStageName(stgName),
		log_handler.WithLogger(dap.logger),
	}

	return log_handler.DebugHandlerFnFactory(opt...)
}

// Constructor ------------------------------

func NewDataUpdaterProcess(config *viper.Viper, logger asLogger.TOLogger, store drivers.Store) processor.Processor {

	strUtil.CoalesceStr(config.GetString(systemDependedFilepathDelimiterParamName), defaultSystemFilepathDelimiter)
	validatorFactory := validator.NewValidator(config, logger)

	outputDirPath := config.GetString(outputDirPathParamName)
	if strUtil.IsEmptyString(outputDirPath) {
		panic(fmt.Sprintf("%s param is mandatory", outputDirPathParamName))
	}
	return &DataUpdaterProcess{
		dataSource:              fswatchdog.NewFsWatchdog(config, logger),
		logger:                  logger,
		processName:             processor.Name(config.GetString("PROCESS_NAME")),
		systemFilepathDelimiter: strUtil.CoalesceStr(config.GetString(systemDependedFilepathDelimiterParamName), defaultSystemFilepathDelimiter),
		validator:               &validatorFactory,
		store:                   store,
		outputDirPath:           outputDirPath,
	}
}

// Helpers ---------------------------------
func (dap *DataUpdaterProcess) getKey(path string, m map[string]interface{}) string {
	splitPath := strings.Split(path, dap.systemFilepathDelimiter)

	subStr := strings.Join(splitPath[(len(splitPath)-2):], dap.systemFilepathDelimiter)
	if strings.HasPrefix(subStr, inputResourcesPref) {
		if v, ok := m["type"]; !ok {
			return unknownKey
		} else {
			switch v.(string) {
			case "repository":
				return inputRepositoryKey
			case "image":
				return inputImageKey
			default:
				return unknownKey
			}
		}
	}

	if strings.HasPrefix(subStr, inputConnectionsPref) {
		return inputConnectionKey
	}

	if strings.HasPrefix(subStr, inputScansPref) {
		return inputScanKey
	}

	return unknownKey
}

// insert write data into db
func (dap *DataUpdaterProcess) insert(key string, unstructuredMap map[string]interface{}) error {
	switch key {
	case inputRepositoryKey:
		insertParams, cErr := convertToInputRepositoryDbObject(unstructuredMap)
		if cErr != nil {
			return cErr
		}
		iRes, dbErr := dap.store.InsertInputRepository(context.Background(), insertParams)
		if dbErr != nil {
			return dbErr
		}
		dap.logger.Debug(fmt.Sprintf("input_resource table upserted result: %+v", iRes))
	case inputImageKey:
		insertParams, cErr := convertToInputImageDbObject(unstructuredMap)
		if cErr != nil {
			return cErr
		}
		iRes, dbErr := dap.store.InsertInputImage(context.Background(), insertParams)
		if dbErr != nil {
			return dbErr
		}
		dap.logger.Debug(fmt.Sprintf("input_resource table upserted result: %+v", iRes))
	case inputConnectionKey:
		insertParams, cErr := convertToInputConnectionDbObject(unstructuredMap)
		if cErr != nil {
			return cErr
		}

		dbErr := dap.store.UpdateConnection(context.Background(), insertParams)
		if dbErr != nil {
			return dbErr
		}
		dap.logger.Debug(fmt.Sprintf("input_image table updated by result: %+v", unstructuredMap))
	case inputScanKey:
		insertParams, cErr := convertToInputScanDbObject(unstructuredMap)
		if cErr != nil {
			return cErr
		}
		iRes, dbErr := dap.store.InsertInputScan(context.Background(), insertParams)
		if dbErr != nil {
			return dbErr
		}
		dap.logger.Debug(fmt.Sprintf("scan table upserted result: %+v", iRes))

	default:
		return fmt.Errorf("key %s is not supported", key)
	}

	return nil

}

// convertors prepare/create DB structures
func convertToInputRepositoryDbObject(m map[string]interface{}) (res dbModel.InsertInputRepositoryParams, err error) {

	err = res.ID.Scan(m["id"].(string))
	if err != nil {
		return
	}

	res.Name = m["name"].(string)

	res.Url = m["url"].(string)

	err = res.Source.Scan(m["source"].(string))
	if err != nil {
		return
	}

	res.CreatedDateTimestamp = time.Unix(int64(m["created_date_timestamp"].(float64)), 0).UTC()

	if t, pErr := time.Parse(time.RFC3339, string(m["last_push"].(string))); pErr != nil {
		err = pErr
		return
	} else {
		res.LastPush = t
	}

	res.Size = int64(m["size"].(float64))

	return
}
func convertToInputImageDbObject(m map[string]interface{}) (res dbModel.InsertInputImageParams, err error) {

	err = res.ID.Scan(m["id"].(string))
	if err != nil {
		return
	}

	res.Name = m["name"].(string)

	res.Url = m["url"].(string)

	err = res.Source.Scan(m["source"].(string))
	if err != nil {
		return
	}

	res.CreatedDateTimestamp = time.Unix(int64(m["created_date_timestamp"].(float64)), 0).UTC()

	res.NumberOfLayers = int32(m["number_of_layers"].(float64))

	err = res.Architecture.Scan(m["architecture"].(string))
	if err != nil {
		return
	}

	return
}
func convertToInputConnectionDbObject(m map[string]interface{}) (res dbModel.UpdateConnectionParams, err error) {

	err = res.RepositoryID.Scan(m["repository_id"].(string))
	if err != nil {
		return
	}

	err = res.ID.Scan(m["image_id"].(string))

	return
}
func convertToInputScanDbObject(m map[string]interface{}) (res dbModel.InsertInputScanParams, err error) {

	res.ScanID = int32(m["scan_id"].(float64))

	err = res.ResourceID.Scan(m["resource_id"].(string))
	if err != nil {
		return
	}

	err = res.ResourceType.Scan(m["resource_type"].(string))
	if err != nil {
		return
	}

	err = res.HighestSeverity.Scan(m["highest_severity"].(string))
	if err != nil {
		return
	}

	res.TotalFindings = int32(m["total_findings"].(float64))

	res.ScanDateTimestamp = time.Unix(int64(m["scan_date_timestamp"].(float64)), 0).UTC()

	return
}
