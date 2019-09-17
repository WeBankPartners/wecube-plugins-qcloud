package securitygroup
import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"time"
	"strings"
	"strconv"
	"sync"
	//vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)
const (
	MAX_SEUCRITY_RULE_NUM = 100
)

//interface definition
type ResourceInstance interface {
	ResourceTypeName()string
	Id()string
	Name()string
	Region()string
	QuerySecurityGroups(providerParams string)([]string,error)
	AssociateSecurityGroups(providerParams string,securityGroups []string)error
}

type ResourceType interface {
	QueryInstancesById(providerParams string,instanceIds []string)(map[string]ResourceInstance,error)
	QueryInstancesByIp(providerParams string,ips []string)(map[string]ResourceInstance,error)
	IsSupportSecurityGroupApi()bool 
}

//resourceType register
var (
	resTypesMutex sync.Mutex
	resourceTypeMap = map[string]ResourceType
)

//resource type  register 
func addNewResourceType(name string, newResourceType ResourceType) error{
	resTypesMutex.Lock()
	defer resTypesMutex.Unlock()

	if _, found := resourceTypeMap[name]; found {
		logrus.Fatalf("resourceType(%s) was registered twice", name)
	}

	resourceTypeMap[name] = newResourceType
}

func getResouceTypeByName(name string) (ResourceType, error) {
	resTypesMutex.Lock()
	defer resTypesMutex.Unlock()

	resType, found := resourceTypeMap[name]
	if !found {
		return nil, fmt.Errorf("resourceType[%s] not found", name)
	}
	return resType, nil
}

type BussinessSecurityGroupPlugin struct{
}

var SecurityGroupActions = make(map[string]plugins.Action)

func init() {
	//plugin registry
	RegisterPlugin("bs-security-group", new(BussinessSecurityGroupPlugin))

	//resourceType registry
	addNewResourceType("mysql",new(MysqlResourceType))

	//action 
	SecurityGroupActions["calc-security-policies"] = new(CalcSecurityPolicyAction)
	SecurityGroupActions["apply-security-policies"] = new(ApplySecurityPolicyAction)
}

func findInstanceByIp(ip string)(ResourceInstance,error){
	regions,err := getRegions()
	if err != nil {
		logrus.Errorf("getRegions meet err=%v\n",err)
		return nil,err
	}

	for _,region:=range regions {
		providerParams,err := getProviderParams(region)
		if err != nil {
			logrus.Fatalf("getProviderParams meet err=%v\n",err)
			return nil,err
		}

		for _,resType:=range resourceTypeList{
			instanceMap,err:=resType.QueryInstancesByIp(providerParams,[]string{ip})
			if err != nil {
				logrus.Errorf("QueryInstancesByIp meet err=%v\n",err)
			    return nil,err
			}
			instance,ok :=instanceMap[ip]
			if ok {
				return instance,nil 
			}
		}
	}

	logrus.Errorf("ip(%s),can't be found",ip)
	return nil,fmt.Errorf("ip(%s),can't be found",ip)
}

//---------------calc security policy action------------------------------//
type CalcSecurityPoliciesRequest struct {
	 Protocol   string `json:"protocol"`
	 SourceIps  []string `json:"source_ips"`
	 DestIps    []string `json:"dest_ips"`
	 DestPort   string `json:"dest_port"`
	 PolicyAction string `json:"policy_action"`
	 PolicyDirections []string `json:"policy_directions"`
	 Description  string `json:"description"`
}
			   
type SecurityPolicy struct {
	Ip string `json:"ip"`
	Type string `json:"type"`
	Id   string `json:"id"`
	Region string `json:"region"`
	SupportSecurityGroupApi bool `json:"support_security_group_api"`
	PeerIp string `json:"peer_ip"`
	Protocol string `json:"protocol"`
	Ports   string  `json:"ports"`
	Action string    `json:"action"`
	Description  string `json:"description"`
	ErrorMsg  string `json:"err_msg,omitempty"`
	UndoReason  string `json:"undo_reason,omitempty"`

	SecurityGroupId string `json:"-"`
}

type CalcSecurityPoliciesResult struct {
	TimeTaken  string `json:"time_taken"`

	IngressPoliciesTotal int `json:"ingress_policies_total"`
	EgressPoliciesTotal int `json:"egress_policies_total"`

	IngressPolicies []SecurityPolicy `json:"ingress_policies,omitempty"`
	EgressPolicies []SecurityPolicy `json:"egress_policies,omitempty"`
}

