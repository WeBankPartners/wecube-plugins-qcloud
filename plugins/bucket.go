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
	BucketActions["create"] = new(BucketActionsCreateAction)
	BucketActions["delete"] = new(BucketActionsDeleteAction)
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

type BucketActionsCreateAction struct {
}

type BucketActionsDeleteAction struct {
}

func (action *BucketActionsCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs BucketInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *BucketActionsDeleteAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs BucketInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *BucketActionsCreateAction) createBucket(bucketInput *BucketInput) (output BucketOutput, err error) {
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
	u, _ := url.Parse(fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", bucketInput.BucketName, bucketInput.AccountAppId, paramsMap["Region"]))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  paramsMap["SecretID"],
			SecretKey: paramsMap["SecretKey"],
		},
	})
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
	logrus.Printf("create bucket:%s success,url:%s \n", bucketInput.BucketName, u.String())
	output.BucketUrl = u.String()
	return output, err
}

func (action *BucketActionsCreateAction) Do(input interface{}) (interface{}, error) {
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

func (action *BucketActionsDeleteAction) deleteBucket(bucketInput *BucketInput) (output BucketOutput, err error) {
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
	u, _ := url.Parse(fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", bucketInput.BucketName, bucketInput.AccountAppId, paramsMap["Region"]))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  paramsMap["SecretID"],
			SecretKey: paramsMap["SecretKey"],
		},
	})
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

func (action *BucketActionsDeleteAction) Do(input interface{}) (interface{}, error) {
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