package securitygroup

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type ClbResourceType struct {
}

type ClbInstance struct {
	Id      string
	Name    string
	Forward uint64
	Region  string
	Vip     string
}

func createClbClient(providerParams string) (client *clb.Client, err error) {
	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		logrus.Errorf("createClbClient GetMapFromProviderParams meet error=%v", err)
		return nil, err
	}

	credential := common.NewCredential(paramsMap["SecretID"], paramsMap["SecretKey"])
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "clb.tencentcloudapi.com"

	return clb.NewClient(credential, paramsMap["Region"], clientProfile)
}

func (resourceType *ClbResourceType) IsSupportEgressPolicy() bool {
	logrus.Infof("ClbResourceType IsSupportEgressPolicy: return=[false]")
	return false
}

func (resourceType *ClbResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	logrus.Infof("ClbResourceType QueryInstancesById: request instanceIds=%++v", instanceIds)

	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		err := fmt.Errorf("instanceIds is empty")

		logrus.Errorf("ClbResourceType QueryInstancesById meet error=%v", err)
		return result, err
	}

	client, _ := createClbClient(providerParams)
	var offset, limit int64 = 0, int64(len(instanceIds))
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	request := clb.NewDescribeLoadBalancersRequest()
	request.LoadBalancerIds = common.StringPtrs(instanceIds)
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeLoadBalancers(request)
	if err != nil {
		logrus.Errorf("ClbResourceType QueryInstancesById DescribeLoadBalancers meet err0r=%v", err)
		return result, err
	}

	for _, lb := range resp.Response.LoadBalancerSet {
		instance := ClbInstance{
			Id:      *lb.LoadBalancerId,
			Name:    *lb.LoadBalancerName,
			Forward: *lb.Forward, // 负载均衡类型标识，1：负载均衡，0：传统型负载均衡。
			Region:  region,
		}
		if len(lb.LoadBalancerVips) > 0 {
			instance.Vip = *lb.LoadBalancerVips[0]
		}
		result[*lb.LoadBalancerId] = instance
	}

	logrus.Infof("ClbResourceType QueryInstancesById: result=%++v", result)
	return result, nil
}