type CalcSecurityPolicyAction struct {
}

func (action *CalcSecurityPolicyAction) ReadParam(param interface{}) (interface{}, error) {
	var input CalcSecurityPoliciesRequest
	err := plugins.UnmarshalJson(param, &input)
	if err != nil {
		logrus.Errorf("CalcSecurityPolicyAction unmarshal failed,err=%v,param=%v",err,param)
		return nil, err
	}
	return input, nil
}

func (action *CalcSecurityPolicyAction) CheckParam(input interface{}) error {
	req, _ := input.(CalcSecurityPoliciesRequest)
	if err:=isValidProtocol(req.Protocol);err != nil{
		return err 
	}

	if err := isValidAction();err != nil {
		return err 
	}
	
	for _,ip:=range req.SourceIps {
		if err := isValidIp(ip);err!=nil {
			return err 
		}
	}

	for _,ip:=range req.DestIps {
		if err := isValidIp(ip);err!=nil {
			return err 
		}
	}

	_,err := getPortsByPolicyFormat(req.DestPort)
	if err != nil {
		return err 
	}

	for _,direction:=range req.PolicyDirections{
		if err := isValidDirection(direction);err!=nil {
			return err 
		}
	}

	return nil
}

func calcPolicies(devIp string,peerIps []string,proto string,ports []string,
	             action string,description)([]SecurityPolicy,error) {
	var finalError error 
	policies := []SecurityPolicy{}
	joinPorts:=strings.Join(ports,",")

	//check if dev exist 
	instance,err:=findInstanceByIp(devIp)
	if err != nil {
		return policies ,err
	}

	restType,err:=getResouceTypeByName(instanace.ResourceTypeName())
	if err != nil {
		return policies ,err
	}

	for _, peerIp:=range peerIps{
		for_,port:=range ports{
			newPolicy:=SecurityPolicy{
				Ip:devIp,
				Type:instance.ResourceTypeName(),
				Id:instance.Id(),
				Region:instance.Region(),
				SupportSecurityGroupApi:restType.IsSupportSecurityGroupApi(),
				PeerIp:peerIp,
				Protocol:proto,
				Ports:joinPorts,
				Action:action,
				Description:description,
			}
			policies = append(policies,newPolicy)
		}
	}
	return policies,nil
}

func (action *CalcSecurityPolicyAction) Do(input interface{}) (interface{}, error) {
	var finalError error
	req, _ := input.(CalcSecurityPoliciesRequest)
	result:=CalcSecurityPoliciesResult{}
	start:=time.Now()
	ports,_:=getPortsByPolicyFormat(req.DestPort)

	//calc egress policies
	if isContainInList(EGRESS_RULE,req.PolicyDirections){
		for _,ip:=range req.SourceIps{
			policies,err:= calcPolicies(ip,req.DestIps,req.Protocol,ports,req.Action,req.Description)
			result.EgressPolicies=append(result.EgressPolicies,policies...)
			if err != nil {
				finalError = fmt.Errorf(%s,finalError.Error()+err.Error())
			}
		}
	}

	//calc ingress policies
	if isContainInList(INGRESS_RULE,req.PolicyDirections){
		for _,ip := range req.DestIps{
			policies,err := calcPolicies(ip,req.SourceIps,req.Protocol,ports,req.Action,req.Description)
			result.IngressPolicies=append(result.EgressPolicies,policies...)
			if err != nil {
				finalError = fmt.Errorf(%s,finalError.Error() + err.Error())
			}
		}
	}

	result.TimeTaken=fmt.Sprintf("%v",time.Since(start))
	result.IngressPoliciesTotal=len(result.IngressPolicies)
	result.EgressPoliciesTotal=len(result.EgressPolicies)

	return result,finalError 
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
	UndoPolicies []SecurityPolicy `json:"undo_policies"`
	FailedPolicies []SecurityPolicy `json:"failed_policies"`
}

type ApplySecurityPoliciesResult struct {
	TimeTaken  string `json:"time_taken"`
	IngressApplyResult ApplyResult `json:"ingress"`
	EgressApplyResult  ApplyResult `json:"egress"`
}

