package plugins

import (
	"fmt"
	"net/http"
	"strings"
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
	ReadParam(*http.Request) (interface{}, error)
	CheckParam(param interface{}) error
	Do(input interface{}) (interface{}, error)
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
	registerPlugin("storage", new(StoragePlugin))
	registerPlugin("security-group", new(SecurityGroupPlugin))
	registerPlugin("subnet", new(SubnetPlugin))
	registerPlugin("nat-gateway", new(NatGatewayPlugin))
	registerPlugin("vpc", new(VpcPlugin))
	registerPlugin("peering-connection", new(PeeringConnectionPlugin))
}

type PluginRequest struct {
	Version      string
	ProviderName string
	Name         string
	Action       string
}

type PluginResponse struct {
	ResultCode string      `json:"result_code"`
	ResultMsg  string      `json:"result_message"`
	Results    interface{} `json:"results"`
}

func CallPluginAction(r *http.Request) (*PluginResponse, error) {
	var pluginResponse = PluginResponse{}
	pluginRequest := parsePluginRequest(r)
	var err error
	defer func() {
		if err != nil {
			logrus.Errorf("plguin[%v]-action[%v] meet error = %v", pluginRequest.Name, pluginRequest.Action, err)
			pluginResponse.ResultCode = "1"
			pluginResponse.ResultMsg = fmt.Sprint(err)
		} else {
			logrus.Infof("plguin[%v]-action[%v] completed", pluginRequest.Name, pluginRequest.Action)
			pluginResponse.ResultCode = "0"
		}
	}()

	logrus.Infof("plguin[%v]-action[%v] start...", pluginRequest.Name, pluginRequest.Action)

	plugin, err := getPluginByName(pluginRequest.Name)
	if err != nil {
		return &pluginResponse, err
	}

	action, err := plugin.GetActionByName(pluginRequest.Action)
	if err != nil {
		return &pluginResponse, err
	}

	logrus.Infof("read parameters from http request = %v", r)
	actionParam, err := action.ReadParam(r)
	if err != nil {
		return &pluginResponse, err
	}

	logrus.Infof("check parameters = %v", actionParam)
	if err = action.CheckParam(actionParam); err != nil {
		return &pluginResponse, err
	}

	logrus.Infof("action do with parameters = %v", actionParam)
	outputs, err := action.Do(actionParam)
	if err != nil {
		return &pluginResponse, err
	}

	pluginResponse.Results = outputs

	return &pluginResponse, nil
}

func parsePluginRequest(r *http.Request) *PluginRequest {
	var pluginInput = PluginRequest{}
	pathStrings := strings.Split(r.URL.Path, "/")
	logrus.Infof("path strings = %v", pathStrings)
	if len(pathStrings) >= 5 {
		pluginInput.Version = pathStrings[1]
		pluginInput.ProviderName = pathStrings[2]
		pluginInput.Name = pathStrings[3]
		pluginInput.Action = pathStrings[4]
	}
	logrus.Infof("parsed request = %v", pluginInput)
	return &pluginInput
}
