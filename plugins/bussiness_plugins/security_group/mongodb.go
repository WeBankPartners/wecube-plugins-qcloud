package securitygroup

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
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
		logrus.Errorf("createBmClient: failed to create Qcloud mongodb client, err=%v", err)
		return nil, err
	}

	credential := common.NewCredential(paramsMap["SecretID"], paramsMap["SecretKey"])
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "mongodb.tencentcloudapi.com"

	return mongodb.NewClient(credential, paramsMap["Region"], clientProfile)
}

func (resourceType *MongodbResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	logrus.Infof("MongodbResourceType QueryInstancesById: request instanceIds=%++v", instanceIds)

	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		err := fmt.Errorf("instanceIds is empty")

		logrus.Errorf("MongodbResourceType QueryInstancesById meet error=%v", err)
		return result, err
	}

	client, _ := createMongodbClient(providerParams)
	var offset, limit uint64 = 0, uint64(len(instanceIds))
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	request := mongodb.NewDescribeDBInstancesRequest()
	request.InstanceIds = common.StringPtrs(instanceIds)
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeDBInstances(request)
	if err != nil {
		logrus.Errorf("MongodbResourceType QueryInstancesById DescribeDBInstances meet error=%v", err)
		return result, err
	}

	if *resp.Response.TotalCount == 0 {
		logrus.Infof("MongodbResourceType QueryInstancesById DescribeDBInstances: Response.TotalCount==0")
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

	logrus.Infof("MongodbResourceType QueryInstancesById: result=%++v", result)
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
		logrus.Errorf("queryMongodbInstances DescribeDBInstances meet error=%v", err)
		return result, 0, err
	}

	if *resp.Response.TotalCount == 0 {
		logrus.Infof("queryMongodbInstances DescribeDBInstances: Response.TotalCount==0")
		return result, 0, nil
	}

	logrus.Infof("queryMongodbInstances: return Response.InstanceDetails=%++v", resp.Response.InstanceDetails)
	return resp.Response.InstanceDetails, *resp.Response.TotalCount, nil
}

func (resourceType *MongodbResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	logrus.Infof("MongodbResourceType QueryInstancesByIp: request ips=%++v", ips)

	var offset, limit uint64 = 0, 100
	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		err := fmt.Errorf("ips is empty")

		logrus.Errorf("MongodbResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}

	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	for {
		mongodbs, total, err := queryMongodbInstances(providerParams, offset, limit)
		if err != nil {
			logrus.Errorf("MongodbResourceType queryMongodbInstances meet error=%v", err)
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
					result[*db.Vip] = instance
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

	logrus.Infof("MongodbResourceType: result=%++v", result)
	return result, nil
}

func (resourceType *MongodbResourceType) IsLoadBalanceType() bool {
	logrus.Infof("MongodbResourceType IsLoadBalanceType: return=[false]")
	return false
}

func (resourceType *MongodbResourceType) IsSupportEgressPolicy() bool {
	logrus.Infof("MongodbResourceType IsSupportEgressPolicy: return=[false]")
	return false
}

func (instance MongodbInstance) ResourceTypeName() string {
	logrus.Infof("MongodbInstance ResourceTypeName: return=[mongodb]")
	return "mongodb"
}

func (instance MongodbInstance) GetId() string {
	logrus.Infof("MongodbInstance GetId: return=[%v]", instance.Id)
	return instance.Id
}

func (instance MongodbInstance) GetName() string {
	logrus.Infof("MongodbInstance GetName: return=[%v]", instance.Name)
	return instance.Name
}

func (instance MongodbInstance) GetRegion() string {
	logrus.Infof("MongodbInstance GetRegion: return=[%v]", instance.Region)
	return instance.Region
}
func (instance MongodbInstance) GetIp() string {
	logrus.Infof("MongodbInstance GetIp: return=[%v]", instance.Vip)
	return instance.Vip
}

func (instance MongodbInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	err := fmt.Errorf("mongodb do not support query security group api")

	logrus.Errorf("MongodbInstance QuerySecurityGroups meet error=%v", err)
	return []string{}, err
}

func (instance MongodbInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	err := fmt.Errorf("mongodb do not support associateSecurityGroup api")

	logrus.Errorf("MongodbInstance AssociateSecurityGroups meet error=%v", err)
	return err
}

func (instance MongodbInstance) IsSupportSecurityGroupApi() bool {
	logrus.Infof("MongodbInstance IsSupportSecurityGroupApi: return=[false]")
	return false
}

func (instance MongodbInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error) {
	err := fmt.Errorf("mongodb do not support backendTarget")

	logrus.Errorf("MongodbInstance GetBackendTargets meet error=%v", err)
	return []ResourceInstance{}, []string{}, err
}
