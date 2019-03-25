package conf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type AppConfig struct {
	HttpPort        string
	CMDBLink        string
	CMDBUserAuthKey string
}

type AppConfigMgr struct {
	Config atomic.Value
}

type Config struct {
	Filename       string
	Items          map[string]string
	LastUpdateTime int64
	RWLock         sync.RWMutex
	NotifyerList   []Notifyer
}

var AppConfMgr = &AppConfigMgr{}
var GobalAppConfig = &AppConfig{}

func InitConfig(file string) {
	conf, err := NewConfig(file)
	if err != nil {
		fmt.Printf("read config file err: %v\n", err)
		return
	}

	GobalAppConfig.HttpPort, err = conf.GetString("httpport")
	if err != nil {
		fmt.Printf("get HttpPort err: %v\n", err)
		return
	}
	GobalAppConfig.CMDBLink, err = conf.GetString("cmdb.link")
	if err != nil {
		fmt.Printf("get CMDBLink err: %v\n", err)
		return
	}
	GobalAppConfig.CMDBUserAuthKey, err = conf.GetString("cmdb.user.auth.key")
	if err != nil {
		fmt.Printf("get CMDBUserAuthKey err: %v\n", err)
		return
	}

	AppConfMgr.Config.Store(GobalAppConfig)
}

func NewConfig(file string) (conf *Config, err error) {
	conf = &Config{
		Filename: file,
		Items:    make(map[string]string, 1024),
	}

	m, err := conf.parse()
	if err != nil {
		fmt.Printf("parse conf error:%v\n", err)
		return
	}

	conf.RWLock.Lock()
	conf.Items = m
	conf.RWLock.Unlock()

	return
}

func (c *Config) parse() (m map[string]string, err error) {
	m = make(map[string]string, 1024)
	f, err := os.Open(c.Filename)
	if err != nil {
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var lineNo int
	for {
		var line string
		byteLine, _, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		line = string(byteLine)
		if er == io.EOF {
			lineParse(&lineNo, &line, &m)
			break
		}
		lineParse(&lineNo, &line, &m)
	}
	return
}

func lineParse(lineNo *int, line *string, m *map[string]string) {
	*lineNo++

	l := strings.TrimSpace(*line)
	if len(l) == 0 || l[0] == '\n' || l[0] == '#' || l[0] == ';' {
		return
	}

	itemSlice := strings.Split(l, "=")
	if len(itemSlice) == 0 {
		fmt.Printf("invalid config, line:%d", lineNo)
		return
	}

	key := strings.TrimSpace(itemSlice[0])
	if len(key) == 0 {
		fmt.Printf("invalid config, line:%d", lineNo)
		return
	}
	if len(key) == 1 {
		(*m)[key] = ""
		return
	}

	value := strings.TrimSpace(itemSlice[1])
	(*m)[key] = value

	return
}

func (c *Config) GetInt(key string) (value int, err error) {
	c.RWLock.RLock()
	defer c.RWLock.RUnlock()

	str, ok := c.Items[key]
	if !ok {
		err = fmt.Errorf("key [%s] not found", key)
	}
	value, err = strconv.Atoi(str)
	return
}

func (c *Config) GetIntDefault(key string, defaultInt int) (value int) {
	c.RWLock.RLock()
	defer c.RWLock.RUnlock()

	str, ok := c.Items[key]
	if !ok {
		value = defaultInt
		return
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		value = defaultInt
	}
	return
}

func (c *Config) GetString(key string) (value string, err error) {
	c.RWLock.RLock()
	defer c.RWLock.RUnlock()

	value, ok := c.Items[key]
	if !ok {
		err = fmt.Errorf("key [%s] not found", key)
	}
	return
}

func (c *Config) GetIStringDefault(key string, defaultStr string) (value string) {
	c.RWLock.RLock()
	defer c.RWLock.RUnlock()

	value, ok := c.Items[key]
	if !ok {
		value = defaultStr
		return
	}
	return
}

type Notifyer interface {
	Callback(*Config)
}
