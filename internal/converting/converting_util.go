package converting

import (
	"fmt"
	strUtil "github.com/lissdx/aqua-security/internal/pkg/utils"
	"go.uber.org/zap/zapcore"
	"strconv"
	"strings"
)

const (
	defaultStrSeparator    = ","
	defaultStrMapSeparator = ":"
)

const (
	invalidInputValue = "%s can't convert the value: %v"
	convertingError   = "%s can't convert the value: %v, error: %s"
)

func StrToStr(toConvert string) (res string, err error) {

	if strUtil.IsEmptyString(toConvert) {
		err = fmt.Errorf(invalidInputValue, "StrToStr(toConvert string)", toConvert)
		return
	}

	return toConvert, nil
}

func StrToStrWithTrim(toConvert string) (res string, err error) {

	res, err = StrToStr(strings.TrimSpace(toConvert))

	if err != nil {
		res = ""
	}

	return
}

func StrToUint64(toConvert string) (res uint64, err error) {

	if strUtil.IsEmptyString(toConvert) {
		err = fmt.Errorf(invalidInputValue, "StrToUint64(toConvert string)", toConvert)
		return
	}

	converted, err := strconv.ParseUint(strings.TrimSpace(toConvert), 10, 64)
	if err != nil {
		err = fmt.Errorf(convertingError, "StrToUint64(toConvert string)", toConvert, err)
		return
	}

	return converted, nil
}

func StrToInt(toConvert string) (res int, err error) {

	if strUtil.IsEmptyString(toConvert) {
		err = fmt.Errorf(invalidInputValue, "StrToInt(toConvert string)", toConvert)
		return
	}

	converted, err := strconv.Atoi(strings.TrimSpace(toConvert))
	if err != nil {
		err = fmt.Errorf(convertingError, "StrToInt(toConvert string)", toConvert, err)
		return
	}

	return converted, nil
}

func StrToStrArray(toConvert string) (res []string, err error) {
	res, err = StrToStrArrayWithSeparator(toConvert, defaultStrSeparator)
	if err != nil {
		err = fmt.Errorf(invalidInputValue, "StrToStrArray(toConvert string)", toConvert)
	}

	return
}

func StrToStrArrayWithSeparator(toConvert string, separator string) (res []string, err error) {
	var converted = make([]string, 0)
	for _, subStr := range strings.Split(toConvert, separator) {
		if len(strings.TrimSpace(subStr)) != 0 {
			converted = append(converted, strings.TrimSpace(subStr))
		}
	}

	if len(converted) == 0 {
		err = fmt.Errorf(invalidInputValue, "StrToStrArrayWithSeparator(toConvert string, separator string)", toConvert)
	} else {
		res = converted
	}

	return
}

func StrToBool(toConvert string) (res bool, err error) {

	res, err = strconv.ParseBool(strings.TrimSpace(toConvert))

	if err != nil {
		err = fmt.Errorf(convertingError, "StrToBool(toConvert string)", toConvert, err)
	}

	return
}

// StrToMapBool converts coma separated string to map[string]bool
// where key should be the str and value should be true
// for example:
// res :=  StrToMapBool("a,b,c")
// res: map[string]bool{"a": true, "b": true, "c": true}
func StrToMapBool(toConvert string) (res map[string]bool, err error) {
	strArray, err := StrToStrArray(toConvert)
	if err != nil {
		err = fmt.Errorf(invalidInputValue, "StrToMapBool(toConvert string) error: ", err)
		return nil, err
	}
	res = make(map[string]bool)
	for _, key := range strArray {
		res[key] = true
	}
	return
}

func StrToZapCoreLevel(toConvert string) (res zapcore.Level, err error) {

	err = res.Set(strings.ToLower(strings.TrimSpace(toConvert)))

	return
}

// StrToMap converts column separated string to map[string]string
// where key should be the str and value should be true
// for example:
// res :=  StrToMap("a:true,b:valB,c:valC")
// res: map[string]string{"a": "true", "b": ""valB, "c": "valC"}
func StrToMap(toConvert string) (res map[string]string, err error) {

	// split env var into list
	// expected list of strings like "key1:val1, key2:val2"
	strList, err := StrToStrArray(toConvert)
	if err != nil {
		err = fmt.Errorf(convertingError, "StrToMap(toConvert string)", toConvert, err)
		return nil, err
	}

	res = map[string]string{}
	// split every tuple into key, value entry
	// update result map
	for _, keyVal := range strList {
		keyValTuple := strings.SplitN(keyVal, defaultStrMapSeparator, 2)
		if len(keyValTuple) != 2 {
			err = fmt.Errorf(convertingError, "StrToMap(toConvert string) invalid key:value pair", toConvert, err)
			return nil, err
		}
		key, errK := StrToStrWithTrim(keyValTuple[0])
		if errK != nil || strUtil.IsEmptyString(key) {
			err = fmt.Errorf(convertingError, "StrToMap(toConvert string) invalid key:value pair", toConvert, errK)
			return nil, err
		}
		val, errV := StrToStrWithTrim(keyValTuple[1])
		if errV != nil || strUtil.IsEmptyString(val) {
			err = fmt.Errorf(convertingError, "StrToMap(toConvert string) invalid key:value pair", toConvert, errV)
			return nil, err
		}
		res[key] = val
	}

	return res, nil
}
