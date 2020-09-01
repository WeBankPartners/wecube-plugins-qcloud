package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"fmt"

	_ "github.com/WeBankPartners/wecube-plugins-qcloud/plugins/bussiness_plugins/security_group"

	"github.com/WeBankPartners/wecube-plugins-qcloud/conf"
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"io/ioutil"
	"encoding/base64"
	"bytes"
	"github.com/dgrijalva/jwt-go"
	"strconv"
)

const (
	CONF_FILE_PATH = "./conf/app.conf"
)

var (
	CoreJwtKey string
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
	CoreJwtKey = os.Getenv("JWT_SIGNING_KEY")
	conf.InitConfig(CONF_FILE_PATH)
}

func initRouter() {
	//path should be defined as "/[package name]/[version]/[plugin]/[action]"
	http.HandleFunc("/", routeDispatcher)
}

func routeDispatcher(w http.ResponseWriter, r *http.Request) {
	if authCore(r.Header.Get("Authorization")) {
		pluginRequest := parsePluginRequest(r)
		pluginResponse, _ := plugins.Process(pluginRequest)
		b, _ := json.Marshal(pluginResponse)
		logrus.Infof("write data to client response=%s", string(b))
		write(w, pluginResponse)
	}else{
		logrus.Warnf("Request token illegal ----------------!!")
		pluginResponse := plugins.PluginResponse{ResultCode:"1",ResultMsg:"Token illegal"}
		write(w, &pluginResponse)
	}
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
		pluginInput.Version = pathStrings[2]
		pluginInput.ProviderName = pathStrings[1]
		pluginInput.Name = pathStrings[len(pathStrings)-2]
		pluginInput.Action = pathStrings[len(pathStrings)-1]
	}
	pluginInput.Parameters = r.Body
	logrus.Infof("parsed request = %v", pluginInput)
	return &pluginInput
}

func authCore(coreToken string) bool {
	_,err := decodeCoreToken(coreToken, CoreJwtKey)
	if err == nil {
		return true
	}
	return false
}

type coreJwtToken struct {
	User    string    `json:"user"`
	Expire  int64     `json:"expire"`
	Roles   []string  `json:"roles"`
}

func decodeCoreToken(token,key string) (result coreJwtToken,err error) {
	if strings.HasPrefix(token, "Bearer") {
		token = token[7:]
	}
	if key == "" || strings.HasPrefix(key, "{{") {
		key = "Platform+Auth+Server+Secret"
	}
	keyBytes,err := ioutil.ReadAll(base64.NewDecoder(base64.RawStdEncoding, bytes.NewBufferString(key)))
	if err != nil {
		logrus.Error("Decode core token fail,base64 decode error", err)
		return result,err
	}
	pToken,err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return keyBytes, nil
	})
	if err != nil {
		logrus.Error("Decode core token fail,jwt parse error", err)
		return result,err
	}
	claimMap,ok := pToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.Error("Decode core token fail,claims to map error", err)
		return result,err
	}
	result.User = fmt.Sprintf("%s", claimMap["sub"])
	result.Expire,err = strconv.ParseInt(fmt.Sprintf("%.0f", claimMap["exp"]), 10, 64)
	if err != nil {
		logrus.Error("Decode core token fail,parse expire to int64 error", err)
		return result,err
	}
	roleListString := fmt.Sprintf("%s", claimMap["authority"])
	roleListString = roleListString[1:len(roleListString)-1]
	result.Roles = strings.Split(roleListString, ",")
	return result,nil
}