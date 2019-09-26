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

func (instance ClbInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, error) {
	instances := []ResourceInstance{}
	//先获取监听器
	listenerId, err := getLbListener(providerParams, instance.Id, instance.Forward, proto, port)
	if err != nil {
		return instances, fmt.Errorf("can't found listener,lb=%s,proto=%s,port=%s", instance.Id, proto, port)
	}

	//获取后端RS
	instances, err = getLbBackendTargets(providerParams, instance.Id, instance.Forward, listenerId, instance.Vip)
	if err != nil {
		logrus.Errorf("getLbBackendTargets meet error:%v", err)
	}

	return instances, err
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

func getClassicLbListenerId(client *clb.Client, lbId string, proto string, port int64) (string, error) {
	request := clb.NewDescribeClassicalLBListenersRequest()
	request.LoadBalancerId = &lbId
	request.Protocol = &proto
	request.ListenerPort = &port

	resp, err := client.DescribeClassicalLBListeners(request)
	if err != nil {
		return "", err
	}

	if len(resp.Response.Listeners) == 0 {
		return "", fmt.Errorf("can't found listenerId by lb(%s),proto(%s),port(%v)", lbId, proto, port)
	}

	return *resp.Response.Listeners[0].ListenerId, nil
}

func getLbListener(providerParams string, id string, forward uint64, protocol string, port string) (string, error) {
	listenerId := ""
	client, _ := createClbClient(providerParams)
	proto := strings.ToUpper(protocol)
	portInt64, err := strconv.ParseInt(port, 10, 64)

	if err != nil {
		return listenerId, fmt.Errorf("%s is invalid port", port)
	}

	if forward != 0 && forward != 1 {
		return listenerId, fmt.Errorf("lb forward(%v) is invalid value", forward)
	}

	if forward == 1 {
		listenerId, err = getAppLbListenerId(client, id, proto, portInt64)
	}

	//classic
	if forward == 0 {
		listenerId, err = getClassicLbListenerId(client, id, proto, portInt64)
	}

	if err != nil {
		logrus.Errorf("getLbListener:meet err=%v", err)
	}

	return listenerId, err
}

func getAppLbBackends(client *clb.Client, lbId string, listenerId string) ([]string, error) {
	request := clb.NewDescribeTargetsRequest()
	instanceIds := []string{}
	listenerIds := []string{listenerId}
	request.LoadBalancerId = &lbId
	request.ListenerIds = common.StringPtrs(listenerIds)

	fmt.Printf("getAppLbBackends id=%v,listenerId=%v\n", lbId, listenerId)
	resp, err := client.DescribeTargets(request)
	if err != nil {
		fmt.Printf("describe target meet err=%v\n", err)
		return instanceIds, err
	}

	if len(resp.Response.Listeners) == 0 {
		fmt.Printf("listenrerNum =0\n")
		return instanceIds, fmt.Errorf("lb(%v) can't found listenerId(%s)", lbId, listenerId)
	}

	for _, target := range resp.Response.Listeners[0].Targets {
		instanceIds = append(instanceIds, *target.InstanceId)
	}
	fmt.Printf("getAppLbBackEnd=%v\n", instanceIds)

	return instanceIds, nil
}

func getClassicLbBackends(client *clb.Client, lbId string, listenerId string) ([]string, error) {
	instanceIds := []string{}
	request := clb.NewDescribeClassicalLBTargetsRequest()
	request.LoadBalancerId = &lbId

	resp, err := client.DescribeClassicalLBTargets(request)
	if err != nil {
		return instanceIds, err
	}

	for _, target := range resp.Response.Targets {
		instanceIds = append(instanceIds, *target.InstanceId)
	}

	return instanceIds, nil
}

func getLbBackendTargets(providerParams string, id string, forward uint64, listenerId string, vip string) ([]ResourceInstance, error) {
	instances := []ResourceInstance{}
	instanceIds := []string{}
	client, err := createClbClient(providerParams)

	if forward != 0 && forward != 1 {
		return instances, fmt.Errorf("lb forward(%v) is invalid value", forward)
	}

	if forward == 1 {
		instanceIds, err = getAppLbBackends(client, id, listenerId)
	}

	if forward == 0 {
		instanceIds, err = getClassicLbBackends(client, id, listenerId)
	}

	if err != nil {
		logrus.Errorf("getLbBackendTargets:meet err=%v", err)
		return instances, err
	}

	cvmType := CvmResourceType{}
	instanceMap, err := cvmType.QueryInstancesById(providerParams, instanceIds)
	if err != nil {
		logrus.Errorf("getLbBackendTargets:query meet err=%v", err)
		return instances, fmt.Errorf("getLbBackendTargets:query meet err=%v", err)
	}

	for _, instance := range instanceMap {
		cvmInstance := instance.(CvmInstance)
		cvmInstance.IsLoadBalancerBackend = true
		cvmInstance.LoadBalanceIp = vip
		instances = append(instances, cvmInstance)
	}

	return instances, err
}
