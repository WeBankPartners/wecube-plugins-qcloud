package securitygroup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

const (
	MAX_SEUCRITY_RULE_NUM = 100
)

var (
	ErrorIpNotFound = errors.New("ip not found")
)

//interface definition
type ResourceInstance interface {
	ResourceTypeName() string
	GetId() string
	GetName() string
	GetRegion() string
	GetIp() string
	QuerySecurityGroups(providerParams string) ([]string, error)
	AssociateSecurityGroups(providerParams string, securityGroups []string) error
	IsSupportSecurityGroupApi() bool
	GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error)
}

type ResourceType interface {
	QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error)
	QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error)
	IsLoadBalanceType() bool
	IsSupportEgressPolicy() bool
}

//resourceType register
var (
	resTypesMutex   sync.Mutex
	resourceTypeMap = make(map[string]ResourceType)
)

//resource type  register
func addNewResourceType(name string, newResourceType ResourceType) error {
	resTypesMutex.Lock()
	defer resTypesMutex.Unlock()

	if _, found := resourceTypeMap[name]; found {
		logrus.Errorf("resourceType(%s) was registered twice", name)
	}

	resourceTypeMap[name] = newResourceType
	return nil
}

func getResouceTypeByName(name string) (ResourceType, error) {
	resTypesMutex.Lock()
	defer resTypesMutex.Unlock()

	resType, found := resourceTypeMap[name]
	if !found {
		err := fmt.Errorf("resourceType[%s] not found", name)

		logrus.Errorf("getResouceTypeByName meet error=%v", err)
		return nil, err
	}

	return resType, nil
}

