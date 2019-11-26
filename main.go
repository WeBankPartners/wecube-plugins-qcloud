package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	_ "github.com/WeBankPartners/wecube-plugins-qcloud/plugins/bussiness_plugins/security_group"

	"github.com/WeBankPartners/wecube-plugins-qcloud/conf"
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
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
	logrus.Infof("Start WeCube-Plungins-Qcloud Service ... ")

	if err := http.ListenAndServe(":"+conf.GobalAppConfig.HttpPort, nil); err != nil {
		logrus.Fatalf("ListenAndServe meet err = %v", err)
	}
}

func initLogger() {
	fileName := "logs/wecube-plugins-qcloud.log"
	logrus.SetReportCaller(true)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logrus.SetOutput(file)
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     7,
		Level:      logrus.InfoLevel,
		Formatter:  &logrus.TextFormatter{DisableTimestamp: false, DisableColors: false},
	})
	logrus.AddHook(rotateFileHook)
}

func initConfig() {
	conf.InitConfig(CONF_FILE_PATH)
}

func initRouter() {
	//path should be defined as "/[package name]/[version]/[plugin]/[action]"
	http.HandleFunc("/qcloud/v1/", routeDispatcher)
}

func routeDispatcher(w http.ResponseWriter, r *http.Request) {
	pluginRequest := parsePluginRequest(r)
	pluginResponse, _ := plugins.Process(pluginRequest)
	logrus.Infof("write data to client response=%++v", pluginResponse)
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
