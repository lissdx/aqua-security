package validator

import (
	"errors"
	"fmt"
	"github.com/lissdx/aqua-security/internal/converting"
	asLogger "github.com/lissdx/aqua-security/internal/pkg/logger"
	strUtil "github.com/lissdx/aqua-security/internal/pkg/utils"
	"github.com/spf13/viper"
	jschema "github.com/xeipuuv/gojsonschema"
	"strings"
)

const ValidatorMapParamName = "VALIDATOR_MAP"

type JsonSchemaValidator struct {
	validatorFnMap   ValidatorFnMap
	validatorGoFnMap ValidatorGoFnMap
	logger           asLogger.TOLogger
}

func (jsv *JsonSchemaValidator) GetGoValidatorFn(validatorKey string) (ValidatorGoFn, error) {
	if vFunc, ok := jsv.validatorGoFnMap[normalizeKey(validatorKey)]; ok {
		return vFunc, nil
	}
	err := fmt.Errorf("validator error. provided key:[%s] have not matching to proper ValidatorFn", validatorKey)
	jsv.logger.Error(err.Error())
	return nil, err
}

func (jsv *JsonSchemaValidator) GetValidatorFn(validatorKey string) (ValidatorFn, error) {
	if vFunc, ok := jsv.validatorFnMap[normalizeKey(validatorKey)]; ok {
		return vFunc, nil
	}
	err := fmt.Errorf("validator error. provided key:[%s] have not matching to proper ValidatorFn", validatorKey)
	jsv.logger.Error(err.Error())
	return nil, err
}

func NewValidator(vConf *viper.Viper, logger asLogger.TOLogger) ValidatorFactory {
	jsValidator := JsonSchemaValidator{
		logger:         logger,
		validatorFnMap: ValidatorFnMap{},
	}

	valMapStr := vConf.GetString(ValidatorMapParamName)
	if !strUtil.IsEmptyString(valMapStr) {

		strMap, err := converting.StrToMap(strings.Trim(strings.TrimSpace(valMapStr), "[]"))
		if err != nil {
			logger.Panic("cannot create JsonSchema validator. %s: %s, error: %s", ValidatorMapParamName, valMapStr, err.Error())
		}

		for k, vFile := range strMap {
			// create JSON schema
			err = jsValidator.Append(k, vFile)
			if err != nil {
				logger.Panic("cannot create JsonSchema validator. %s: %s, error: %s", ValidatorMapParamName, valMapStr, err.Error())
			}
		}
	}

	return &jsValidator
}

func NewGoValidator(vConf *viper.Viper, logger asLogger.TOLogger) ValidatorGoFactory {
	jsValidator := JsonSchemaValidator{
		logger:           logger,
		validatorGoFnMap: ValidatorGoFnMap{},
	}

	valMapStr := vConf.GetString(ValidatorMapParamName)
	if !strUtil.IsEmptyString(valMapStr) {

		strMap, err := converting.StrToMap(strings.Trim(strings.TrimSpace(valMapStr), "[]"))
		if err != nil {
			logger.Panic("cannot create JsonSchema validator. %s: %s, error: %s", ValidatorMapParamName, valMapStr, err.Error())
		}

		for k, vFile := range strMap {
			// create JSON schema
			err = jsValidator.AppendGo(k, vFile)
			if err != nil {
				logger.Panic("cannot create JsonSchema validator. %s: %s, error: %s", ValidatorMapParamName, valMapStr, err.Error())
			}
		}
	}

	return &jsValidator
}

func (jsv *JsonSchemaValidator) Append(keyStr string, jsonSchemaPath string) error {
	// init src jsonschema
	loader := jschema.NewReferenceLoader(jsonSchemaPath)
	jsSchema, err := jschema.NewSchema(loader)
	if err != nil {
		jsErr := fmt.Errorf("can't create json schema: %s", err.Error())
		jsv.logger.Error(jsErr.Error())
		return jsErr
	}

	// init validation function
	var validatorFn ValidatorFn = func(data string) (bool, error) {
		strLoader := jschema.NewStringLoader(data)
		result, jsErr := jsSchema.Validate(strLoader)
		if jsErr != nil {
			jsv.logger.Error("schema validation error: %s", jsErr.Error())
			return false, jsErr
		}
		if !result.Valid() {
			jsv.logger.Error("the document is not valid. see errors :\n")
			errRes := errors.New("the document is not valid. see errors :\n")

			for _, rErr := range result.Errors() {
				// Err implements the ResultError interface
				errRes = fmt.Errorf("%w; %s", errRes, rErr)
				jsv.logger.Error("- %s", rErr)
			}
			return false, errRes
		}

		return true, nil
	}

	// update map
	jsv.validatorFnMap[normalizeKey(keyStr)] = validatorFn

	return nil
}

func (jsv *JsonSchemaValidator) AppendGo(keyStr string, jsonSchemaPath string) error {
	// init src jsonschema
	loader := jschema.NewReferenceLoader(jsonSchemaPath)
	jsSchema, err := jschema.NewSchema(loader)
	if err != nil {
		jsErr := fmt.Errorf("can't create json schema: %s", err.Error())
		jsv.logger.Error(jsErr.Error())
		return jsErr
	}

	// init validation function
	var validatorGoFn ValidatorGoFn = func(data interface{}) (bool, error) {
		goLoader := jschema.NewGoLoader(data)
		result, jsErr := jsSchema.Validate(goLoader)
		if jsErr != nil {
			jsv.logger.Error("schema validation error: %s", jsErr.Error())
			return false, jsErr
		}
		if !result.Valid() {
			jsv.logger.Error("the document is not valid. see errors :\n")
			errRes := errors.New("the document is not valid. see errors :\n")

			for _, rErr := range result.Errors() {
				// Err implements the ResultError interface
				errRes = fmt.Errorf("%w; %s", errRes, rErr)
				jsv.logger.Error("- %s", rErr)
			}
			return false, errRes
		}

		return true, nil
	}

	// update map
	jsv.validatorGoFnMap[normalizeKey(keyStr)] = validatorGoFn

	return nil
}

func normalizeKey(keyStr string) string {
	return strUtil.NormalizeStringToLower(keyStr)
}
