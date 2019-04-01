package cmdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"git.webank.io/wecube-plugins/conf"
	"github.com/sirupsen/logrus"
)

const (
	CMDB_OPERATE_CI_PATH        = "/cmdb/api/operateCi.json"
	CMDB_GET_TEMPLATE_DATA_PATH = "/cmdb/api/getIntegrateTemplateData.json"

	CMDB_STATE_CREATED    = "Created"
	CMDB_STATE_REGISTERED = "Registered"

	MAX_LIMIT_VALUE = 50000
)

type CiType int

const (
	CI_TYPE_NORMAL CiType = iota
	CI_TYPE_INTEGRATE
)

type CmdbCiQueryParam struct {
	Filter        map[string]string
	ResultColumn  []string
	Offset        int
	Limit         int
	OrderBy       string
	Order         string
	PluginCode    string
	PluginVersion string
}

type CmdbRequest struct {
	UserAuthKey           string                   `json:"userAuthKey"`
	Type                  string                   `json:"type,omitempty"`
	Action                string                   `json:"action,omitempty"`
	Filters               []map[string]interface{} `json:"filters,omitempty"`
	Filter                map[string]interface{}   `json:"filter,omitempty"`
	Parameters            []map[string]interface{} `json:"parameters,omitempty"`
	Parameter             map[string]interface{}   `json:"parameter,omitempty"`
	ResultColumn          []string                 `json:"resultColumn,omitempty"`
	IsPaging              bool                     `json:"isPaging,omitempy"`
	StartIndex            int                      `json:"startIndex,omitempy"`
	PageSize              int                      `json:"pageSize,omitempty"`
	Orderby               map[string]interface{}   `json:"orderby,omitempty"`
	EnableCasCadingDelete bool                     `json:"enableCascadingDelete,omitempty"`
	LimitRowCount         int                      `json:"limitRowCount,omitempty"`
	PluginCode            string                   `json:"pluginCode,omitempty"`
	PluginVersion         string                   `json:"pluginVersion,omitempty"`
}

type CmdbResponse struct {
	Headers CmdbHeaders `json:"headers"`
	Data    CmdbData    `json:"data"`
}

type CmdbHeaders struct {
	RetCode        int         `json:"retCode"`
	StartIndex     interface{} `json:"startIndex"`
	TotalRows      interface{} `json:"totalRows"`
	RetDetail      string      `json:"retDetail"`
	PageSize       interface{} `json:"pageSize"`
	Msg            string      `json:"msg"`
	PermissionType string      `json:"permissionType"`
	ContentRows    int         `json:"contentRows"`
}

type CmdbData struct {
	Header  []Header    `json:"header"`
	Content interface{} `json:"content"`
}

type Header struct {
	EnName          string      `json:"enName"`
	DataType        string      `json:"dataType"`
	DisplaySeq      int         `json:"display_seq"`
	Name            string      `json:"name"`
	SearchSeq       int         `json:"search_seq"`
	IDAdmCiTypeAttr int         `json:"idAdmCiTypeAttr"`
	DisplayType     interface{} `json:"display_type"`
	PermissionType  string      `json:"permissionType"`
	IsUnique        interface{} `json:"is_unique"`
	IsNone          interface{} `json:"is_none"`
	IsSystem        interface{} `json:"is_system"`
	IsSearch        interface{} `json:"is_search"`
	RefType         int         `json:"refType"`
	RefUrl          string      `json:"refUrl"`
	IsDisplay       string      `json:"is_display"`
	Description     string      `json:"description"`
}

type CmdbClient struct {
	HttpClient *http.Client
	Host       string
	Path       string
}

func NewCmdbClient(path string) *CmdbClient {
	client := &CmdbClient{
		HttpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		Host: conf.GobalAppConfig.CMDBLink,
		Path: path,
	}

	return client
}

func isValidPointer(response interface{}) error {
	if nil == response {
		return errors.New("input param should not be nil")
	}

	if kind := reflect.ValueOf(response).Type().Kind(); kind != reflect.Ptr {
		return errors.New("input param should be pointer type")
	}

	return nil
}

func (client *CmdbClient) DoPostHttpRequest(request interface{}, response interface{}) ([]byte, error) {
	if err := isValidPointer(response); err != nil {
		return nil, err
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		logrus.Errorf("Cmdb DoPostHttpRequest Marshal request failed err=%v", err)
		return nil, err
	}

	url := client.Host + client.Path
	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))
	if err != nil {
		logrus.Errorf("Cmdb DoPostHttpRequest NewRequest failed err=%v", err)
		return nil, err
	}

	logrus.Debugf("Http request url = %s,requestData = %s", url, string(requestBytes))

	httpResponse, err := client.HttpClient.Do(httpRequest)
	if err != nil {
		logrus.Errorf("Cmdb DoPostHttpRequest Do failed err=%v", err)
		return nil, err
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		logrus.Errorf("Cmdb DoPostHttpRequest httpResponse.StatusCode != 200,status=%v", httpResponse.StatusCode)
		return nil, fmt.Errorf("Cmdb DoPostHttpRequest httpResponse.StatusCode != 200,statusCode=%v", httpResponse.StatusCode)
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		logrus.Errorf("DoRequest failed!read http response failed: url=%s,err=%v", url, err)
		return body, err
	}

	logrus.Debugf("Http response, url =%s,response=%s", url, string(body))

	err = json.Unmarshal(body, response)
	if err != nil {
		logrus.Errorf("Cmdb DoPostHttpRequest unmarshal failed err=%v,body=%s", err, string(body))
	}

	return body, err
}

