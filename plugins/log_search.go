package plugins

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

//LogActions define
var LogActions = make(map[string]Action)

func init() {
	LogActions["search"] = new(LogSearchAction)
}

//LogInputs .
type LogInputs struct {
	Inputs []LogInput `json:"inputs,omitempty"`
}

//LogInput .
type LogInput struct {
	Guid       string `json:"guid,omitempty"`
	KeyWord    string `json:"key_word,omitempty"`
	LineNumber int    `json:"line_number,omitempty"`
}

//LogOutputs .
type LogOutputs struct {
	Outputs []LogOutput `json:"outputs,omitempty"`
}

//LogOutput .
type LogOutput struct {
	Guid string     `json:"guid,omitempty"`
	Logs [][]string `json:"logs,omitempty"`
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

//LogSearchAction .
type LogSearchAction struct {
}

//ReadParam .
func (action *LogSearchAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs LogInput
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *LogSearchAction) CheckParam(input interface{}) error {
	log, ok := input.(LogInput)
	if !ok {
		return fmt.Errorf("LogSearchAAction:input type=%T not right", input)
	}

	if log.KeyWord == "" {
		return errors.New("LogSearchAAction input KeyWord can not be empty")
	}

	return nil
}

//Do .
func (action *LogSearchAction) Do(input interface{}) (interface{}, error) {
	logs, _ := input.(LogInputs)
	var logoutputs LogOutputs

	for k := 0; k < len(logs.Inputs); k++ {
		output, err := action.SearchLineNumber(&logs.Inputs[k])
		if err != nil {
			return nil, err
		}

		if len(output) > 0 {
			for i := 0; i < len(output); i++ {
				if output[i] == "" {
					continue
				}

				lineinfo, err := action.Search(logs.Inputs[k].LineNumber, output[i])
				if err != nil {
					return nil, err
				}

				var out LogOutput
				out.Guid = logs.Inputs[k].Guid
				out.Logs = append(out.Logs, lineinfo)
				logoutputs.Outputs = append(logoutputs.Outputs, out)
			}
		}
		logrus.Infof("all keyword relate information = %v are getted", logs.Inputs[k].KeyWord)
	}

	return &logoutputs, nil
}

//Search .
func (action *LogSearchAction) Search(searchLine int, LineNumber string) ([]string, error) {
	if searchLine == 0 {
		searchLine = 10
	}

	sh := "cat -n logs/wecube-plugins.log |tail -n +"
	startLine, needLine := CountLineNumber(searchLine, LineNumber)
	sh += startLine + " | head -n " + needLine

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

	output, err := LogReadLine(cmd, stdout)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//SearchLineNumber .
func (action *LogSearchAction) SearchLineNumber(input *LogInput) ([]string, error) {

	keystring := []string{}
	if strings.Contains(input.KeyWord, ",") {
		keystring = strings.Split(input.KeyWord, ",")
	}

	sh := "cat -n logs/wecube-plugins.log "
	if len(keystring) > 1 {
		for _, key := range keystring {
			sh += "|grep " + key
		}
	} else {
		sh += "|grep " + input.KeyWord
	}
	sh += " |awk '{print $1}';echo $1 "
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

	output, err := LogReadLine(cmd, stdout)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//LogReadLine .
func LogReadLine(cmd *exec.Cmd, stdout io.ReadCloser) ([]string, error) {

	var linelist []string
	outputBuf := bufio.NewReader(stdout)

	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			if err.Error() != "EOF" {
				logrus.Info("readline is error")
				return []string{}, err
			}
		}

		linelist = append(linelist, string(output))
	}

	if err := cmd.Wait(); err != nil {
		return []string{}, err
	}

	return linelist, nil
}

//CountLineNumber .
func CountLineNumber(wLine int, rLine string) (string, string) {

	rline, _ := strconv.Atoi(rLine)

	var num int

	var startLineNumber int
	if rline <= wLine {
		startLineNumber = 1
		num = wLine + rline
	} else {
		startLineNumber = rline - wLine
		num = 2*wLine + 1
	}

	line1 := strconv.Itoa(startLineNumber)

	line2 := strconv.Itoa(num)

	return line1, line2
}
