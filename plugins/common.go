package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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

func UnmarshalJson(r *http.Request, object interface{}) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("parse http request (%v) meet error (%v)", r, err)
	}

	if err = json.Unmarshal(bodyBytes, object); err != nil {
		return fmt.Errorf("unmarshal http request (%v) meet error (%v)", r, err)
	}
	return nil
}