func (action *ApplySecurityPolicyAction) ReadParam(param interface{}) (interface{}, error) {
	var input ApplySecurityPoliciesRequest
	err := plugins.UnmarshalJson(param, &input)
	if err != nil {
		logrus.Errorf("ApplySecurityPolicyAction:unmarshal failed,err=%v,param=%v",err,param)
		return nil, err
	}
	return input, nil
}

func (action *ApplySecurityPolicyAction) CheckParam(input interface{}) error {
	req, _ := input.(ApplySecurityPoliciesRequest)

	for _,policy:=req.IngressPolicies{
		if policy.Ip == "" || policy.InstanceId =="" {
			return errors.New("ingress policy have empty value")
		}
	}

	for _,policy:=req.EgressPolicies{
		if policy.Ip == "" || policy.InstanceId =="" {
			return errors.New("egress policy have empty value")
		}
	}

	return nil 
}

func (action *ApplySecurityPolicyAction) Do(input interface{}) (interface{}, error){
	var err error
	req, _ := input.(ApplySecurityPoliciesRequest)
	result :=ApplySecurityPoliciesResult{}
	start :=time.Now()

	result.IngressApplyResult = applyPolicies(IngressPolicies []SecurityPolicy,INGRESS_RULE)
	result.EgressApplyResult = applyPolicies(IngressPolicies []SecurityPolicy,EGRESS_RULE)

	result.Taken = fmt.Sprintf("%v",time.Since(start))
	if result.IngressResult.FailedTotal > 0 || result.EgressApplyResult.FailedTotal > 0{
		err=errors.New("have some failed polices,please check policy applied detail")
	}

	return result,err
}

func fillSecuityPoliciesWithErrMsg(policies []*SecurityPolicy,err error ){
	for _,policy:=range policies{
		policy.ErrorMsg= err.Error()
	}
}

func applyPolicies(policies []SecurityPolicy,direction string)ApplyResult {
	result:=ApplyResult{}
	instanceMap:=make(map[string][]*SecurityPolicy)

	for i,_ := range policies{
		if  policies[i].SupportSecurityGroupApi == true {
			key:=policies[i].Ip
			instanceMap[key]=append(instanceMap[key],&policies[i])
		}else {
			policy.UndoReason=fmt.Sprintf("instanceType(%s) do not support security group api",policy.Type)
			result.UndoPolicies=append(result.UndoPolicies,policy)
		}
	}

	for _,policies := range instanceMap{
		resType,err := getResouceTypeByName(policies[0].Type)
		if err != nil {
			fillSecuityPoliciesWithErrMsg(policies,err)
			continue
		}
		
		providerParams,err:=getProviderParams(policies[0].Region)
		if err != nil {
			fillSecuityPoliciesWithErrMsg(policies,err)
			continue
		}
		
		instances,err:=resType.QueryInstancesById(providerParams,[]string{policies[0].InstanceId})
		if err != nil {
			fillSecuityPoliciesWithErrMsg(policies,err)
			continue
		}
		if len(instances) == 0 {
			fillSecuityPoliciesWithErrMsg(policies,errors.New("can't found instanceId(%s)",policies[0].InstanceId))
			continue
		}

		existSecurityGroups,err:=instances[0].QuerySecurityGroups()
		if err != nil {
			fillSecuityPoliciesWithErrMsg(policies,err)
			continue
		}
		
		newSecurityGroups,err:=createPolicies(providerParams,existSecurityGroups,policies,direction)
		if err != nil {
			destroyPolicies(providerParams,policies,direction)
			fillSecuityPoliciesWithErrMsg(policies,err)
			continue
		}

		if len(newSecurityGroups) > 0 {
			groups:=[]string{}
			groups=append(groups,newSecurityGroups...)
			groups=append(groups,existSecurityGroups...)

			if err = instances[0].AssociateSecurityGroups(providerParams,groups);err != nil {
				destroyPolicies(providerParams,policies,direction)
				bindError:=fmt.Errorf("resourceType(%s) instance(%s) AssociateSecurityGroups[%v] meet err=%v",policies[0].Type,policies[0].Ip,groups,err)
				fillSecuityPoliciesWithErrMsg(policies,bindError)
				continue
			}
		}
	}

	for _,policies := range instanceMap{
		for _,policy := range policies {
			if policy.ErrorMsg == "" {
				result.SuccessPolicies = append(result.SuccessPolicies,*policy)
			}else {
				result.FailedPolicies = append(result.FailedPolicies,*policy)
			}
		}
    }
	result.PoliciesTotal  = len(policies)
	result.SuccessTotal = len(result.SuccessPolicies)
	result.FaildTotal = len(result.FailedPolicies)

	return result
}