func unmarshalJson(source interface{}, target interface{}) error {
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

type BussinessSecurityGroupPlugin struct {
}

func (plugin *BussinessSecurityGroupPlugin) GetActionByName(actionName string) (plugins.Action, error) {
	logrus.Infof("BussinessSecurityGroupPlugin GetActionByName: request actionName=%v", actionName)

	action, found := SecurityGroupActions[actionName]

	if !found {
		err := fmt.Errorf("Bussiness Security Group plugin,action = %s not found", actionName)
		logrus.Errorf("BussinessSecurityGroupPlugin GetActionByName meet error=%v", err)
		return nil, err
	}

	return action, nil
}

var SecurityGroupActions = make(map[string]plugins.Action)

func init() {
	//plugin registry
	plugins.RegisterPlugin("bs-security-group", new(BussinessSecurityGroupPlugin))

	//resourceType registry
	addNewResourceType("cvm", new(CvmResourceType))
	addNewResourceType("clb", new(ClbResourceType))
	addNewResourceType("mysql", new(MysqlResourceType))
	/*addNewResourceType("bm", new(BmResourceType))
	addNewResourceType("bmlb", new(BmlbResourceType))
	addNewResourceType("mariadb", new(MariadbResourceType))
	addNewResourceType("redis", new(RedisResourceType))
	addNewResourceType("mongodb", new(MongodbResourceType))*/

	//action
	SecurityGroupActions["calc-security-policies"] = new(CalcSecurityPolicyAction)
	SecurityGroupActions["apply-security-policies"] = new(ApplySecurityPolicyAction)
}

type QueryIpsResult struct {
	Err         error
	InstanceMap map[string]ResourceInstance
}

func queryOneRegionInstanceByIps(providerParams string, region string, ips []string, ch chan QueryIpsResult) {
	result := QueryIpsResult{
		Err:         nil,
		InstanceMap: make(map[string]ResourceInstance),
	}
	start := time.Now()
	defer func() {
		logrus.Infof("queryOneRegionInstanceByIps region(%s) ips (%v) taken %v,result=%++v", region, ips, time.Since(start), result)
	}()

	rtnIps := 0
	for _, resType := range resourceTypeMap {
		instanceMap, err := resType.QueryInstancesByIp(providerParams, ips)
		logrus.Infof("findInstanceByIp QueryInstancesByIp instanceMap:%++v", instanceMap)
		if err != nil {
			result.Err = err
			logrus.Errorf("findInstanceByIp QueryInstancesByIp meet error=%v\n", err)
			break
		}

		for key, value := range instanceMap {
			result.InstanceMap[key] = value
			rtnIps++
		}
		if rtnIps == len(ips) {
			break
		}
	}

	ch <- result
}

func getResourceAllIp(sourceIps []string, destIps []string) (map[string]ResourceInstance, error) {
	totalMap := make(map[string]ResourceInstance)
	chResult := make(chan QueryIpsResult)
	regions, err := getRegions()
	if err != nil {
		logrus.Errorf("findInstanceByIp getRegions meet err=%v\n", err)
		return nil, err
	}

	ips := []string{}
	ipmap := make(map[string]bool)
	for _, ip := range sourceIps {
		ipmap[ip] = true
	}
	for _, ip := range destIps {
		ipmap[ip] = true
	}
	for key := range ipmap {
		ips = append(ips, key)
	}

	for _, region := range regions {
		providerParams, err := getProviderParams(region)
		if err != nil {
			return totalMap, err
		}
		go queryOneRegionInstanceByIps(providerParams, region, ips, chResult)
	}

	returnedIp := 0
	for _, _ = range regions {
		result := <-chResult
		if result.Err != nil {
			return totalMap, result.Err
		}
		for key, value := range result.InstanceMap {
			totalMap[key] = value
			returnedIp++
		}

		if returnedIp == len(ips) {
			return totalMap, nil
		}
	}

	return totalMap, nil
}

func findInstanceByIp(ip string, ipMap map[string]ResourceInstance) (ResourceInstance, error) {
	instance, ok := ipMap[ip]
	if ok {
		return instance, nil
	}
	return nil, fmt.Errorf("Ip(%s) not found", ip)
}

//---------------calc security policy action------------------------------//
type CalcSecurityPoliciesRequest struct {
	Protocol         string   `json:"protocol"`
	SourceIps        []string `json:"source_ips"`
	DestIps          []string `json:"dest_ips"`
	DestPort         string   `json:"dest_port"`
	PolicyAction     string   `json:"policy_action"`
	PolicyDirections []string `json:"policy_directions"`
	Description      string   `json:"description"`
}

type SecurityPolicy struct {
	Ip                      string `json:"ip"`
	Type                    string `json:"type"`
	Id                      string `json:"id"`
	Region                  string `json:"region"`
	SupportSecurityGroupApi bool   `json:"support_security_group_api"`
	PeerIp                  string `json:"peer_ip"`
	Protocol                string `json:"protocol"`
	Ports                   string `json:"ports"`
	Action                  string `json:"action"`
	Description             string `json:"description"`
	ErrorMsg                string `json:"err_msg,omitempty"`
	UndoReason              string `json:"undo_reason,omitempty"`
	SecurityGroupId         string `json:"security_group_id,omitempty"`
}

type CalcSecurityPoliciesResult struct {
	TimeTaken string `json:"time_taken"`

	IngressPoliciesTotal int `json:"ingress_policies_total"`
	EgressPoliciesTotal  int `json:"egress_policies_total"`

	IngressPolicies []SecurityPolicy `json:"ingress_policies"`
	EgressPolicies  []SecurityPolicy `json:"egress_policies"`
}

type CalcSecurityPolicyAction struct {
}

func (action *CalcSecurityPolicyAction) ReadParam(param interface{}) (interface{}, error) {
	var input CalcSecurityPoliciesRequest
	err := unmarshalJson(param, &input)
	if err != nil {
		logrus.Errorf("CalcSecurityPolicyAction ReadParam UnmarshalJson: failed to unmarsh, err=%v, param=%v", err, param)
		return nil, err
	}

	logrus.Infof("CalcSecurityPolicyAction ReadParam: return=%++v", input)
	return input, nil
}

func (action *CalcSecurityPolicyAction) CheckParam(input interface{}) error {
	req, _ := input.(CalcSecurityPoliciesRequest)
	if err := isValidProtocol(req.Protocol); err != nil {
		logrus.Errorf("CalcSecurityPolicyAction CheckParam isValidProtocol meet error=%v", err)
		return err
	}

	if err := isValidAction(req.PolicyAction); err != nil {
		logrus.Errorf("CalcSecurityPolicyAction CheckParam isValidAction meet error=%v", err)
		return err
	}

	for _, ip := range req.SourceIps {
		if err := isValidIp(ip); err != nil {
			logrus.Errorf("CalcSecurityPolicyAction CheckParam isValidIp meet error=%v", err)
			return err
		}
	}

	for _, ip := range req.DestIps {
		if err := isValidIp(ip); err != nil {
			logrus.Errorf("CalcSecurityPolicyAction CheckParam isValidIp meet error=%v", err)
			return err
		}
	}

	_, err := getPortsByPolicyFormat(req.DestPort)
	if err != nil {
		logrus.Errorf("CalcSecurityPolicyAction CheckParam getPortsByPolicyFormat meet error=%v", err)
		return err
	}

	for _, direction := range req.PolicyDirections {
		if err := isValidDirection(direction); err != nil {
			logrus.Errorf("CalcSecurityPolicyAction CheckParam isValidDirection meet error=%v", err)
			return err
		}
	}

	return nil
}

func newPolicies(instance ResourceInstance, myIp string, peerIp string, proto string, port string, action string, desc string) ([]SecurityPolicy, error) {
	logrus.Infof("newPolicies: request instance=%++v, myIp=%v, peerIp=%v, protocol=%v, port=%v, action=%v, description=%v", instance, myIp, peerIp, proto, port, action, desc)

	policies := []SecurityPolicy{}
	resType, _ := getResouceTypeByName(instance.ResourceTypeName())

	//非LB设备
	if false == resType.IsLoadBalanceType() {
		newPolicy := SecurityPolicy{
			Ip:                      myIp,
			Type:                    instance.ResourceTypeName(),
			Id:                      instance.GetId(),
			Region:                  instance.GetRegion(),
			SupportSecurityGroupApi: instance.IsSupportSecurityGroupApi(),
			PeerIp:                  peerIp,
			Protocol:                proto,
			Ports:                   port,
			Action:                  action,
			Description:             desc,
		}
		policies := append(policies, newPolicy)

		logrus.Infof("newPolicies: return policies=%++v", policies)
		return policies, nil
	}

	//LB设备
	providerParams, _ := getProviderParams(instance.GetRegion())
	splitPorts := strings.Split(port, ",")

	for _, splitPort := range splitPorts {
		if _, err := strconv.Atoi(splitPort); err != nil {
			err := fmt.Errorf("loadbalancer do not support port format like %s", port)

			logrus.Errorf("newPolicies strconv.Atoi meet error=%v", err)
			return policies, err
		}
		instances, ports, err := instance.GetBackendTargets(providerParams, proto, splitPort)
		if err != nil {
			logrus.Errorf("newPolicies GetBackendTargets meet error=%v", err)
			return policies, err
		}
		if len(instances) == 0 {
			err := fmt.Errorf("loadbalancer(%s) port (%v) do not have any backends", instance.GetIp(), splitPort)
			logrus.Errorf("newPolicies GetBackendTargets meet error=%v", err)
			return policies, err
		}

		for i, backendInstance := range instances {
			newPolicy := SecurityPolicy{
				Ip:                      backendInstance.GetIp(),
				Type:                    backendInstance.ResourceTypeName(),
				Id:                      backendInstance.GetId(),
				Region:                  backendInstance.GetRegion(),
				SupportSecurityGroupApi: backendInstance.IsSupportSecurityGroupApi(),
				PeerIp:                  peerIp,
				Protocol:                proto,
				Ports:                   ports[i],
				Action:                  action,
				Description:             desc,
			}
			policies = append(policies, newPolicy)
		}
	}

	logrus.Infof("newPolicies: return policies=%++v", policies)
	return policies, nil
}

func calcPolicies(devIp string, ipMap map[string]ResourceInstance, peerIps []string, proto string, ports []string,
	action string, description string, direction string) ([]SecurityPolicy, error) {
	logrus.Infof("calcPolicies: reuqest devIp=%v, peerIps=%++v, protocol=%v, ports=%++v, action=%v, description=%v, direction=%v", devIp, peerIps, proto, ports, action, description, direction)

	policies := []SecurityPolicy{}

	//check if dev exist
	instance, err := findInstanceByIp(devIp, ipMap)
	if err != nil {
		logrus.Errorf("calcPolicies findInstanceByIp meet error=%v", err)
		return policies, err
	}

	resType, err := getResouceTypeByName(instance.ResourceTypeName())
	if err != nil {
		logrus.Errorf("calcPolicies getResouceTypeByName meet error=%v", err)
		return policies, err
	}

	if direction == EGRESS_RULE {
		if false == resType.IsSupportEgressPolicy() {
			err := fmt.Errorf("%s is %s device,do not support egress", devIp, instance.ResourceTypeName())
			logrus.Errorf("calcPolicies IsSupportEgressPolicy meet error=%v", err)
			return policies, err
		}
	}

	for _, peerIp := range peerIps {
		peerInstance, err := findInstanceByIp(peerIp, ipMap)
		logrus.Infof("calcPolicies findInstanceByip peerIp=%s, instance=%++v, err=%v\n", peerIp, peerInstance, err)
		if err == nil {
			peerResType, _ := getResouceTypeByName(peerInstance.ResourceTypeName())
			if direction == INGRESS_RULE && nil != peerResType && peerResType.IsLoadBalanceType() {
				err := fmt.Errorf("对端设备(%s) 是负载均衡设备,入栈规则不支持对端IP为负载均衡设备", peerIp)
				logrus.Infof("calcPolicies getResouceTypeByName meet error=%v", err)
				return policies, err
			}
		}

		for _, port := range ports {
			newPolicies, err := newPolicies(instance, devIp, peerIp, proto, port, action, description)
			if err != nil {
				logrus.Errorf("calcPolicies newPolicies meet error=%v", err)
				return policies, err
			}
			if len(newPolicies) > 0 {
				policies = append(policies, newPolicies...)
			}
		}
	}

	logrus.Infof("calcPolicies: retuern policies=%++v", policies)
	return policies, nil
}

func (action *CalcSecurityPolicyAction) Do(input interface{}) (interface{}, error) {
	req, _ := input.(CalcSecurityPoliciesRequest)
	logrus.Infof("CalcSecurityPolicyAction Do: request input=%++v", input)

	result := CalcSecurityPoliciesResult{}
	start := time.Now()
	ports, _ := getPortsByPolicyFormat(req.DestPort)
	logrus.Infof("CalcSecurityPolicyAction Do: ports=%++v", ports)

	ipMaps, err := getResourceAllIp(req.SourceIps, req.DestIps)
	logrus.Infof("CalcSecurityPolicyAction Do getResourceAllIp: len(ipMaps)=%v ipMaps=%++v", len(ipMaps), ipMaps)

	if err != nil {
		logrus.Infof("getResourceAllIp meet err=%v", err)
		result.TimeTaken = fmt.Sprintf("%v", time.Since(start))
		return result, err
	}
	//calc egress policies
	if isContainInList(EGRESS_RULE, req.PolicyDirections) {
		for _, ip := range req.SourceIps {
			policies, err := calcPolicies(ip, ipMaps, req.DestIps, req.Protocol, ports, req.PolicyAction, req.Description, EGRESS_RULE)
			result.EgressPolicies = append(result.EgressPolicies, policies...)
			if err != nil {
				result.TimeTaken = fmt.Sprintf("%v", time.Since(start))
				return result, err
			}
		}
	}

	//calc ingress policies
	if isContainInList(INGRESS_RULE, req.PolicyDirections) {
		for _, ip := range req.DestIps {
			policies, err := calcPolicies(ip, ipMaps, req.SourceIps, req.Protocol, ports, req.PolicyAction, req.Description, INGRESS_RULE)
			result.IngressPolicies = append(result.IngressPolicies, policies...)
			if err != nil {
				result.TimeTaken = fmt.Sprintf("%v", time.Since(start))
				return result, err
			}
		}
	}

	result.TimeTaken = fmt.Sprintf("%v", time.Since(start))
	result.IngressPoliciesTotal = len(result.IngressPolicies)
	result.EgressPoliciesTotal = len(result.EgressPolicies)

	logrus.Infof("CalcSecurityPolicyAction Do: return result=%++v", result)
	return result, nil
}

//---------------apply security policy action------------------------------//
type ApplySecurityPolicyAction struct {
}

type ApplySecurityPoliciesRequest struct {
	IngressPolicies []SecurityPolicy `json:"ingress_policies"`
	EgressPolicies  []SecurityPolicy `json:"egress_policies"`
}

type ApplyResult struct {
	PoliciesTotal int `json:"policies_total"`

	SuccessTotal int `json:"success_policies_total"`
	UndoTotal    int `json:"undo_policies_total"`
	FailedTotal  int `json:"failed_policies_total"`

	SuccessPolicies []SecurityPolicy `json:"success_policies"`
	UndoPolicies    []SecurityPolicy `json:"undo_policies"`
	FailedPolicies  []SecurityPolicy `json:"failed_policies"`
}

type ApplySecurityPoliciesResult struct {
	TimeTaken          string      `json:"time_taken"`
	IngressApplyResult ApplyResult `json:"ingress"`
	EgressApplyResult  ApplyResult `json:"egress"`
}

func (action *ApplySecurityPolicyAction) ReadParam(param interface{}) (interface{}, error) {
	var input ApplySecurityPoliciesRequest
	err := unmarshalJson(param, &input)
	if err != nil {
		logrus.Errorf("ApplySecurityPolicyAction:unmarshal failed,err=%v,param=%v", err, param)
		return nil, err
	}
	logrus.Infof("ApplySecurityPolicyAction ReadParam: input=%++v", input)
	return input, nil
}

func (action *ApplySecurityPolicyAction) CheckParam(input interface{}) error {
	req, _ := input.(ApplySecurityPoliciesRequest)
	logrus.Infof("ApplySecurityPolicyAction CheckParam: req=%++v", req)

	for _, policy := range req.IngressPolicies {
		if policy.Ip == "" || policy.Id == "" {
			return errors.New("ingress policy have empty value")
		}
	}

	for _, policy := range req.EgressPolicies {
		if policy.Ip == "" || policy.Id == "" {
			return errors.New("egress policy have empty value")
		}
	}

	return nil
}

func (action *ApplySecurityPolicyAction) Do(input interface{}) (interface{}, error) {
	var err error
	req, _ := input.(ApplySecurityPoliciesRequest)
	result := ApplySecurityPoliciesResult{}
	start := time.Now()
	logrus.Infof("ApplySecurityPolicyAction Do: req=%++v", req)

	result.IngressApplyResult = applyPolicies(req.IngressPolicies, INGRESS_RULE)
	result.EgressApplyResult = applyPolicies(req.EgressPolicies, EGRESS_RULE)

	result.TimeTaken = fmt.Sprintf("%v", time.Since(start))
	if result.IngressApplyResult.FailedTotal > 0 || result.EgressApplyResult.FailedTotal > 0 {
		err = errors.New("have some failed polices,please check policy applied detail")
	}

	logrus.Infof("ApplySecurityPolicyAction Do: result=%++v", result)
	return result, err
}

func fillSecuityPoliciesWithErrMsg(policies []*SecurityPolicy, err error) {
	for _, policy := range policies {
		policy.ErrorMsg = err.Error()
	}
}

func applyPolicies(policies []SecurityPolicy, direction string) ApplyResult {
	logrus.Infof("applyPolicies: input policies=%++v direction=%++v", policies, direction)

	result := ApplyResult{}
	instanceMap := make(map[string][]*SecurityPolicy)

	for i, _ := range policies {
		if strings.HasPrefix(policies[i].Type, "clb-cvm") {
			policies[i].Type = "cvm"
		}

		if policies[i].SupportSecurityGroupApi == true {
			key := policies[i].Ip
			instanceMap[key] = append(instanceMap[key], &policies[i])
		} else {
			policies[i].UndoReason = fmt.Sprintf("instanceType(%s) do not support security group api", policies[i].Type)
			result.UndoPolicies = append(result.UndoPolicies, policies[i])
		}
	}
	logrus.Infof("applyPolicies: instanceMap=%++v", instanceMap)

	for _, policies := range instanceMap {
		resType, err := getResouceTypeByName(policies[0].Type)
		if err != nil {
			logrus.Errorf("applyPolicies getResouceTypeByName meet error=%v", err)
			fillSecuityPoliciesWithErrMsg(policies, err)
			continue
		}

		providerParams, err := getProviderParams(policies[0].Region)
		if err != nil {
			logrus.Errorf("applyPolicies getProviderParams meet error=%v", err)
			fillSecuityPoliciesWithErrMsg(policies, err)
			continue
		}

		instances, err := resType.QueryInstancesById(providerParams, []string{policies[0].Id})
		if err != nil {
			logrus.Errorf("applyPolicies QueryInstancesById meet error=%v", err)
			fillSecuityPoliciesWithErrMsg(policies, err)
			continue
		}
		if len(instances) == 0 {
			err := fmt.Errorf("can't found instanceId(%s)", policies[0].Id)
			logrus.Errorf("applyPolicies QueryInstancesById meet error=%v", err)

			fillSecuityPoliciesWithErrMsg(policies, err)
			continue
		}
		instance := instances[policies[0].Id]
		logrus.Infof("applyPolicies instance=%++v", instance)

		existSecurityGroups, err := instance.QuerySecurityGroups(providerParams)
		if err != nil {
			logrus.Errorf("applyPolicies QuerySecurityGroups meet error=%v", err)
			fillSecuityPoliciesWithErrMsg(policies, err)
			continue
		}

		logrus.Infof("applyPolicies existSecurityGroups=%++v", existSecurityGroups)
		newSecurityGroups, err := createPolicies(providerParams, existSecurityGroups, policies, direction)
		if err != nil {
			logrus.Errorf("applyPolicies createPolicies meet error=%v", err)

			destroyPolicies(providerParams, policies, direction)
			fillSecuityPoliciesWithErrMsg(policies, err)
			continue
		}
		logrus.Infof("applyPolicies newSecurityGroups:%v", newSecurityGroups)

		if len(newSecurityGroups) > 0 {
			groups := []string{}
			groups = append(groups, newSecurityGroups...)
			groups = append(groups, existSecurityGroups...)

			if err = instance.AssociateSecurityGroups(providerParams, groups); err != nil {
				logrus.Errorf("applyPolicies AssociateSecurityGroups meet error=%v", err)

				destroyPolicies(providerParams, policies, direction)
				bindError := fmt.Errorf("resourceType(%s) instance(%s) AssociateSecurityGroups[%v] meet err=%v", policies[0].Type, policies[0].Ip, groups, err)
				fillSecuityPoliciesWithErrMsg(policies, bindError)
				continue
			}
		}
	}

	for _, policies := range instanceMap {
		for _, policy := range policies {
			if policy.ErrorMsg == "" {
				result.SuccessPolicies = append(result.SuccessPolicies, *policy)
			} else {
				result.FailedPolicies = append(result.FailedPolicies, *policy)
			}
		}
	}
	result.PoliciesTotal = len(policies)
	result.SuccessTotal = len(result.SuccessPolicies)
	result.FailedTotal = len(result.FailedPolicies)

	logrus.Infof("applyPolicies: result=%++v", result)
	return result
}

//自动构建的安全组的名称格式ip-auoto-1,ip_auto_2
func getAutoCreatedSecurityGroups(ip string, allSecurityGroupsNames, allSecurityGroupsIds []string) ([]string, int, error) {
	logrus.Infof("getAutoCreatedSecurityGroups: input ip=%v allSecurityGroupsNames=%++v allSecurityGroupsIds=%++v", ip, allSecurityGroupsNames, allSecurityGroupsIds)

	var err error
	maxAutoCreatedNum := 0
	createdSecurityGroups := []string{}
	nums := []int{}
	for i, securityGroup := range allSecurityGroupsNames {
		elements := strings.Split(securityGroup, "-")
		if len(elements) == 3 {
			if elements[0] == ip && elements[1] == "auto" {
				if num, err := strconv.Atoi(elements[2]); err == nil {
					createdSecurityGroups = append(createdSecurityGroups, allSecurityGroupsIds[i])
					nums = append(nums, num)
					if maxAutoCreatedNum < num {
						maxAutoCreatedNum = num
					}
				}
			}
		}
	}
	createdSecurityGroups, err = sortSecurityGroupsIds(nums, createdSecurityGroups)
	if err != nil {
		logrus.Errorf("getAutoCreatedSecurityGroups sortSecurityGroupsIds meet error=%v", err)
		return createdSecurityGroups, maxAutoCreatedNum + 1, err
	}
	logrus.Infof("getAutoCreatedSecurityGroups createdSecurityGroups:%v", createdSecurityGroups)

	return createdSecurityGroups, maxAutoCreatedNum + 1, nil
}

func sortSecurityGroupsIds(num []int, securityGroupsIds []string) ([]string, error) {
	logrus.Infof("sortSecurityGroupsIds; input num=%++v securityGroupsIds=%++v", num, securityGroupsIds)

	if len(num) != len(securityGroupsIds) {
		err := fmt.Errorf("sortSecurityGroupsIds error: lengths of two arrays is not equal")
		logrus.Errorf("sortSecurityGroupsIds meet error=%v", err)

		return []string{}, err
	}
	flag := 1
	for i := 0; i < len(num) && flag == 1; i++ {
		flag = 0
		for j := 0; j < len(num)-i-1; j++ {
			if num[j] > num[j+1] {
				num[j], num[j+1] = num[j+1], num[j]
				securityGroupsIds[j], securityGroupsIds[j+1] = securityGroupsIds[j+1], securityGroupsIds[j]
				flag = 1
			}
		}
	}

	logrus.Infof("sortSecurityGroupsIds: return securityGroupsIds=%++v", securityGroupsIds)
	return securityGroupsIds, nil
}

func getSecurityGroupFreePolicyNum(providerParams string, securityGroup string, direction string) (int, error) {
	logrus.Infof("getSecurityGroupFreePolicyNum: input securityGroup=%v direction=%v", securityGroup, direction)

	policiesSet, err := plugins.QuerySecurityGroupPolicies(providerParams, securityGroup)
	if err != nil {
		logrus.Errorf("getSecurityGroupFreePolicyNum meet error=%v\n", err)
		return 0, err
	}

	if strings.EqualFold(direction, INGRESS_RULE) {
		return MAX_SEUCRITY_RULE_NUM - len(policiesSet.Ingress), nil
	}

	return MAX_SEUCRITY_RULE_NUM - len(policiesSet.Egress), nil
}

func getSecurityGroupNames(providerParams string, securityGroupIds []string) ([]string, error) {
	logrus.Infof("getSecurityGroupNames: input securityGroupIds=%++v", securityGroupIds)

	securityGroupNames := []string{}
	idNameMap := make(map[string]string)
	securityGroupSet, err := plugins.QuerySecurityGroups(providerParams, securityGroupIds)
	if err != nil {
		logrus.Errorf("getSecurityGroupNames QuerySecurityGroups meet error=%v", err)
		return securityGroupNames, err
	}

	for _, securityGroup := range securityGroupSet {
		idNameMap[*securityGroup.SecurityGroupId] = *securityGroup.SecurityGroupName
	}

	for _, id := range securityGroupIds {
		if name, ok := idNameMap[id]; ok {
			securityGroupNames = append(securityGroupNames, name)
		} else {
			err := fmt.Errorf("can't found groupId(%s) detail", id)
			logrus.Errorf("getSecurityGroupNames meet error=%v", err)

			return securityGroupNames, err
		}
	}

	logrus.Infof("getSecurityGroupNames: return securityGroupNames=%++v", securityGroupNames)
	return securityGroupNames, nil
}

//format ip-auto-2
func createNewAutomationSecurityGroups(providerParams string, ip string, newCreatedSecurityGroupNum int, auotNumIndex int) ([]string, error) {
	logrus.Infof("createNewAutomationSecurityGroups: input ip=%v newCreatedSecurityGroupNum=%v auotNumIndex=%v", ip, newCreatedSecurityGroupNum, auotNumIndex)

	newSecurityGroupIds := []string{}
	for i := 0; i < newCreatedSecurityGroupNum; i++ {
		securityGroupName := fmt.Sprintf("%s-auto-%d", ip, auotNumIndex+i)
		securityGroupId, err := plugins.CreateSecurityGroup(providerParams, securityGroupName, "automation created")
		if err != nil {
			logrus.Errorf("createNewAutomationSecurityGroups CreateSecurityGroup meet err=%v", err)
			return newSecurityGroupIds, err
		}
		newSecurityGroupIds = append(newSecurityGroupIds, securityGroupId)
	}

	logrus.Errorf("createNewAutomationSecurityGroups: return newSecurityGroupIds=%++v", newSecurityGroupIds)
	return newSecurityGroupIds, nil
}

func newSecurityPolicySet(policies []*SecurityPolicy, direction string, isSetPolicyIndex bool) vpc.SecurityGroupPolicySet {
	logrus.Infof("newSecurityPolicySet: input policies=%++v direction=%v isSetPolicyIndex=%v", policies, direction, isSetPolicyIndex)

	securityPolicies := []*vpc.SecurityGroupPolicy{}
	var policyIndex int64 = 0

	for _, policy := range policies {
		action := strings.ToUpper(policy.Action)
		securityPolicy := vpc.SecurityGroupPolicy{
			Protocol:          &policy.Protocol,
			Port:              &policy.Ports,
			CidrBlock:         &policy.PeerIp,
			Action:            &action,
			PolicyDescription: &policy.Description,
		}
		if isSetPolicyIndex {
			securityPolicy.PolicyIndex = &policyIndex
		}
		securityPolicies = append(securityPolicies, &securityPolicy)
	}

	securityGroupPolicySet := vpc.SecurityGroupPolicySet{}
	if strings.EqualFold(direction, INGRESS_RULE) {
		securityGroupPolicySet.Ingress = securityPolicies
	} else {
		securityGroupPolicySet.Egress = securityPolicies
	}

	logrus.Infof("newSecurityPolicySet: return securityGroupPolicySet=%++v", securityGroupPolicySet)
	return securityGroupPolicySet
}

func addPoliciesToSecurityGroup(providerParams string, securityGroupId string, policies []*SecurityPolicy, direction string) error {
	logrus.Infof("addPoliciesToSecurityGroup: input securityGroupId=%v policies=%++v direction=%v", securityGroupId, policies, direction)

	req := vpc.NewCreateSecurityGroupPoliciesRequest()
	req.SecurityGroupId = &securityGroupId
	var err error

	if len(policies) == 0 {
		return nil
	}
	defer func() {
		if err != nil {
			logrus.Errorf("addPoliciesToSecurityGroup add policy to securityGroup(%s) meet err =%v", securityGroupId, err)
			errMsg := fmt.Sprintf("addPoliciesToSecurityGroup add policy to securityGroup(%s) meet err =%v", securityGroupId, err)
			for _, policy := range policies {
				policy.ErrorMsg = errMsg
			}
		}
	}()

	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	client, err := plugins.CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		logrus.Errorf("addPoliciesToSecurityGroup CreateVpcClient meet error=%v", err)
		return err
	}

	securityGroupPolicySet := newSecurityPolicySet(policies, direction, true)
	req.SecurityGroupPolicySet = &securityGroupPolicySet
	if _, err = client.CreateSecurityGroupPolicies(req); err == nil {
		for _, policy := range policies {
			policy.SecurityGroupId = securityGroupId
		}
	}

	return err
}

