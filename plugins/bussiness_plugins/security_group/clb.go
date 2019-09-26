package securitygroup

import (
	"errors"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"strconv"
	"strings"
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
		return nil, err
	}

	credential := common.NewCredential(paramsMap["SecretID"], paramsMap["SecretKey"])
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "clb.tencentcloudapi.com"

	return clb.NewClient(credential, paramsMap["Region"], clientProfile)
}

func (resourceType *ClbResourceType) IsSupportEgressPolicy() bool {
	return false
}

func (resourceType *ClbResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)
	client, _ := createClbClient(providerParams)
	var offset, limit int64 = 0, int64(len(instanceIds))
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	request := clb.NewDescribeLoadBalancersRequest()
	request.LoadBalancerIds = common.StringPtrs(instanceIds)
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeLoadBalancers(request)
	if err != nil {
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
	return result, nil
}

func (resourceType *ClbResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)
	client, _ := createClbClient(providerParams)

	var offset, limit int64 = 0, int64(len(ips))
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	request := clb.NewDescribeLoadBalancersRequest()
	request.LoadBalancerVips = common.StringPtrs(ips)
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeLoadBalancers(request)
	if err != nil {
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
	return result, nil
}

func (resourceType *ClbResourceType) IsLoadBalanceType() bool {
	return true
}

func (instance ClbInstance) ResourceTypeName() string {
	return "clb"
}

func (instance ClbInstance) GetId() string {
	return instance.Id
}

func (instance ClbInstance) GetName() string {
	return instance.Name
}

func (instance ClbInstance) GetIp() string {
	return instance.Vip
}
func (instance ClbInstance) GetRegion() string {
	return instance.Region
}

func (instance ClbInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	return []string{}, errors.New("clb do not support query security groups function")
}

func (instance ClbInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	return errors.New("clb do not associate security groups function")
}

func (instance ClbInstance) IsSupportSecurityGroupApi() bool {
	return false
}

func (instance ClbInstance) GetBackendTargets(providerParams string, protocol string, port string) ([]ResourceInstance, []string, error) {
	instances := []ResourceInstance{}
	client, _ := createClbClient(providerParams)
	proto := strings.ToUpper(protocol)
	portInt64, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return instances, []string{}, fmt.Errorf("%s is invalid port", port)
	}

	instanceIds := []string{}
	ports := []int64{}
	if instance.Forward == 1 {
		instanceIds, ports, err = getAppLbBackends(client, instance.Id, proto, portInt64)
	}

	if instance.Forward == 0 {
		instanceIds, ports, err = getClassicLbBackends(client, instance.Id, proto, portInt64)
	}

	if err != nil {
		return instances, []string{}, err
	}

	portsStr := []string{}
	cvmType := CvmResourceType{}

	instanceMap, err := cvmType.QueryInstancesById(providerParams, instanceIds)
	if err != nil {
		logrus.Errorf("getLbBackendTargets:query meet err=%v", err)
		return instances, []string{}, fmt.Errorf("getLbBackendTargets:query meet err=%v", err)
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

	return instances, portsStr, err
}

func getAppLbListenerId(client *clb.Client, lbId string, proto string, port int64) (string, error) {
	request := clb.NewDescribeListenersRequest()
	request.LoadBalancerId = &lbId
	request.Protocol = &proto
	request.Port = &port

	resp, err := client.DescribeListeners(request)
	if err != nil {
		return "", err
	}

	if len(resp.Response.Listeners) == 0 {
		return "", fmt.Errorf("can't found listenerId by lb(%s),proto(%s),port(%v)", lbId, proto, port)
	}
	return *resp.Response.Listeners[0].ListenerId, nil
}

func getClassicLbListenerId(client *clb.Client, lbId string, proto string, port int64) (string, int64, error) {
	request := clb.NewDescribeClassicalLBListenersRequest()
	request.LoadBalancerId = &lbId
	request.Protocol = &proto
	request.ListenerPort = &port

	resp, err := client.DescribeClassicalLBListeners(request)
	if err != nil {
		return "", 0, err
	}

	if len(resp.Response.Listeners) == 0 {
		return "", 0, fmt.Errorf("can't found listenerId by lb(%s),proto(%s),port(%v)", lbId, proto, port)
	}

	return *resp.Response.Listeners[0].ListenerId, *resp.Response.Listeners[0].InstancePort, nil
}

func getAppLbBackends(client *clb.Client, lbId string, protocol string, port int64) ([]string, []int64, error) {
	instanceIds := []string{}
	ports := []int64{}
	listenerId, err := getAppLbListenerId(client, lbId, protocol, port)
	if err != nil {
		return instanceIds, ports, err
	}

	listenerIds := []string{listenerId}
	request := clb.NewDescribeTargetsRequest()
	request.ListenerIds = common.StringPtrs(listenerIds)
	request.LoadBalancerId = &lbId

	resp, err := client.DescribeTargets(request)
	if err != nil {
		fmt.Printf("describe target meet err=%v\n", err)
		return instanceIds, ports, err
	}

	if len(resp.Response.Listeners) == 0 {
		return instanceIds, ports, fmt.Errorf("lb(%v) can't found listenerId(%s)", lbId, listenerId)
	}

	for _, target := range resp.Response.Listeners[0].Targets {
		instanceIds = append(instanceIds, *target.InstanceId)
		ports = append(ports, *target.Port)
	}

	return instanceIds, ports, nil
}

func getClassicLbBackends(client *clb.Client, lbId string, protocol string, port int64) ([]string, []int64, error) {
	instanceIds := []string{}
	ports := []int64{}

	_, listenerPort, err := getClassicLbListenerId(client, lbId, protocol, port)
	if err != nil {
		return instanceIds, ports, err
	}

	request := clb.NewDescribeClassicalLBTargetsRequest()
	request.LoadBalancerId = &lbId
	resp, err := client.DescribeClassicalLBTargets(request)
	if err != nil {
		return instanceIds, ports, err
	}

	for _, target := range resp.Response.Targets {
		instanceIds = append(instanceIds, *target.InstanceId)
		ports = append(ports, listenerPort)
	}

	return instanceIds, ports, nil
}
