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
	LogActions["searchlog"] = new(LogSearchLogAction)
	// LogActions["searchdetail"] = new(LogSearchDetailAction)
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
		return fmt.Errorf("LogSearchAction:input type=%T not right", input)
	}

	for _, log := range logs.Inputs {
		if log.KeyWord == "" {
			return errors.New("LogSearchAction input KeyWord can not be empty")
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
	}

	logrus.Infof("all keyword relate information = %v are getted", logs.Inputs)

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

		// str := string(output)
		// str1 := strings.Replace(str, "\t", "  ", -1)

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

//GetLogFileNameAndLineNumberByKeyword .
func (action *LogSearchAction) GetLogFileNameAndLineNumberByKeyword(input *LogInput) (info []LogFileNameLineInfo, err error) {

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

			if output[k] == "" {
				continue
			}
			if !strings.Contains(output[k], ":") {
				continue
			}

			fileline := strings.Split(output[k], ":")

			//单个日志文件的情况，不会输出文件名
			if !strings.Contains(output[k], "log") {

				lineinfos["wecube-plugins.log"] = append(lineinfos["wecube-plugins.log"], fileline[0])

			} else {
				//多个日志文件的情况，会输出文件名
				lineinfos[fileline[0]] = append(lineinfos[fileline[0]], fileline[1])
			}
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

//LogSearchLogAction .
type LogSearchLogAction struct {
}

//SearchLogInputs .
type SearchLogInputs struct {
	Inputs []SearchLogInput `json:"inputs,omitempty"`
}

//SearchLogInput .
type SearchLogInput struct {
	Guid       string `json:"guid,omitempty"`
	KeyWord    string `json:"key_word,omitempty"`
	LineNumber int    `json:"line_number,omitempty"`
}

//SearchLogOutputs .
type SearchLogOutputs struct {
	Outputs []SearchLogOutput `json:"outputs,omitempty"`
}

//SearchLogOutput .
type SearchLogOutput struct {
	FileName string `json:"name,omitempty"`
	Line     string `json:"line,omitempty"`
	Log      string `json:"log,omitempty"`
}

//ReadParam .
func (action *LogSearchLogAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SearchLogInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *LogSearchLogAction) CheckParam(input interface{}) error {
	logs, ok := input.(SearchLogInputs)
	if !ok {
		return fmt.Errorf("LogSearchAction:input type=%T not right", input)
	}

	for _, log := range logs.Inputs {
		if log.KeyWord == "" {
			return errors.New("LogSearchAction input KeyWord can not be empty")
		}
	}

	return nil
}

//Do .
func (action *LogSearchLogAction) Do(input interface{}) (interface{}, error) {
	logs, _ := input.(SearchLogInputs)

	var logoutputs SearchLogOutputs

	for i := 0; i < len(logs.Inputs); i++ {
		output, err := action.SearchLog(&logs.Inputs[i])
		if err != nil {
			return nil, err
		}

		loginfo, _ := output.(SearchLogOutputs)

		for k := 0; k < len(loginfo.Outputs); k++ {
			logoutputs.Outputs = append(logoutputs.Outputs, loginfo.Outputs[k])
		}

	}

	return &logoutputs, nil
}

//SearchLog .
func (action *LogSearchLogAction) SearchLog(input *SearchLogInput) (interface{}, error) {

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

	sh += " |awk '{print $1}';echo $1 "
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

	//获取输出中的文件名和行号
	var infos SearchLogOutputs

	if len(output) > 0 {
		for k := 0; k < len(output); k++ {
			var info SearchLogOutput

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

			//单个日志文件的情况，不会输出文件名
			if !strings.Contains(fileline[0], ":") {
				info.FileName = "wecube-plugins.log"
				info.Line = fileline[0]
			} else {
				f := strings.Split(fileline[0], ":")
				info.FileName = f[0]
				info.Line = f[1]
			}

			logrus.Info("fileline info ==========>>>>>", fileline[1])

			if len(fileline) == 2 {
				info.Log = "time=" + fileline[1]
				logrus.Info("fileline = 2 =====here=====>>>>>")
			}
			if len(fileline) > 2 {
				info.Log = "time="
				for j := 1; j < len(fileline); j++ {
					info.Log += fileline[j]
				}
				logrus.Info("fileline > 2 =====here=====>>>>>")
			}

			infos.Outputs = append(infos.Outputs, info)
		}
	}

	return infos, nil
}
