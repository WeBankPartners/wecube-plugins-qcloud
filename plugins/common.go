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

type Filter struct{
	Name string   
	Values []string{}
}

func isValidValue(inputValue string,validValues []string)error{
	for _,validValue:=range validValues {
		if validValue == inputValue {
			return nil 
		}
	}
	return fmt.Errorf("%s is not valid value in(%++v)",inputValue,validValues)
}

func transferStringArrayToStringPointerArray(stringArray []string)[]*string{
	outputs:=[]*string{}
	for _,value:=stringArray{
		temp:=value
		outputs:=append(outputs,&temp)
	}
	return outputs
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

func UnmarshalJson(source interface{}, target interface{}) error {
	reader, ok := source.(io.Reader)
	if !ok {
		return fmt.Errorf("the source to be unmarshaled is not a io.reader type")
	}

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("parse http request (%v) meet error (%v)", reader, err)
	}

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