func UnmarshalContent(cmdbContent interface{}, destinationStruct interface{}) error {
	if err := isValidPointer(destinationStruct); err != nil {
		return err
	}

	body, err := json.Marshal(cmdbContent)
	if err != nil {
		logrus.Errorf("Cmdb UnmarshalContent marshal failed,err=%v ", err)
		return err
	}

	err = json.Unmarshal(body, destinationStruct)
	if err != nil {
		logrus.Errorf("Cmdb UnmarshalContent unmarshal failed,err=%v,body=%v ", err, string(body))
	}

	return err
}

func GetErrorFromResponse(resp *CmdbResponse) error {
	if resp.Headers.RetCode == 0 {
		return nil
	}
	logrus.Errorf("Cmdb GetErrorFromResponse meet error=%s", resp.Headers.Msg)
	return errors.New(resp.Headers.Msg)
}

func OperateCi(req *CmdbRequest) (*CmdbResponse, []byte, error) {
	resp := &CmdbResponse{}
	client := NewCmdbClient(CMDB_OPERATE_CI_PATH)

	if req.UserAuthKey == "" {
		req.UserAuthKey = conf.GobalAppConfig.CMDBUserAuthKey
	}

	buf, err := client.DoPostHttpRequest(req, resp)
	if err != nil {
		logrus.Errorf("Cmdb OperateCi meet error=%s,request=%v", err, req)
		return resp, buf, err
	}
	return resp, buf, GetErrorFromResponse(resp)
}

func getIntegrateTemplateData(req *CmdbRequest) (*CmdbResponse, []byte, error) {
	resp := &CmdbResponse{}
	client := NewCmdbClient(CMDB_GET_TEMPLATE_DATA_PATH)

	if req.UserAuthKey == "" {
		req.UserAuthKey = conf.GobalAppConfig.CMDBUserAuthKey
	}

	logrus.Infof("req=%v", req)
	buf, err := client.DoPostHttpRequest(req, resp)
	if err != nil {
		logrus.Errorf("Cmdb getIntegrateTemplateData meet error=%s,request=%v", err, req)
		return resp, buf, err
	}

	return resp, buf, GetErrorFromResponse(resp)
}

func GetStringFromInterface(v interface{}) string {
	data := fmt.Sprintf("%v", v)
	return data
}

func ExtractColumnFromStruct(s interface{}) []string {
	columns := []string{}
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			column := t.Field(i).Tag.Get("json")
			columns = append(columns, strings.Split(column, ",")[0])
		}
	}
	return columns
}

func GetIntFromInterface(v interface{}) (int, error) {
	data := fmt.Sprintf("%v", v)
	return strconv.Atoi(data)
}

func GetMapFromStruct(dataStruct interface{}) (map[string]interface{}, error) {
	var dataMap map[string]interface{}
	bytes, err := json.Marshal(dataStruct)
	if err != nil {
		logrus.Errorf("GetMapFromStruct meet error=%v", err)
		return dataMap, err
	}

	err = json.Unmarshal(bytes, &dataMap)
	if err != nil {
		logrus.Errorf("GetMapFromStruct meet Error,bytes=%s,err=%v", string(bytes), err)
	}
	return dataMap, err
}

func GetMapFromProviderParams(providerParams string) (map[string]string, error) {
	rtnMap := make(map[string]string)
	params := strings.Split(providerParams, ";")

	if len(params) == 0 {
		return rtnMap, nil
	}

	for _, param := range params {
		afterTrimParam := strings.Trim(param, " ")
		kv := strings.Split(afterTrimParam, "=")
		if len(kv) == 2 {
			rtnMap[kv[0]] = kv[1]
		} else {
			return rtnMap, fmt.Errorf("GetMapFromProviderParams meet illegal format param=%s", param)
		}
	}
	return rtnMap, nil
}

func callCiFunctionByCiType(req *CmdbRequest, ciType CiType) (*CmdbResponse, error) {
	resp := &CmdbResponse{}
	var err error

	if ciType == CI_TYPE_NORMAL {
		resp, _, err = OperateCi(req)
	} else if ciType == CI_TYPE_INTEGRATE {
		resp, _, err = getIntegrateTemplateData(req)
	} else {
		return nil, fmt.Errorf("callCiFunctionByCiType not valid citype=%v", ciType)
	}

	return resp, err
}

