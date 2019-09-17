package securitygroup
import (
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
)

//resource type 
type MysqlResourceType  struct {

}

func (resourceType *MysqlResourceType)QueryInstancesById(providerParams string,instanceIds []string)(map[string]ResourceInstance,error){
	result:=make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		return result,nil 
	}

	filter:=plugins.Filter{
		Name:"instanceId",
		Values:instanceIds,
	}
	paramsMap, _:= plugins.GetMapFromProviderParams(providerParams)
	items,err:=plugins.QueryMysqlInstance(providerParams,filter)
	if err != nil {
		return result,err 
	}

	for _,item:=range items{
		instance:=MysqlInstance{
			Id: item.InstanceId,
	        Name:item.InstanceName,
			Vip:item.Vip,
			Region:paramsMap["Region"],
		}
		result[item.Vip] = instance
	}

	return result,nil 
}

func (resourceType *MysqlResourceType)QueryInstancesByIp(providerParams string,ips []string)(map[string]ResourceInstance,error){
	result:=make(map[string]ResourceInstance)

	if len(instanceIds) == 0 {
		return result,nil 
	}

	filter:=plugin.Filter{
		Name:"vip",
		Values:ips,
	}

	items,err:=plugins.QueryMysqlInstance(providerParams,filter)
	if err != nil {
		return result,err 
	}

	paramsMap, _:= plugins.GetMapFromProviderParams(providerParams)
	for _,item:=range items{
		instance:=MysqlInstance{
			Id: item.InstanceId,
	        Name:item.InstanceName,
			Vip:item.Vip,
			Region:paramsMap["Region"],
		}
		result[item.Vip] = &instance
	}

	return result,nil 
}

func (resourceType *MysqlResourceType) IsSupportSecurityGroupApi()bool {
	return true 
}


//resource instance 
type MysqlInstance struct {
	Id string
	Name string
	Vip string
	Region string 
}

func (instance *MysqlInstance)Id()string{
	return instance.Id
}

func (instance *MysqlInstance)Name()string{
	return instance.Name
}

func (instance *MysqlInstance)QuerySecurityGroups(providerParams string)([]string,error){
	return plugins.QueryMySqlInstanceSecurityGroups(providerParams,instance.Id)
}

func (inst *MysqlInstance)AssociateSecurityGroups(providerParams string,securityGroups []string)error{
	return plugins.BindMySqlInstanceSecurityGroups(providerParams,instance.Id,securityGroups)
}

func (inst *MysqlInstance)ResourceTypeName()string{
	return "mysql"
}
func (inst *MysqlInstance )Region()string{
	return instance.Region
}





