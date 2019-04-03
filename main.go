package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"git.webank.io/wecube-plugins/conf"
	"git.webank.io/wecube-plugins/plugins"
	"github.com/sirupsen/logrus"
)

const (
	CONF_FILE_PATH = "./conf/app.conf"
)

func main() {
	logrus.Infof("Start WeCube-Plungins Service ... ")
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)

	initConfig()

	initRouter()

	if err := http.ListenAndServe(":"+conf.GobalAppConfig.HttpPort, nil); err != nil {
		logrus.Fatalf("ListenAndServe meet err = %v", err)
	}
}

func initConfig() {
	conf.InitConfig(CONF_FILE_PATH)
}

func initRouter() {
	//path should be define as "/[version]/[provider]/[plugin]/[action]"
	http.HandleFunc("/v1/qcloud/vm/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/vm/terminate", routeDispatcher)
	http.HandleFunc("/v1/qcloud/storage/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/storage/terminate", routeDispatcher)
	http.HandleFunc("/v1/qcloud/subnet/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/subnet/terminate", routeDispatcher)
	http.HandleFunc("/v1/qcloud/security-group/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/security-group/terminate", routeDispatcher)
	http.HandleFunc("/v1/qcloud/nat-gateway/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/nat-gateway/terminate", routeDispatcher)
	http.HandleFunc("/v1/qcloud/vpc/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/vpc/terminate", routeDispatcher)
	http.HandleFunc("/v1/qcloud/peering-connection/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/peering-connection/terminate", routeDispatcher)
}

func routeDispatcher(w http.ResponseWriter, r *http.Request) {
	var err error
	var workflowParam plugins.WorkflowParam

	defer func() {
		if err != nil {
			logrus.Error(err)
			OutputJson(w, &workflowParam, fmt.Sprint(err))
		}
	}()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bodyBytes, &workflowParam); err != nil {
		return
	}
	pathStrings := strings.Split(r.URL.Path, "/")
	if len(pathStrings) >= 5 {
		workflowParam.PluginVersion = pathStrings[1]
		workflowParam.ProviderName = pathStrings[2]
		workflowParam.PluginName = pathStrings[3]
		workflowParam.ApplicationAction = pathStrings[4]
	}

	go plugins.CallPluginAction(workflowParam)

	OutputJson(w, &workflowParam, "")
}

const (
	RESULT_CODE_SUCCESSFUL = "0"
	RESULT_CODE_ERROR      = "1"
)

func OutputJson(w http.ResponseWriter, res *plugins.WorkflowParam, resultMsg string) {
	w.Header().Set("content-type", "application/json")

	if resultMsg == "" {
		res.ResultCode = RESULT_CODE_SUCCESSFUL
	} else {
		res.ResultCode = RESULT_CODE_ERROR
		res.ResultMsg = resultMsg
	}
	b, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Write(b)
}
