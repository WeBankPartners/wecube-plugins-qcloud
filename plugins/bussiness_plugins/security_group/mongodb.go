package securitygroup

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20180408"
)

type MongodbResourceType struct {
}

type MongodbInstance struct {
	Id     string
	Name   string
	Region string
	Vip    string
}

func createMongodbClient(providerParams string) (client *mongodb.Client, err error) {
	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		return nil, err
	}

	credential := common.NewCredential(paramsMap["SecretID"], paramsMap["SecretKey"])
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "mongodb.tencentcloudapi.com"

	return mongodb.NewClient(credential, paramsMap["Region"], clientProfile)
}

func (resourceType *MongodbResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)
	client, _ := createMongodbClient(providerParams)
	var offset, limit uint64 = 0, uint64(len(instanceIds))
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	request := mongodb.NewDescribeDBInstancesRequest()
	request.InstanceIds = common.StringPtrs(instanceIds)
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeDBInstances(request)
	if err != nil {
		return result, err
	}

	if *resp.Response.TotalCount == 0 {
		return result, nil
	}

	for _, mongodb := range resp.Response.InstanceDetails {
		instance := MongodbInstance{
			Id:     *mongodb.InstanceId,
			Name:   *mongodb.InstanceName,
			Region: region,
			Vip:    *mongodb.Vip,
		}
		result[*mongodb.InstanceId] = instance
	}
	return result, nil
}

func queryMongodbInstances(providerParams string, offset uint64, limit uint64) ([]*mongodb.MongoDBInstanceDetail, uint64, error) {
	client, _ := createMongodbClient(providerParams)
	result := []*mongodb.MongoDBInstanceDetail{}
	request := mongodb.NewDescribeDBInstancesRequest()
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeDBInstances(request)
	if err != nil {
		return result, 0, err
	}

	if *resp.Response.TotalCount == 0 {
		return result, 0, nil
	}

	return resp.Response.InstanceDetails, *resp.Response.TotalCount, nil
}

func (resourceType *MongodbResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	var offset, limit uint64 = 0, 100
	result := make(map[string]ResourceInstance)
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	for {
		mongodbs, total, err := queryMongodbInstances(providerParams, offset, limit)
		if err != nil {
			return result, err
		}

		for _, db := range mongodbs {
			for _, ip := range ips {
				if ip == *db.Vip {
					instance := MongodbInstance{
						Id:     *db.InstanceId,
						Name:   *db.InstanceName,
						Region: region,
						Vip:    *db.Vip,
					}
					result[ip] = instance
					break
				}
			}
		}
		if total > offset+limit {
			offset = offset + limit
		} else {
			break
		}
	}

	return result, nil
}

func (resourceType *MongodbResourceType) IsLoadBalanceType() bool {
	return false
}

func (resourceType *MongodbResourceType) IsSupportEgressPolicy() bool {
	return false
}

func (instance MongodbInstance) ResourceTypeName() string {
	return "mongodb"
}

func (instance MongodbInstance) GetId() string {
	return instance.Id
}

func (instance MongodbInstance) GetName() string {
	return instance.Name
}

func (instance MongodbInstance) GetRegion() string {
	return instance.Region
}
func (instance MongodbInstance) GetIp() string {
	return instance.Vip
}

func (instance MongodbInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	return []string{}, fmt.Errorf("mongodb do not support query security group api")
}

func (instance MongodbInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	return fmt.Errorf("mongodb do not support associateSecurityGroup api")
}

func (instance MongodbInstance) IsSupportSecurityGroupApi() bool {
	return false
}

func (instance MongodbInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error) {
	return []ResourceInstance{}, []string{}, fmt.Errorf("mongodb do not support backendTarget")
}
