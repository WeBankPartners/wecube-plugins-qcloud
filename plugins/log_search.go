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
	logOutput, err := action.GetKeyWordLineNumber(&log)
	if err != nil {
		return nil, err
	}

	logrus.Info("linenumber is ====================>", len(logOutput))

	// logOutput, err := action.GetKeyWord(&log, output)
	// if err != nil {
	// 	return nil, err
	// }

	logrus.Infof("all keyword relate information = %v are getted", log.KeyWord)
	return &logOutput, nil
}

//GetKeyWord .
func (action *LogGetKeyWordAction) GetKeyWord(input *LogInput, LineNumber []string) (interface{}, error) {
	if input.LineNumber == "" {
		input.LineNumber = "10"
	}

	var outputs []LogOutputs

	for i := 0; i < len(LineNumber); i++ {
		sh := "cat -n wecube-plugins.log |tail -n +1 | head -n 10"

		// sh := "cat -n wecube-plugins.log |tail -n +"

		// startLine, needLine := CountLineNumber(input.LineNumber, LineNumber[i])

		// sh += startLine + " | head -n " + needLine

		// logrus.Info("commandsss =================> ", sh)

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

		//按行读取
		output, err := LogReadLine(stdout)
		if err != nil {
			return nil, err
		}

		if len(output) > 0 {
			var out LogOutputs
			out.Outputs = output
			outputs = append(outputs, out)
		}
	}

	return outputs, nil
}

//GetKeyWordLineNumber .
func (action *LogGetKeyWordAction) GetKeyWordLineNumber(input *LogInput) ([]string, error) {

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

	//按行读取
	output, err := LogReadLine(stdout)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//LogReadLine .
func LogReadLine(stdout io.ReadCloser) ([]string, error) {

	var linelist []string
	outputBuf := bufio.NewReader(stdout)

	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			logrus.Info("readline is error")
			return []string{}, err
		}
		linelist = append(linelist, string(output))
	}

	return linelist, nil
}

//CountLineNumber .
func CountLineNumber(wLine string, rLine string) (string, string) {

	wline, _ := strconv.Atoi(wLine)
	rline, _ := strconv.Atoi(rLine)

	num := 2 * wline

	var startLineNumber int
	if rline <= wline {
		startLineNumber = 1
	} else {
		startLineNumber = rline - wline
	}

	line1 := strconv.Itoa(startLineNumber)

	line2 := strconv.Itoa(num)

	return line1, line2
}
