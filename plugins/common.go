package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

const (
	CHARGE_TYPE_PREPAID = "PREPAID"
)

type CallBackParameter struct {
	Parameter string `json:"callbackParameter,omitempty"`
}

type Filter struct {
	Name   string
	Values []string
}

func TransLittleCamelcaseToShortLineFormat(inputValue string) (string, error) {
	str := ""
	for i := 0; i < len(inputValue); i++ {

		ch := inputValue[i]
		if ch < 65 || (ch > 90 && ch < 97) || ch > 126 {
			return str, fmt.Errorf("wrong character")
		}
		if ch < 'a' {
			str = fmt.Sprintf("%s-%c", str, ch+32)
		} else {
			str = fmt.Sprintf("%s%c", str, ch)
		}
	}
	return str, nil
}

func IsValidValue(inputValue string, validValues []string) error {
	for _, validValue := range validValues {
		if validValue == inputValue {
			return nil
		}
	}
	return fmt.Errorf("%s is not valid value in(%++v)", inputValue, validValues)
}

func GetRegionFromProviderParams(providerParams string) (string, error) {
	paramMap, err := GetMapFromProviderParams(providerParams)
	if err != nil {
		return "", err
	}

	region, ok := paramMap["Region"]
	if !ok {
		return region, fmt.Errorf("region not found in providerParams")
	}
	return region, nil
}

func GetMapFromProviderParams(providerParams string) (map[string]string, error) {
	rtnMap := make(map[string]string)
	params := strings.Split(providerParams, ";")

	if len(params) == 0 {
		return rtnMap, nil
	}

	for _, param := range params {
		afterTrimParam := strings.Trim(param, " ")
		kv := strings.Split(afterTrimParam, "=")
		if len(kv) == 2 {
			rtnMap[kv[0]] = kv[1]
		} else {
			return rtnMap, fmt.Errorf("GetMapFromProviderParams meet illegal format param=%s", param)
		}
	}
	return rtnMap, nil
}

type CommonInputs struct {
	Inputs []interface{} `json:"inputs,omitempty"`
}

func UnmarshalJson(source interface{}, target interface{}) error {
	reader, ok := source.(io.Reader)
	if !ok {
		return fmt.Errorf("the source to be unmarshaled is not a io.reader type")
	}

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("parse http request (%v) meet error (%v)", reader, err)
	}

	// commonInputs := CommonInputs{}
	// if err = json.Unmarshal(bodyBytes, &commonInputs); err != nil {
	// 	return fmt.Errorf("unmarshal http request (%v) meet error (%v)", reader, err)
	// }
	// if len(commonInputs.Inputs) == 0 {
	// 	return fmt.Errorf("empty inputs")
	// }

	if err = json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("unmarshal http request (%v) meet error (%v)", reader, err)
	}
	return nil
}

func ExtractJsonFromStruct(s interface{}) map[string]string {
	fields := make(map[string]string)
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i).Tag.Get("json")
			fields[strings.Split(field, ",")[0]] = t.Field(i).Type.String()
		}
	}
	return fields
}
