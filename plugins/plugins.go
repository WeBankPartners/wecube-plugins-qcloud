package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	pluginsMutex sync.Mutex
	plugins      = make(map[string]Plugin)
)

type Plugin interface {
	GetActionByName(actionName string) (Action, error)
}

type Action interface {
	BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error)
	CheckParam(param interface{}) error
	Do(param interface{}, workflowParam *WorkflowParam) error
}

func registerPlugin(name string, plugin Plugin) {
	pluginsMutex.Lock()
	defer pluginsMutex.Unlock()

	if _, found := plugins[name]; found {
		logrus.Fatalf("cloud provider %q was registered twice", name)
	}

	plugins[name] = plugin
}

func getPluginByName(name string) (Plugin, error) {
	pluginsMutex.Lock()
	defer pluginsMutex.Unlock()
	plugin, found := plugins[name]
	if !found {
		return nil, fmt.Errorf("plugin[%s] not found", name)
	}
	return plugin, nil
}

func init() {
	registerPlugin("vm", new(VmPlugin))

}

type WorkflowParam struct {
	AckPath             string `json:"ackPath"`
	AckServer           string `json:"ackServer"`
	ApplicationName     string `json:"applicationName"`
	ApplicationAction   string `json:"applicationAction"`
	ProcessDefinitionID string `json:"processDefinitionId"`
	ProcessExecutionID  string `json:"processExecutionId"`
	ProcessInstanceID   string `json:"processInstanceId"`
	RequestID           string `json:"requestId"`

	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`

	PluginVersion string
	PluginName    string
	ProviderName  string
}

func CallPluginAction(workflowParam WorkflowParam) {
	var err error
	defer func() {
		if err != nil {
			logrus.Errorf("plguin[%v]-action[%v] meet error = %v", workflowParam.PluginName, workflowParam.ApplicationAction, err)
			workflowParam.ResultCode = "1"
			workflowParam.ResultMsg = fmt.Sprint(err)
		} else {
			logrus.Infof("plguin[%v]-action[%v] completed", workflowParam.PluginName, workflowParam.ApplicationAction)
			workflowParam.ResultCode = "0"
		}
		callbackWorkflow(&workflowParam)
	}()

	logrus.Infof("plguin[%v]-action[%v] start...", workflowParam.PluginName, workflowParam.ApplicationAction)

	plugin, err := getPluginByName(workflowParam.PluginName)
	if err != nil {
		return
	}

	action, err := plugin.GetActionByName(workflowParam.ApplicationAction)
	if err != nil {
		return
	}

	logrus.Infof("get CMDB parameters with process instance id = %v", workflowParam.ProcessInstanceID)
	actionParam, err := action.BuildParamFromCmdb(&workflowParam)
	if err != nil {
		return
	}
	logrus.Infof("CMDB parameters results = %v", actionParam)

	if err = action.CheckParam(actionParam); err != nil {
		return
	}
	logrus.Info("CMDB parameters are passed validation")

	logrus.Infof("action with parameters = %v", actionParam)
	if err = action.Do(actionParam, &workflowParam); err != nil {
		return
	}

	return
}

func callbackWorkflow(workflowParam *WorkflowParam) {
	requestBytes, err := json.Marshal(workflowParam)
	if err != nil {
		logrus.Errorf("callbackWorkflow Marshal failed err=%v", err)
		return
	}

	url := "http://" + workflowParam.AckServer + workflowParam.AckPath
	contentType := "application/json"
	logrus.Debugf("callbackWorkflow request url = %s,requestData = %s", url, string(requestBytes))

	response, err := http.Post(url, contentType, bytes.NewReader(requestBytes))
	if err != nil {
		logrus.Errorf("callbackWorkflow Post failed err=%v", err)
		return
	}

	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != 200 {
		logrus.Errorf("callbackWorkflow response.StatusCode != 200,statusCode=%v", response.StatusCode)
		return
	}

	logrus.Infof("callback workflow has been done, response is [%v]", string(bytes))
	return
}