//自动构建的安全组的名称格式ip-auoto-1,ip_auto_2
func getAutoCreatedSecurityGroups(ip string,allSecurityGroups []string)([]string,int){
	 maxAutoCreatedNum := 0
	 createdSecurityGroups:=[]string{}
	 for _,securityGroup:=range allSecurityGroups{
		 elements:=strings.Split(createdSecurityGroups,"-")
		 if len(elements) ==3 {
			 if elements[0]==ip  && elements[1]="auto" {
				 if num,err:=strconv.Atoi(elements[3]);err == nil {
					createdSecurityGroups=append(createdSecurityGroups,securityGroup)
					maxAutoCreatedNum = num
				 }
			 }
		 }
	 }

	 return createdSecurityGroups,maxAutoCreatedNum + 1
}

func getSecurityGroupFreePolicyNum(providerParams string,securityGroup string,direction string)(int,error){
	policiesSet,err:=plugins.QuerySecurityGroupPolicies(providerParams,securityGroup)
	if err != nil {
		logrus.Errorf("getSecurityGroupFreePolicyNum meet err=%v\n",err)
		return 0,err 
	}

	if strings.EqualFold(direction,INGRESS_RULE){
		return MAX_SEUCRITY_RULE_NUM - len(*policiesSet.Ingress)
	}

	return  MAX_SEUCRITY_RULE_NUM - len(*policiesSet.Egress)
}

func  getSecurityGroupNames(providerParams string,securityGroupIds []string)([]string,error){
	securityGroupNames :=[]string{}
	idNameMap :=make (map[string]string)
	securityGroupSet,err := plugins.QuerySecurityGroups(providerParams,securityGroupIds)
	if err != nil {
		return securityGroupNames,err
	}

    for securityGroup:=range securityGroupSet {
		idNameMap[securityGroup.SecurityGroupId]=securityGroup.SecurityGroupName
	}
	
	for _,id:=range securityGroupIds{
		if name, ok :=idNameMap[id];ok {
			securityGroupNames:=append(securityGroupNames,name)
		}else {
			return securityGroupNames,fmt.Errorf("can't found groupId(%s) detail",id)
		}
	}

	return securityGroupNames,nil 
}

//format ip-auto-2
func createNewAutomationSecurityGroups(providerParams string,ip string,newCreatedSecurityGroupNum int,auotNumIndex int)([]string,error){
   newSecurityGroupIds:=[]string{}
   for i:=0;i<newCreatedSecurityGroupNum;i++ {
	   securityGroupName:=fmt.Sprintf("%s-auto-%d,",ip,auotNumIndex+i)
	   securityGroupId,err:=plugins.CreateSecurityGroup(providerParam ,securityGroupName,"automation created")
	   if err != nil {
		   logrus.Errorf("CreateSecurityGroup meet err=%v",err)
		   return newSecurityGroupIds,err
	   }
	   newSecurityGroupIds=append(newSecurityGroupIds,securityGroupId)
   }

   return newSecurityGroupIds,nil 
}
func newSecurityPolicySet(policies []*SecurityPolicy,direction string)vpc.SecurityGroupPolicySet{
	securityPolicies:=[]*vpc.SecurityGroupPolicy
	for _,policy:=range policies {
		action:=strings.ToUpper(policy.Action)
		securityPolicy:=vpc.SecurityGroupPolicy{
			PolicyIndex:&policyIndex,
			Protocol:&policy.Protocol,
			Port:&policy.Ports,
			CidrBlock:&policy.PeerIp,
			Action:&action,
			PolicyDescription:&policy.Description,
		}
		securityPolicies=append(securityPolicies,&securityPolicy)
	}

	securityGroupPolicySet=vpc.SecurityGroupPolicySet{}
	if strings.EqualFold(direction,INGRESS_RULE){
		securityGroupPolicySet.Ingress = securityPolicies
	}else {
		securityGroupPolicySet.Egress = securityPolicies
	}

	return securityGroupPolicySet
}

