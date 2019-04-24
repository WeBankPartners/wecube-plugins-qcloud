package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"git.webank.io/wecube-plugins/conf"
	"git.webank.io/wecube-plugins/plugins"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

const (
	CONF_FILE_PATH = "./conf/app.conf"
)

func init() {
	initConfig()
	initLogger()
	initRouter()
}

func main() {
	logrus.Infof("Start WeCube-Plungins Service ... ")

	go LogTest()

	if err := http.ListenAndServe(":"+conf.GobalAppConfig.HttpPort, nil); err != nil {
		logrus.Fatalf("ListenAndServe meet err = %v", err)
	}
}

func initLogger() {
	fileName := "logs/wecube-plugins.log"
	logrus.SetReportCaller(true)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logrus.SetOutput(file)
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:  fileName,
		MaxSize:   5,
		MaxAge:    7,
		Level:     logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{DisableTimestamp: false, DisableColors: false},
	})
	logrus.AddHook(rotateFileHook)
}

func initConfig() {
	conf.InitConfig(CONF_FILE_PATH)
}

func initRouter() {
	//path should be defined as "/[version]/[provider]/[plugin]/[action]"
	http.HandleFunc("/v1/qcloud/vm/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/vm/start", routeDispatcher)
	http.HandleFunc("/v1/qcloud/vm/stop", routeDispatcher)
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
	http.HandleFunc("/v1/qcloud/route-table/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/route-table/terminate", routeDispatcher)
	http.HandleFunc("/v1/qcloud/mysql-vm/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/mysql-vm/terminate", routeDispatcher)
	http.HandleFunc("/v1/qcloud/mysql-vm/restart", routeDispatcher)
	http.HandleFunc("/v1/qcloud/redis/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/log/search", routeDispatcher)
}

func routeDispatcher(w http.ResponseWriter, r *http.Request) {
	pluginRequest := parsePluginRequest(r)
	pluginResponse, _ := plugins.Process(pluginRequest)
	write(w, pluginResponse)
}

func write(w http.ResponseWriter, output *plugins.PluginResponse) {
	w.Header().Set("content-type", "application/json")
	b, err := json.Marshal(output)
	if err != nil {
		logrus.Errorf("write http response (%v) meet error (%v)", output, err)
	}
	w.Write(b)
}

func parsePluginRequest(r *http.Request) *plugins.PluginRequest {
	var pluginInput = plugins.PluginRequest{}
	pathStrings := strings.Split(r.URL.Path, "/")
	logrus.Infof("path strings = %v", pathStrings)
	if len(pathStrings) >= 5 {
		pluginInput.Version = pathStrings[1]
		pluginInput.ProviderName = pathStrings[2]
		pluginInput.Name = pathStrings[3]
		pluginInput.Action = pathStrings[4]
	}
	pluginInput.Parameters = r.Body
	logrus.Infof("parsed request = %v", pluginInput)
	return &pluginInput
}

func LogTest() {
	for {
		logrus.Info("this is a test for log file, through this function we can see the new log finename is what")
	}
}
