package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	SECRET_ID       = "Your Qcloud Secret Id"
	SECRET_KEY      = "Your Qcloud Secret Key"
	PLUGIN_HOST_URL = "http://10.107.117.154:8081"
)

var resourceIds = make(map[string]string)

type Outputs struct {
	Outputs []Output `json:"outputs,omitempty"`
}

type Output struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type PluginResponse struct {
	ResultCode string  `json:"result_code"`
	ResultMsg  string  `json:"result_message"`
	Results    Outputs `json:"results"`
}

func CallPlugin(name, action, input string) map[string]string {
	output, err := http.Post(PLUGIN_HOST_URL+"/v1/qcloud/"+name+"/"+action, "application/json", strings.NewReader(input))

	if err != nil {
		logrus.Errorf("call plugin server meet error = %v", err)
		//panic("failed")
	}

	pluginResponse := PluginResponse{}
	err = UnmarshalJson(output.Body, &pluginResponse)
	if err != nil {
		logrus.Errorf("unmarshal plugin response meet error = %v", err)
		//panic("failed")
	}

	if pluginResponse.ResultCode == "1" {
		logrus.Errorf("call plugin meet error = %v", pluginResponse.ResultMsg)
		//panic("failed")
	}

	outputMap := make(map[string]string)
	for i := 0; i < len(pluginResponse.Results.Outputs); i++ {
		outputMap[pluginResponse.Results.Outputs[i].Guid] = pluginResponse.Results.Outputs[i].Id
	}
	logrus.Infof("resource (ids=%v) have been handled, plugin = %v, action = %v", outputMap, name, action)

	return outputMap
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