func (resourceType *ClbResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	logrus.Infof("ClbResourceType QueryInstancesByIp: request ips=%++v", ips)

	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		err := fmt.Errorf("ips is empty")

		logrus.Errorf("ClbResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}

	client, _ := createClbClient(providerParams)
	var offset, limit int64 = 0, int64(len(ips))
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	request := clb.NewDescribeLoadBalancersRequest()
	request.LoadBalancerVips = common.StringPtrs(ips)
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeLoadBalancers(request)
	if err != nil {
		logrus.Errorf("ClbResourceType QueryInstancesByIp DescribeLoadBalancers meet error=%v", err)
		return result, err
	}

	for _, lb := range resp.Response.LoadBalancerSet {
		instance := ClbInstance{
			Id:      *lb.LoadBalancerId,
			Name:    *lb.LoadBalancerName,
			Forward: *lb.Forward, // 负载均衡类型标识，1：负载均衡，0：传统型负载均衡。
			Region:  region,
		}
		if len(lb.LoadBalancerVips) > 0 {
			instance.Vip = *lb.LoadBalancerVips[0]
			result[instance.Vip] = instance
		}
	}

	logrus.Infof("ClbResourceType QueryInstancesById: result=%++v", result)
	return result, nil
}

func (resourceType *ClbResourceType) IsLoadBalanceType() bool {
	logrus.Infof("ClbResourceType IsLoadBalanceType: return=[true]")
	return true
}

func (instance ClbInstance) ResourceTypeName() string {
	logrus.Infof("ClbInstance ResourceTypeName: return=[clb]")
	return "clb"
}

func (instance ClbInstance) GetId() string {
	logrus.Infof("ClbInstance GetId: return=[%v]", instance.Id)
	return instance.Id
}

func (instance ClbInstance) GetName() string {
	logrus.Infof("ClbInstance GetName: return=[%v]", instance.Name)
	return instance.Name
}

func (instance ClbInstance) GetIp() string {
	logrus.Infof("ClbInstance GetIp: return=[%v]", instance.Vip)
	return instance.Vip
}
func (instance ClbInstance) GetRegion() string {
	logrus.Infof("ClbInstance GetRegion: return=[%v]", instance.Region)
	return instance.Region
}

func (instance ClbInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	err := errors.New("clb do not support query security groups function")

	logrus.Errorf("ClbInstance QuerySecurityGroups meet error=%v", err)
	return []string{}, err
}

func (instance ClbInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	err := errors.New("clb do not support query security groups function")

	logrus.Errorf("ClbInstance AssociateSecurityGroups meet error=%v", err)
	return err
}

func (instance ClbInstance) IsSupportSecurityGroupApi() bool {
	logrus.Infof("ClbInstance IsSupportSecurityGroupApi: return=[false]")
	return false
}

func (instance ClbInstance) GetBackendTargets(providerParams string, protocol string, port string) ([]ResourceInstance, []string, error) {
	logrus.Infof("ClbInstance GetBackendTargets: reuqest protocol=%v, port=%v", protocol, port)

	instances := []ResourceInstance{}
	client, _ := createClbClient(providerParams)
	proto := strings.ToUpper(protocol)
	portInt64, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		err := fmt.Errorf("%s is invalid port", port)

		logrus.Errorf("ClbInstance GetBackendTargets ParseInt meet error=%v", err)
		return instances, []string{}, err
	}

	instanceIds := []string{}
	ports := []int64{}
	if instance.Forward == 1 {
		instanceIds, ports, err = getAppLbBackends(client, instance.Id, proto, portInt64)
		if err != nil {
			logrus.Errorf("ClbInstance GetBackendTargets getAppLbBackends meet error=%v", err)
			return instances, []string{}, err
		}
	}

	if instance.Forward == 0 {
		instanceIds, ports, err = getClassicLbBackends(client, instance.Id, proto, portInt64)
		if err != nil {
			logrus.Errorf("ClbInstance GetBackendTargets getClassicLbBackends meet error=%v", err)
			return instances, []string{}, err
		}
	}

	portsStr := []string{}
	cvmType := CvmResourceType{}

	instanceMap, err := cvmType.QueryInstancesById(providerParams, instanceIds)
	if err != nil {
		logrus.Errorf("ClbInstance GetBackendTargets QueryInstancesById meet error=%v", err)
		return instances, []string{}, err
	}

	for i, instanceId := range instanceIds {
		ins, ok := instanceMap[instanceId]
		if !ok {
			continue
		}
		cvmInstance := ins.(CvmInstance)
		cvmInstance.IsLoadBalancerBackend = true
		cvmInstance.LoadBalanceIp = instance.Vip
		instances = append(instances, cvmInstance)
		portsStr = append(portsStr, fmt.Sprintf("%v", ports[i]))
	}

	logrus.Infof("ClbInstance GetBackendTargets: return results=%++v, ports=%++v", instances, portsStr)
	return instances, portsStr, err
}

func getAppLbListenerId(client *clb.Client, lbId string, proto string, port int64) (string, error) {
	logrus.Infof("getAppLbListenerId: requst lbId=%v, protocol=%v, port=%v", lbId, proto, port)

	request := clb.NewDescribeListenersRequest()
	request.LoadBalancerId = &lbId
	request.Protocol = &proto
	request.Port = &port

	resp, err := client.DescribeListeners(request)
	if err != nil {
		logrus.Errorf("getAppLbListenerId DescribeListeners meet error=%v", err)
		return "", err
	}

	if len(resp.Response.Listeners) == 0 {
		err := fmt.Errorf("can't found listenerId by lb(%s),proto(%s),port(%v)", lbId, proto, port)

		logrus.Errorf("getAppLbListenerId DescribeListeners meet error=%v", err)
		return "", err
	}

	logrus.Infof("getAppLbListenerId: return ListenerId=%v", *resp.Response.Listeners[0].ListenerId)
	return *resp.Response.Listeners[0].ListenerId, nil
}

