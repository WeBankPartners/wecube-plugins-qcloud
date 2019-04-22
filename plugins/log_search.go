package plugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/sirupsen/logrus"
)

//LogActions define
var LogActions = make(map[string]Action)

func init() {
	LogActions["getkeyword"] = new(LogGetKeyWordAction)
}

//LogInput .
type LogInput struct {
	KeyWord    string `json:"key_word,omitempty"`
	LineNumber string `json:"line_number,omitempty"`
}

//LogOutputs .
type LogOutputs struct {
	Outputs []string `json:"outputs,omitempty"`
}

//LogPlugin .
type LogPlugin struct {
}

//GetActionByName .
func (plugin *LogPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := LogActions[actionName]

	if !found {
		return nil, fmt.Errorf("Log plugin,action = %s not found", actionName)
	}

	return action, nil
}

//LogGetKeyWordAction .
type LogGetKeyWordAction struct {
}

//ReadParam .
func (action *LogGetKeyWordAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs LogInput
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *LogGetKeyWordAction) CheckParam(input interface{}) error {
	log, ok := input.(LogInput)
	if !ok {
		return fmt.Errorf("LogGetKeyWordAction:input type=%T not right", input)
	}

	if log.KeyWord == "" {
		return errors.New("LogGetKeyWordAction input KeyWord can not be empty")
	}

	return nil
}

//Do .
func (action *LogGetKeyWordAction) Do(input interface{}) (interface{}, error) {
	log, _ := input.(LogInput)
	logOutput, err := action.GetKeyWord(&log)
	if err != nil {
		return nil, err
	}

	logrus.Infof("all keyword relate information = %v are getted", log.KeyWord)
	return &logOutput, nil
}

//GetKeyWord .
func (action *LogGetKeyWordAction) GetKeyWord(input *LogInput) (interface{}, error) {
	if input.LineNumber == "" {
		input.LineNumber = "10"
	}

	sh := "cat logs/wecube-plugins.log |grep " + input.KeyWord + " -C " + input.LineNumber
	fmt.Println("command shell==============>", sh)
	cmd := exec.Command("/bin/sh", "-c", sh)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("can not obtain stdout pipe for command: %s \n", err)
		return []string{}, err
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("conmand start is error: %s \n", err)
		return []string{}, err
	}

	//读取输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Printf("ReadAll stdout error: %s \n", err)
		return []string{}, err
	}

	var outputs LogOutputs
	outputs.Outputs = append(outputs.Outputs, string(bytes))

	return outputs, nil
}
