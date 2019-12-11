package plugins

import (
	"bufio"
	"bytes"
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
	LogActions["searchdetail"] = new(LogSearchDetailAction)
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

//SearchInputs .
type SearchInputs struct {
	Inputs []SearchInput `json:"inputs,omitempty"`
}

//SearchInput .
type SearchInput struct {
	CallBackParameter
	Guid       string `json:"guid,omitempty"`
	KeyWord    string `json:"key_word,omitempty"`
	LineNumber int    `json:"line_number,omitempty"`
}

//SearchOutputs .
type SearchOutputs struct {
	Outputs []SearchOutput `json:"outputs,omitempty"`
}

//SearchOutput .
type SearchOutput struct {
	CallBackParameter
	Result
	FileName string `json:"file_name,omitempty"`
	Line     string `json:"line_number,omitempty"`
	Log      string `json:"log,omitempty"`
}

//ReadParam .
func (action *LogSearchAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SearchInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func logSearchCheckParam(log *SearchInput) error {
	if log.KeyWord == "" {
		return errors.New("LogSearchAction input KeyWord can not be empty")
	}

	return nil
}

//Do .
func (action *LogSearchAction) Do(input interface{}) (interface{}, error) {
	logs, _ := input.(SearchInputs)
	var logoutputs SearchOutputs
	var finalErr error

	for i := 0; i < len(logs.Inputs); i++ {
		output, err := action.Search(&logs.Inputs[i])
		if err != nil {
			finalErr = err
		}

		loginfo, _ := output.(SearchOutputs)
		for k := 0; k < len(loginfo.Outputs); k++ {
			logoutputs.Outputs = append(logoutputs.Outputs, loginfo.Outputs[k])
		}

	}

	return &logoutputs, finalErr
}

//Search .
func (action *LogSearchAction) Search(input *SearchInput) (interface{}, error) {
	sh := "cd logs && "

	keystring := []string{}
	if strings.Contains(input.KeyWord, ",") {
		keystring = strings.Split(input.KeyWord, ",")

		sh += "grep -rin '" + keystring[0] + "' *.log"

		for i := 1; i < len(keystring); i++ {
			sh += "|grep '" + keystring[i] + "'"
		}

	} else {
		sh += "grep -rin '" + input.KeyWord + "' *.log"
	}

	cmd := exec.Command("/bin/sh", "-c", sh)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("can not obtain stdout pipe for command when get log filename: %s \n", err)
		return nil, err
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("conmand start is error when get log filename: %s \n", err)
		return nil, err
	}

	output, err := LogReadLine(cmd, stdout)
	if err != nil {
		return nil, err
	}

	//get filename and lineinfo
	var infos SearchOutputs

	if len(output) > 0 {
		for k := 0; k < len(output); k++ {
			var info SearchOutput
			info.CallBackParameter.Parameter = input.CallBackParameter.Parameter
			info.Result.Code = RESULT_CODE_SUCCESS
			if output[k] == "" {
				continue
			}

			if !strings.Contains(output[k], ":time=") {
				continue
			}

			fileline := strings.Split(output[k], ":time=")

			if fileline[1] == "" {
				continue
			}

			//single log file
			if !strings.Contains(fileline[0], ":") {
				info.FileName = "wecube-plugins-qcloud.log"
				info.Line = fileline[0]
			} else {
				f := strings.Split(fileline[0], ":")
				info.FileName = f[0]
				info.Line = f[1]
			}

			if len(fileline) == 2 {
				info.Log = "time=" + fileline[1]
			}

			if len(fileline) > 2 {
				info.Log = "time="
				for j := 1; j < len(fileline); j++ {
					info.Log += fileline[j]
				}
			}

			infos.Outputs = append(infos.Outputs, info)
		}
	}

	return infos, nil
}

//LogSearchDetailAction .
type LogSearchDetailAction struct {
}

//SearchDetailInputs .
type SearchDetailInputs struct {
	Inputs []SearchDetailInput `json:"inputs,omitempty"`
}

//SearchDetailInput .
type SearchDetailInput struct {
	CallBackParameter
	FileName        string `json:"file_name,omitempty"`
	LineNumber      string `json:"line_number,omitempty"`
	RelateLineCount int    `json:"relate_line_count,omitempty"`
}

//SearchDetailOutputs .
type SearchDetailOutputs struct {
	Outputs []SearchDetailOutput `json:"outputs,omitempty"`
}

//SearchDetailOutput .
type SearchDetailOutput struct {
	CallBackParameter
	Result
	FileName   string `json:"file_name,omitempty"`
	LineNumber string `json:"line_number,omitempty"`
	Logs       string `json:"logs,omitempty"`
}

//ReadParam .
func (action *LogSearchDetailAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SearchDetailInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func logSearchDetailCheckParam(log *SearchDetailInput) error {
	if log.FileName == "" {
		return errors.New("LogSearchDetailAction input finename can not be empty")
	}
	if log.LineNumber == "" {
		return errors.New("LogSearchDetailAction input LineNumber can not be empty")
	}

	return nil
}

//Do .
func (action *LogSearchDetailAction) Do(input interface{}) (interface{}, error) {
	logs, _ := input.(SearchDetailInputs)
	var finalErr error
	var logoutputs SearchDetailOutputs

	for i := 0; i < len(logs.Inputs); i++ {
		var info SearchDetailOutput
		info.CallBackParameter.Parameter = logs.Inputs[i].CallBackParameter.Parameter
		info.Result.Code = RESULT_CODE_SUCCESS
		if err := logSearchDetailCheckParam(&logs.Inputs[i]); err != nil {
			info.Result.Code = RESULT_CODE_ERROR
			info.Result.Message = err.Error()
			finalErr = err
			logoutputs.Outputs = append(logoutputs.Outputs, info)
			continue
		}

		text, err := action.SearchDetail(&logs.Inputs[i])
		if err != nil {
			info.Result.Code = RESULT_CODE_ERROR
			info.Result.Message = err.Error()
			finalErr = err
			logoutputs.Outputs = append(logoutputs.Outputs, info)
			continue
		}

		info.FileName = logs.Inputs[i].FileName
		info.LineNumber = logs.Inputs[i].LineNumber
		info.Logs = text

		logoutputs.Outputs = append(logoutputs.Outputs, info)
	}

	return &logoutputs, finalErr
}

//SearchDetail .
func (action *LogSearchDetailAction) SearchDetail(input *SearchDetailInput) (string, error) {
	if input.RelateLineCount == 0 {
		input.RelateLineCount = 10
	}

	startLine, _ := strconv.Atoi(input.LineNumber)
	shellCmd := fmt.Sprintf("cd logs && cat -n %s |sed -n \"%d,%dp\" ", input.FileName, startLine, startLine+input.RelateLineCount)
	contextText, err := runCmd(shellCmd)
	if err != nil {
		return "", err
	}

	return contextText, nil
}

//LogReadLine .
func LogReadLine(cmd *exec.Cmd, stdout io.ReadCloser) ([]string, error) {

	linelist := []string{}
	outputBuf := bufio.NewReader(stdout)

	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			if err.Error() != "EOF" {
				logrus.Info("readline is error")
				return []string{}, nil
			}
		}

		linelist = append(linelist, string(output))
	}

	if err := cmd.Wait(); err != nil {
		return []string{}, nil
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

func runCmd(shellCommand string) (string, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", shellCommand)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		logrus.Errorf("runCmd (%s) meet err=%v,stderr=%v", shellCommand, err, stderr.String())
		return stderr.String(), nil
	}

	return stdout.String(), nil
}