func addPoliciesToSecurityGroup(providerParams string,securityGroupId string,policies []*SecurityPolicy,direction string)error{
	req := vpc.NewCreateSecurityGroupPoliciesRequest()
	req.SecurityGroupId = &securityGroupId
	var policyIndex int64 = 0
	var err error 
	securityPolicies:=newSecurityPolicySet(policies)

	if len(policies) == 0  {
		return nil 
	}
	defer func(){
		if err != nil {
				logrus.Errorf("add policy to securityGroup(%s) meet err =%v",securityGroupId,err)
				errMsg:=fmt.Sprintf("add policy to securityGroup(%s) meet err =%v",securityGroupId,err)
				for _,policy:=range policies {
					policy.ErrorMsg = errMsg
				}
		}
	}()

	paramsMap, err :=plugins.GetMapFromProviderParams(providerParams)
	client,err:=plugins.CreateVpcClientcreateVpcClient()
	if err != nil {
		return err
	}

	securityGroupPolicySet:=newSecurityPolicySet(policies,direction)
	req.SecurityGroupPolicySet = &securityGroupPolicySet
	if _, err = client.CreateSecurityGroupPolicies(req);err == nil{
		for _,policy:=range policies {
			policy.SecurityGroupId = securityGroupId
		}
	}

	return err 
}

func createPolicies(providerParams string,existSecurityGroups []string ,policies []*SecurityPolicy,direction string)([]string,error){
	newCreatedSecurityGroups:=[]string{}
	availablePolicyNumMap:=make(map[string]int)
	freePoliciesNum:=0
	
	if len(policies) == 0 {
		return newCreatedSecurityGroups,nil 
	}

	securityGroupsNames,err := getSecurityGroupNames(providerParams,existSecurityGroups)
	if err != nil {
		return newCreatedSecurityGroups,err
	}

	createdSecurityGroups,autoCreatedStartIndex := getAutoCreatedSecurityGroups(policies[0].Ip,securityGroupsNames)
	//计算已经存在的安全组中还能插入多少条
	for _,securityGroup := range createdSecurityGroups {
		freeNum,err:=getSecurityGroupFreePolicyNum(providerParams,securityGroup,direction)
		if err != nil {
			return newCreatedSecurityGroups,err
		}
		freePolicyNumMap[securityGroup] = freeNum
		freePoliciesNum += freeNum
	}

	//计算需要新创建几个安全组
	if (freePoliciesNum > len(policies)){
		needCreateSecurityGroupNum:= (len(policies) - freePoliciesNum + MAX_SEUCRITY_RULE_NUM -1)/MAX_SEUCRITY_RULE_NUM
		newSecurityGroups,err = createNewAutomationSecurityGroups(providerParams,policies[0].Ip,needCreateSecurityGroupNum,autoCreatedStartIndex)
		if err != nil {
			return newCreatedSecurityGroups,err
		}

		for _,securityGroup:=range newSecurityGroups{
			freePolicyNumMap[securityGroup] = MAX_SEUCRITY_RULE_NUM
		}
	}
	
	//开始将策略加到安全组中
	offset,limit :=0,0
	for securityGroup,freeNum:=range freePolicyNumMap{
		if len(policies) - offset >freeNum {
			limit = freeNum
		}else {
			limit = len(policies) - offset
		}

		if err:=addPoliciesToSecurityGroup(providerParams,securityGroup,policies[offset:offset+limit],direction);err != nil {
            return newCreatedSecurityGroups,err
		}

		for i:=offset;i<offset+limit;i++{
			policies[i].SecurityGroup = securityGroup
		}
		offset+=limit
	}
}

func destroyPolicies(providerParams string,policies []*SecurityPolicy,direction string){
	securityGroupMap:=make(map[string][]*SecurityPolicy)
	for _,policy:=range policies{
		securityGroupMap[policy.SecurityGroupId] = append(securityGroupMap[policy.SecurityGroupId],policy)
	}

	paramsMap, err :=plugins.GetMapFromProviderParams(providerParams)
	client,err:=plugins.CreateVpcClientcreateVpcClient()
	if err != nil {
		return err
	}

	for securityGroupId,policies:=range securityGroupMap{
		securityGroupPolicySet:=newSecurityPolicySet(policies,direction)
		req:=vpc.NewDeleteSecurityGroupPoliciesRequest() 
		req.SecurityGroupId =&securityGroupId
		req.SecurityGroupPolicySet = &securityGroupPolicySet

		_,err:=client.DeleteSecurityGroupPolicies(req)
		if err != nil {
			logrus.Errorf("DeleteSecurityGroupPolicies meet err=%v,req=%++v",err,*req)
		}
	}
}