func getClassicLbListenerId(client *clb.Client, lbId string, proto string, port int64) (string, int64, error) {
	logrus.Infof("getClassicLbListenerId: requst lbId=%v, protocol=%v, port=%v", lbId, proto, port)

	request := clb.NewDescribeClassicalLBListenersRequest()
	request.LoadBalancerId = &lbId
	request.Protocol = &proto
	request.ListenerPort = &port

	resp, err := client.DescribeClassicalLBListeners(request)
	if err != nil {
		logrus.Errorf("getClassicLbListenerId DescribeClassicalLBListeners meet error=%v", err)
		return "", 0, err
	}

	if len(resp.Response.Listeners) == 0 {
		err := fmt.Errorf("can't found listenerId by lb(%s),proto(%s),port(%v)", lbId, proto, port)

		logrus.Errorf("getClassicLbListenerId DescribeClassicalLBListeners meet error=%v", err)
		return "", 0, err
	}

	logrus.Infof("getClassicLbListenerId: return ListenerId=%v", *resp.Response.Listeners[0].ListenerId)
	return *resp.Response.Listeners[0].ListenerId, *resp.Response.Listeners[0].InstancePort, nil
}

func getAppLbBackends(client *clb.Client, lbId string, protocol string, port int64) ([]string, []int64, error) {
	logrus.Infof("getAppLbBackends: requst lbId=%v, protocol=%v, port=%v", lbId, protocol, port)

	instanceIds := []string{}
	ports := []int64{}
	listenerId, err := getAppLbListenerId(client, lbId, protocol, port)
	if err != nil {
		logrus.Errorf("getAppLbBackends getAppLbListenerId meet error=%v", err)
		return instanceIds, ports, err
	}

	listenerIds := []string{listenerId}
	request := clb.NewDescribeTargetsRequest()
	request.ListenerIds = common.StringPtrs(listenerIds)
	request.LoadBalancerId = &lbId

	resp, err := client.DescribeTargets(request)
	if err != nil {
		logrus.Errorf("getAppLbBackends DescribeTargets meet error=%v", err)
		return instanceIds, ports, err
	}

	if len(resp.Response.Listeners) == 0 {
		err := fmt.Errorf("lb(%v) can't found listenerId(%s)", lbId, listenerId)
		logrus.Errorf("getAppLbBackends DescribeTargets meet error=%v", err)
		return instanceIds, ports, err
	}

	for _, target := range resp.Response.Listeners[0].Targets {
		instanceIds = append(instanceIds, *target.InstanceId)
		ports = append(ports, *target.Port)
	}

	logrus.Infof("getAppLbBackends: return instanceIds=%++v, ports=%++v", instanceIds, ports)
	return instanceIds, ports, nil
}

func getClassicLbBackends(client *clb.Client, lbId string, protocol string, port int64) ([]string, []int64, error) {
	logrus.Infof("getClassicLbBackends: requst lbId=%v, protocol=%v, port=%v", lbId, protocol, port)

	instanceIds := []string{}
	ports := []int64{}

	_, listenerPort, err := getClassicLbListenerId(client, lbId, protocol, port)
	if err != nil {
		logrus.Errorf("getClassicLbBackends getClassicLbListenerId meet error=%v", err)
		return instanceIds, ports, err
	}

	request := clb.NewDescribeClassicalLBTargetsRequest()
	request.LoadBalancerId = &lbId
	resp, err := client.DescribeClassicalLBTargets(request)
	if err != nil {
		logrus.Errorf("getClassicLbBackends DescribeClassicalLBTargets meet error=%v", err)
		return instanceIds, ports, err
	}

	for _, target := range resp.Response.Targets {
		instanceIds = append(instanceIds, *target.InstanceId)
		ports = append(ports, listenerPort)
	}

	logrus.Infof("getAppLbBackends: return instanceIds=%++v, ports=%++v", instanceIds, ports)
	return instanceIds, ports, nil
}