func createPolicies(providerParams string, existSecurityGroups []string, policies []*SecurityPolicy, direction string) ([]string, error) {
	logrus.Infof("createPolicies: input existSecurityGroups=%++v policies=%++v direction=%v", existSecurityGroups, policies, direction)

	newSecurityGroups := []string{}
	freePolicyNumMap := make(map[string]int)
	freePoliciesNum := 0
	securityGroupsIds := []string{}

	if len(policies) == 0 {
		return newSecurityGroups, nil
	}

	securityGroupsNames, err := getSecurityGroupNames(providerParams, existSecurityGroups)
	if err != nil {
		logrus.Errorf("createPolicies getSecurityGroupNames meet error=%v", err)
		return newSecurityGroups, err
	}
	logrus.Infof("createPolicies getSecurityGroupNames: securityGroupsNames:%v", securityGroupsNames)

	createdSecurityGroups, autoCreatedStartIndex, err := getAutoCreatedSecurityGroups(policies[0].Ip, securityGroupsNames, existSecurityGroups)
	if err != nil {
		logrus.Errorf("createPolicies getAutoCreatedSecurityGroups meet error=%v", err)
		return newSecurityGroups, err
	}
	logrus.Infof("createPolicies createdSecurityGroups=%v, autoCreatedStartIndex=%v", createdSecurityGroups, autoCreatedStartIndex)

	//计算已经存在的安全组中还能插入多少条
	for _, securityGroup := range createdSecurityGroups {
		freeNum, err := getSecurityGroupFreePolicyNum(providerParams, securityGroup, direction)
		if err != nil {
			logrus.Errorf("createPolicies getSecurityGroupFreePolicyNum meet error=%v", err)
			return newSecurityGroups, err
		}
		freePolicyNumMap[securityGroup] = freeNum
		freePoliciesNum += freeNum
	}
	securityGroupsIds = append(securityGroupsIds, createdSecurityGroups...)

	//计算需要新创建几个安全组
	if freePoliciesNum < len(policies) {
		newSecurityGroupNum := (len(policies) - freePoliciesNum + MAX_SEUCRITY_RULE_NUM - 1) / MAX_SEUCRITY_RULE_NUM
		newSecurityGroups, err = createNewAutomationSecurityGroups(providerParams, policies[0].Ip, newSecurityGroupNum, autoCreatedStartIndex)
		if err != nil {
			logrus.Errorf("createPolicies createNewAutomationSecurityGroups meet error=%v", err)
			return newSecurityGroups, err
		}
		logrus.Infof("createPolicies newSecurityGroups=%v", newSecurityGroups)
		securityGroupsIds = append(securityGroupsIds, newSecurityGroups...)

		for _, securityGroup := range newSecurityGroups {
			freePolicyNumMap[securityGroup] = MAX_SEUCRITY_RULE_NUM
		}
	}

	logrus.Infof("createPolicies freePolicyNumMap=%v", freePolicyNumMap)
	//开始将策略加到安全组中
	offset, limit := 0, 0

	//for securityGroup, freeNum := range freePolicyNumMap {
	for _, securityGroupId := range securityGroupsIds {
		freeNum := freePolicyNumMap[securityGroupId]
		if len(policies)-offset > freeNum {
			limit = freeNum
		} else {
			limit = len(policies) - offset
		}
		if err := addPoliciesToSecurityGroup(providerParams, securityGroupId, policies[offset:offset+limit], direction); err != nil {
			logrus.Errorf("createPolicies addPoliciesToSecurityGroup meet error=%v", err)
			return newSecurityGroups, err
		}

		for i := offset; i < offset+limit; i++ {
			policies[i].SecurityGroupId = securityGroupId
		}
		offset += limit
	}

	logrus.Infof("createPolicies: return newSecurityGroups=%++v", newSecurityGroups)
	return newSecurityGroups, nil
}

func destroyPolicies(providerParams string, policies []*SecurityPolicy, direction string) error {
	logrus.Infof("destroyPolicies: input policies=%++v direction=%v", policies, direction)

	securityGroupMap := make(map[string][]*SecurityPolicy)
	for _, policy := range policies {
		securityGroupMap[policy.SecurityGroupId] = append(securityGroupMap[policy.SecurityGroupId], policy)
		logrus.Infof("destroyPolicies policy=%++v", *policy)
	}

	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	client, err := plugins.CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		logrus.Errorf("destroyPolicies CreateVpcClient meet error=%v", err)
		return err
	}

	for securityGroupId, policies := range securityGroupMap {
		securityGroupPolicySet := newSecurityPolicySet(policies, direction, false)
		req := vpc.NewDeleteSecurityGroupPoliciesRequest()
		req.SecurityGroupId = &securityGroupId
		req.SecurityGroupPolicySet = &securityGroupPolicySet

		_, err := client.DeleteSecurityGroupPolicies(req)
		if err != nil {
			logrus.Errorf("destroyPolicies DeleteSecurityGroupPolicies meet error=%v,req=%++v", err, *req)
			return err
		}
	}

	return nil
}