func isOrderParamValid(order string) error {
	validValues := []string{"", "asc", "desc"}

	for _, valid := range validValues {
		if valid == order {
			return nil
		}
	}

	logrus.Errorf("Invalid order(%s) param", order)
	return fmt.Errorf("invalid order:%s", order)
}

func listEntries(ciType CiType, ciName string, queryParam *CmdbCiQueryParam, results interface{}) (int, error) {
	if err := isValidPointer(results); err != nil {
		return 0, err
	}

	filter := make(map[string]interface{})
	for k, v := range queryParam.Filter {
		filter[k] = v
	}

	req := CmdbRequest{
		Type:          ciName,
		Action:        "select",
		Filter:        filter,
		StartIndex:    queryParam.Offset,
		PageSize:      queryParam.Limit,
		ResultColumn:  queryParam.ResultColumn,
		IsPaging:      true,
		PluginCode:    queryParam.PluginCode,
		PluginVersion: queryParam.PluginVersion,
	}
	logrus.Debugf("req: %++v", req)

	if err := isOrderParamValid(queryParam.Order); err != nil {
		return 0, err
	}

	if queryParam.Order != "" && queryParam.OrderBy != "" {
		req.Orderby = make(map[string]interface{})
		req.Orderby[queryParam.OrderBy] = queryParam.Order
	}

	resp, err := callCiFunctionByCiType(&req, ciType)
	if err != nil {
		logrus.Errorf("ListCiEntries meet error err=%v,req=%++v,ciType=%v", err, req, ciType)
		return 0, err
	}

	total, err := GetIntFromInterface(resp.Headers.TotalRows)
	if err != nil {
		logrus.Errorf("ListCiEntries:get totalRow meet error err=%v,totalRow=%v", err, resp.Headers.TotalRows)
		return 0, err
	}

	if total > 0 {
		err = UnmarshalContent(resp.Data.Content, results)
		if err != nil {
			logrus.Errorf("ListCiEntries UnmarshalContent meet error err=%v", err)
		}
	}

	return total, err
}

func DeleteCiEntryByGuid(guid, pluginCode, pluginVersion string, ciName string, bEnableCasCadingDelete bool) error {
	filter := make(map[string]interface{})
	filter["guid"] = guid

	req := CmdbRequest{
		Type:                  ciName,
		Action:                "delete",
		EnableCasCadingDelete: bEnableCasCadingDelete,
		LimitRowCount:         1,
		Filter:                filter,
		PluginCode:            pluginCode,
		PluginVersion:         pluginVersion,
	}

	_, _, err := OperateCi(&req)
	if err != nil {
		logrus.Errorf("DeleteCiEntryByGuid meet error err=%v,req=%++v", err, req)
	}

	return err
}

type CiInsertGuid struct {
	Guid string `json:"guid,omitempty"`
}

func GetNameAndGuidFromReferenceId(referenceId []map[string]string) (name string, guid string, err error) {
	if len(referenceId) != 1 {
		err = fmt.Errorf("GetNameAndGuidFromIdMapArray,len(id)=%d not equal 1", len(referenceId))
		return
	}
	name = referenceId[0]["v"]
	guid = referenceId[0]["k"]

	return
}

func UpdateCiEntryByGuid(ciName string, guid, pluginCode, pluginVersion string, ciEntries ...interface{}) error {
	return updateCiEntryByGuid(ciName, guid, pluginCode, pluginVersion, ciEntries...)
}

func updateCiEntryByGuid(ciName string, guid, pluginCode, pluginVersion string, ciEntries ...interface{}) error {
	parameters := []map[string]interface{}{}
	for _, ciEntries := range ciEntries {
		parameter, err := GetMapFromStruct(ciEntries)
		if err != nil {
			return err
		}
		parameters = append(parameters, parameter)
	}
	fileter := make(map[string]interface{})
	fileter["guid"] = guid
	req := CmdbRequest{
		Type:          ciName,
		Action:        "update",
		Filter:        fileter,
		Parameters:    parameters,
		PluginCode:    pluginCode,
		PluginVersion: pluginVersion,
	}
	resp, _, err := OperateCi(&req)
	if err != nil || resp.Headers.RetCode != 0 {
		logrus.Errorf("UpdateMultiCiEntries meet error err=%v", err)
		return err
	}
	return err
}

func GetOperateCi(request []byte) (response *CmdbResponse, origin []byte, err error) {
	params := CmdbRequest{}
	if err = json.Unmarshal(request, &params); err != nil {
		logrus.Errorf("unmarshal operaterci request failed, err=%s", err)
		return nil, nil, err
	}

	resp, bytes, err := OperateCi(&params)
	return resp, bytes, err
}

func GetIntegrateTemplateData(params *CmdbRequest) (response *CmdbResponse, origin []byte, err error) {
	resp, bytes, err := getIntegrateTemplateData(params)
	return resp, bytes, err
}

func ListIntegrateEntries(ciName string, queryParam *CmdbCiQueryParam, results interface{}) (int, error) {
	return listEntries(CI_TYPE_INTEGRATE, ciName, queryParam, results)
}
