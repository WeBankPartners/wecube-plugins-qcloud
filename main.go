package main

import (
	"encoding/json"
	"net/http"

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
	http.HandleFunc("/v1/qcloud/route-table/create", routeDispatcher)
	http.HandleFunc("/v1/qcloud/route-table/terminate", routeDispatcher)
}

func routeDispatcher(w http.ResponseWriter, r *http.Request) {
	pluginResponse, _ := plugins.CallPluginAction(r)
	write(w, pluginResponse)
}

func write(w http.ResponseWriter, output *plugins.PluginResponse) {
	w.Header().Set("content-type", "application/json")
	b, err := json.Marshal(output)
	if err != nil {
		logrus.Error("write http response (%v) meet error (%v)", output, err)
	}
	w.Write(b)
}
