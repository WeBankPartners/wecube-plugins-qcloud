package plugins

import (
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/url"
	"net/http"
	"context"
	"strings"
)

type BucketPlugin struct {
}

func (plugin *BucketPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := BucketActions[actionName]
	if !found {
		return nil, fmt.Errorf("Bucket plugin,action = %s not found", actionName)
	}

	return action, nil
}

var BucketActions = make(map[string]Action)

func init() {
	BucketActions["create"] = new(BucketCreateAction)
	BucketActions["delete"] = new(BucketDeleteAction)
}

type BucketInputs struct {
	Inputs []BucketInput `json:"inputs,omitempty"`
}

type BucketInput struct {
	CallBackParameter
	Guid             string `json:"guid,omitempty"`
	BucketName       string `json:"bucket_name,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	Location         string `json:"location"`
	APISecret        string `json:"api_secret"`
	AccountAppId     string `json:"account_app_id"`
	IsPublic         string `json:"is_public"`
	ForceDelete      string `json:"force_delete"`
}

type BucketOutputs struct {
	Outputs []BucketOutput `json:"outputs,omitempty"`
}

type BucketOutput struct {
	CallBackParameter
	Result
	RequestId    string `json:"request_id,omitempty"`
	Guid         string `json:"guid,omitempty"`
	BucketName   string `json:"bucket_name,omitempty"`
	BucketUrl    string `json:"bucket_url,omitempty"`
}

type BucketCreateAction struct {
}

type BucketDeleteAction struct {
}

func (action *BucketCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs BucketInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *BucketDeleteAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs BucketInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func getCosClient(name,appId,region,secretID,secretKey,bucketUrl string) (client *cos.Client,cosUrl string) {
	if bucketUrl == "" {
		bucketUrl = fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", name, appId, region)
	}
	u, _ := url.Parse(bucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return client,bucketUrl
}

func (action *BucketCreateAction) createBucket(bucketInput *BucketInput) (output BucketOutput, err error) {
	output.Guid = bucketInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = bucketInput.CallBackParameter.Parameter
	output.BucketName = bucketInput.BucketName

	if bucketInput.Location != "" && bucketInput.APISecret != "" {
		bucketInput.ProviderParams = fmt.Sprintf("%s;%s", bucketInput.Location, bucketInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(bucketInput.ProviderParams)
	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()
	client,bucketUrl := getCosClient(bucketInput.BucketName, bucketInput.AccountAppId, paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"], "")
	cosAcl := "private"
	isPublic := strings.ToLower(bucketInput.IsPublic)
	if isPublic == "y" || isPublic == "yes" || isPublic == "true" {
		cosAcl = "public-read"
	}
	opt := cos.BucketPutOptions{XCosACL:cosAcl}
	_,err = client.Bucket.Put(context.Background(), &opt)
	if err != nil {
		err = fmt.Errorf("create bucket:%s error ---> %v", bucketInput.BucketName, err)
		return output, err
	}
	logrus.Printf("create bucket:%s success,url:%s \n", bucketInput.BucketName, bucketUrl)
	output.BucketUrl = bucketUrl
	return output, err
}

func (action *BucketCreateAction) Do(input interface{}) (interface{}, error) {
	buckets, _ := input.(BucketInputs)
	outputs := BucketOutputs{}
	var finalErr error

	for _, bucket := range buckets.Inputs {
		bucketOutput, err := action.createBucket(&bucket)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, bucketOutput)
	}

	logrus.Infof("all buckets = %v are created", buckets)
	return &outputs, finalErr
}

func (action *BucketDeleteAction) deleteBucket(bucketInput *BucketInput) (output BucketOutput, err error) {
	output.Guid = bucketInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = bucketInput.CallBackParameter.Parameter
	output.BucketName = bucketInput.BucketName

	if bucketInput.Location != "" && bucketInput.APISecret != "" {
		bucketInput.ProviderParams = fmt.Sprintf("%s;%s", bucketInput.Location, bucketInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(bucketInput.ProviderParams)
	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()
	client,_ := getCosClient(bucketInput.BucketName, bucketInput.AccountAppId, paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"], "")
	// force
	forceDelete := strings.ToLower(bucketInput.ForceDelete)
	if forceDelete == "y" || forceDelete == "yes" || forceDelete == "true" {
		getResult,_,getErr := client.Bucket.Get(context.Background(), &cos.BucketGetOptions{MaxKeys:1000})
		if getErr != nil {
			err = fmt.Errorf("force delete bucket:%s fail, get bucket objects error ---> %v ", bucketInput.BucketName, err)
			return output,err
		}
		if len(getResult.Contents) > 0 {
			var tmpObjects []cos.Object
			for i,v := range getResult.Contents {
				logrus.Printf("contents %d: %s \n", i, v.Key)
				tmpObjects = append(tmpObjects, cos.Object{Key:v.Key})
			}
			delOpt := &cos.ObjectDeleteMultiOptions{
				Objects: tmpObjects,
			}
			_, _, err = client.Object.DeleteMulti(context.Background(), delOpt)
			if err != nil {
				err = fmt.Errorf("force delete bucket:%s fail, delete objects error ---> %v ", bucketInput.BucketName, err)
				return output, err
			}
		}
	}
	_,err = client.Bucket.Delete(context.Background())
	if err != nil {
		err = fmt.Errorf("delete bucket:%s error ---> %v ", bucketInput.BucketName, err)
		return output, err
	}
	return output, err
}

func (action *BucketDeleteAction) Do(input interface{}) (interface{}, error) {
	buckets, _ := input.(BucketInputs)
	outputs := BucketOutputs{}
	var finalErr error

	for _, bucket := range buckets.Inputs {
		bucketOutput, err := action.deleteBucket(&bucket)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, bucketOutput)
	}

	logrus.Infof("all buckets = %v are delete", buckets)
	return &outputs, finalErr
}

func SetBucketAcl(region,secretID,secretKey,bucketUrl,uin,permission string) error {
	client,_ := getCosClient("","",region,secretID,secretKey,bucketUrl)
	bucketAclResult,_,err := client.Bucket.GetACL(context.Background())
	if err != nil {
		return fmt.Errorf("get bucket owner id fail,error: %v ", err)
	}
	var readGrant,writeGrant,fullControlGrant string
	if len(bucketAclResult.AccessControlList) > 0 {
		userClient,_ := createUserClient(region,secretID,secretKey)
		users,_ := ListSubUsers(userClient)
		for _,v := range bucketAclResult.AccessControlList {
			logrus.Infof("access control ---> permission:%s id:%s type:%s ", v.Permission, v.Grantee.ID, v.Grantee.Type)
			tmpGranteeId := v.Grantee.ID
			var newList []string
			for _,vv := range strings.Split(tmpGranteeId, ",") {
				if strings.Contains(vv, uin) {
					continue
				}
				userExist := false
				for _,vvv := range users {
					if strings.Contains(vv, vvv.Uin) {
						userExist = true
						break
					}
				}
				if !userExist {
					logrus.Infof("user uin:%s this user not exist ")
					continue
				}
				newList = append(newList, fmt.Sprintf("id=\"%s\"", vv))
			}
			tmpGranteeId = strings.Join(newList, ",")
			switch strings.ToLower(v.Permission) {
				case "read": readGrant = tmpGranteeId
				case "write": writeGrant = tmpGranteeId
				case "full_control": fullControlGrant = tmpGranteeId
			}
		}
	}
	ownerUin := bucketAclResult.Owner.ID[strings.LastIndex(bucketAclResult.Owner.ID, "/")+1:]
	grantId := fmt.Sprintf("id=\"qcs::cam::uin/%s:uin/%s\"", ownerUin, uin)
	switch permission {
		case "read": readGrant = appendGrantId(readGrant, grantId)
		case "write": writeGrant = appendGrantId(writeGrant, grantId)
		case "full_control": fullControlGrant = appendGrantId(fullControlGrant, grantId)
	}
	logrus.Infof("--------> read:%s  write:%s  full_control:%s", readGrant, writeGrant, fullControlGrant)
	opt := &cos.BucketPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosGrantRead: readGrant,
			XCosGrantWrite: writeGrant,
			XCosGrantFullControl: fullControlGrant,
		},
	}
	_,err = client.Bucket.PutACL(context.Background(), opt)
	if err != nil {
		logrus.Errorf("set bucket acl with grant:%s error %v ", grantId, err)
	}
	return err
}

func appendGrantId(old,grant string) string {
	if old == "" {
		return grant
	}
	return fmt.Sprintf("%s,%s", old, grant)
}