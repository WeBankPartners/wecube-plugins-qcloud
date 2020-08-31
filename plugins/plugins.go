package plugins

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	PROVIDER_NAME = "qcloud"
	VERSION       = "v1"
)

var (
	pluginsMutex sync.Mutex
	plugins      = make(map[string]Plugin)
)

type Plugin interface {
	GetActionByName(actionName string) (Action, error)
}

type Action interface {
	ReadParam(param interface{}) (interface{}, error)
	//CheckParam(param interface{}) error
	Do(param interface{}) (interface{}, error)
}

func RegisterPlugin(name string, plugin Plugin) {
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
	RegisterPlugin("vm", new(VmPlugin))
	RegisterPlugin("storage", new(StoragePlugin))
	RegisterPlugin("security-group", new(SecurityGroupPlugin))
	RegisterPlugin("security-policy", new(SecurityPolicyPlugin))
	RegisterPlugin("subnet", new(SubnetPlugin))
	RegisterPlugin("nat-gateway", new(NatGatewayPlugin))
	RegisterPlugin("vpc", new(VpcPlugin))
	RegisterPlugin("peering-connection", new(PeeringConnectionPlugin))
	RegisterPlugin("route-table", new(RouteTablePlugin))
	RegisterPlugin("mysql", new(MysqlVmPlugin))
	RegisterPlugin("redis", new(RedisPlugin))
	//RegisterPlugin("log", new(LogPlugin))
	RegisterPlugin("elastic-nic", new(ElasticNicPlugin))
	RegisterPlugin("eip", new(EIPPlugin))
	RegisterPlugin("mariadb", new(MariadbPlugin))
	RegisterPlugin("route-policy", new(RoutePolicyPlugin))
	RegisterPlugin("clb", new(ClbPlugin))
	RegisterPlugin("cbs", new(CbsPlugin))
	RegisterPlugin("clb-target", new(ClbTargetPlugin))
	RegisterPlugin("bucket", new(BucketPlugin))
	RegisterPlugin("user", new(UserPlugin))
}

type PluginRequest struct {
	Version      string
	ProviderName string
	Name         string
	Action       string
	Parameters   interface{}
}

type PluginResponse struct {
	ResultCode string      `json:"resultCode"`
	ResultMsg  string      `json:"resultMessage"`
	Results    interface{} `json:"results"`
}

func Process(pluginRequest *PluginRequest) (*PluginResponse, error) {
	var pluginResponse = PluginResponse{}
	var err error
	defer func() {
		if err != nil {
			logrus.Errorf("plguin[%v]-action[%v] meet error = %v", pluginRequest.Name, pluginRequest.Action, err)
			pluginResponse.ResultCode = "1"
			pluginResponse.ResultMsg = fmt.Sprint(err)
		} else {
			logrus.Infof("plguin[%v]-action[%v] completed", pluginRequest.Name, pluginRequest.Action)
			pluginResponse.ResultCode = "0"
			pluginResponse.ResultMsg = "success"
		}
	}()

	logrus.Infof("plguin[%v]-action[%v] start...", pluginRequest.Name, pluginRequest.Action)

	if pluginRequest.ProviderName != PROVIDER_NAME {
		err = fmt.Errorf("ProviderName[%v] is wrong", pluginRequest.ProviderName)
		return nil, err
	}

	if pluginRequest.Version != VERSION {
		err = fmt.Errorf("Version[%v] is wrong", pluginRequest.Version)
		return nil, err
	}

	plugin, err := getPluginByName(pluginRequest.Name)
	if err != nil {
		return &pluginResponse, err
	}

	action, err := plugin.GetActionByName(pluginRequest.Action)
	if err != nil {
		return &pluginResponse, err
	}

	logrus.Infof("read parameters from http request = %v", pluginRequest.Parameters)
	actionParam, err := action.ReadParam(pluginRequest.Parameters)
	if err != nil {
		return &pluginResponse, err
	}

	logrus.Infof("action do with parameters = %v", actionParam)
	pluginResponse.Results, err = action.Do(actionParam)

	return &pluginResponse, err
}
