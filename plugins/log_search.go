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

//LogFileNameLineInfo .
type LogFileNameLineInfo struct {
	FileName string   `json:"name,omitempty"`
	Line     []string `json:"line,omitempty"`
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
	var inputs LogInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *LogSearchAction) CheckParam(input interface{}) error {
	logs, ok := input.(LogInputs)
	if !ok {
		return fmt.Errorf("LogSearchAAction:input type=%T not right", input)
	}

	for _, log := range logs.Inputs {
		if log.KeyWord == "" {
			return errors.New("LogSearchAAction input KeyWord can not be empty")
		}
	}

	return nil
}

//Do .
func (action *LogSearchAction) Do(input interface{}) (interface{}, error) {
	logs, _ := input.(LogInputs)
	var logoutputs LogOutputs

	for k := 0; k < len(logs.Inputs); k++ {
		//获取到文件名和行号的信息
		output, err := action.GetLogFileNameAndLineNumberByKeyword(&logs.Inputs[k])
		if err != nil {
			return nil, err
		}

		var out LogOutput
		out.Guid = logs.Inputs[k].Guid

		if len(output) == 0 {
			continue
		}

		for i := 0; i < len(output); i++ {
			if output[i].FileName == "" {
				continue
			}

			if len(output[i].Line) == 0 {
				continue
			}

			for j := 0; j < len(output[i].Line); j++ {
				lineinfo, err := action.Search(output[i].FileName, logs.Inputs[k].LineNumber, output[i].Line[j])
				if err != nil {
					return nil, err
				}

				out.Logs = append(out.Logs, lineinfo)
			}
		}

		if len(out.Logs) > 0 {
			logoutputs.Outputs = append(logoutputs.Outputs, out)
		}

		logrus.Infof("all keyword relate information = %v are getted", logs.Inputs[k].KeyWord)
	}

	return &logoutputs, nil
}

//Search .
func (action *LogSearchAction) Search(filename string, searchLine int, LineNumber string) ([]string, error) {
	if searchLine == 0 {
		searchLine = 10
	}

	// sh := "cat -n wecube-plugins.log |tail -n +"
	sh := "cd logs && cat -n " + filename + " |tail -n +"
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

		str := string(output)
		str1 := strings.Replace(str, "\t", "  ", -1)

		linelist = append(linelist, str1)
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

//GetLogFileNameAndLineNumberByKeyword .
func (action *LogSearchAction) GetLogFileNameAndLineNumberByKeyword(input *LogInput) (info []LogFileNameLineInfo, err error) {

	sh := "cd logs && "

	keystring := []string{}
	if strings.Contains(input.KeyWord, ",") {
		logrus.Info("come here  111====>>>>>>>")
		keystring = strings.Split(input.KeyWord, ",")
		logrus.Info("come here  22222====>>>>>>>", keystring)

		sh += "grep -rin '" + keystring[0] + "' *.log"
		logrus.Info("come here  3333====>>>>>>>", sh)

		for i := 1; i <= len(keystring); i++ {
			sh += "|grep '" + keystring[i] + "'"
		}

	} else {
		sh += "grep -rin '" + input.KeyWord + "' *.log"
	}

	logrus.Info("come here  4444====>>>>>>>", sh)

	sh += " |awk '{print $1}';echo $1 "
	cmd := exec.Command("/bin/sh", "-c", sh)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("can not obtain stdout pipe for command when get log filename: %s \n", err)
		return []LogFileNameLineInfo{}, err
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("conmand start is error when get log filename: %s \n", err)
		return []LogFileNameLineInfo{}, err
	}

	output, err := LogReadLine(cmd, stdout)
	if err != nil {
		return nil, err
	}

	//获取输出中的文件名和行号
	var infos []LogFileNameLineInfo

	lineinfos := make(map[string][]string)

	if len(output) > 0 {
		for k := 0; k < len(output); k++ {
			// var info LogFileNameLineInfo

			if output[k] == "" {
				continue
			}
			if !strings.Contains(output[k], ":") {
				continue
			}

			fileline := strings.Split(output[k], ":")

			//单个日志文件的情况，不会输出文件名
			if !strings.Contains(output[k], "log") {

				// fileline := strings.Split(output[k], ":")
				lineinfos["wecube-plugins.log"] = append(lineinfos["wecube-plugins.log"], fileline[0])
				// info.FileName = "wecube-plugins.log"
				// info.Line = appen(info.Line, fileline[0])
				// infos = append(infos, info)
			} else {
				//多个日志文件的情况，会输出文件名
				// fileline := strings.Split(output[k], ":")
				lineinfos[fileline[0]] = append(lineinfos[fileline[0]], fileline[1])
			}

			// info.FileName = fileline[0]
			// info.Line = append(info.Line, fileline[1])

			// infos = append(infos, info)
		}
	}

	for filename, message := range lineinfos {
		var info LogFileNameLineInfo
		info.FileName = filename
		info.Line = message

		infos = append(infos, info)
	}

	return infos, nil
}
