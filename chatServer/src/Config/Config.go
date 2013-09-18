package Config

import (
	"Model"
	"Utils"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/**项目总体配置结构体**/
type VsConfig struct {
	XMLName         xml.Name `xml:"Config"`
	ServerIp        string
	LogicServerList string
}

/**项目总体配置结构体对象**/
var vsConfig VsConfig

/**配置文件路径**/
var ConfigFile string

/**日志文件路径**/
var logFile string

/**错误日志文件路径**/
var logErrFile string

/**逻辑服务器列表**/
var logicServerList []*Model.ServerLevel

/**逻辑服务器列表map**/
var logicServerMapList []map[string]string

/**
flag.StringVar的作用是解析外域传来的参数，这俩句可以解析./main -c *** -l XXX 后面的参数***和XXX，即ConfigFile和LogFile的地址
若想用默认的，直接./main即可
**/
func init() {
	flag.StringVar(&logFile, "l", "", "log file path and name")
	flag.StringVar(&logErrFile, "e", "", "logErr file path and name")
	flag.StringVar(&ConfigFile, "c", "", "config file path and name")
	if ConfigFile == "" {
		ConfigFile = "./serverconf.xml"
	}
	if logFile == "" {
		logFile = fmt.Sprintf("./%s.log", time.Now().String())
	}
	if logErrFile == "" {
		logErrFile = "./logErr.log"
	}
	ParseXml(ConfigFile)
}

/**
解析xml文件
**/
func ParseXml(configFile string) {
	file, err := os.Open(configFile)
	if err != nil {
		Utils.LogPanicErr(err)
		return
	}
	xmlObj := xml.NewDecoder(file)
	err = xmlObj.Decode(&vsConfig)
	if err != nil {
		Utils.LogPanicErr(err)
		return
	}
	LogicServerListArray := strings.Split(vsConfig.LogicServerList, ";")
	count := len(LogicServerListArray)
	logicServerList = make([]*Model.ServerLevel, count)
	logicServerMapList = make([]map[string]string, count)
	for i := 0; i < count; i++ {
		serverLevel := new(Model.ServerLevel)
		perServerArray := strings.Split(LogicServerListArray[i], "|")
		serverLevel.ServerId, err = strconv.Atoi(perServerArray[0])
		if err != nil {
			Utils.LogErr(err)
			return
		}
		serverLevel.ServerAddress = perServerArray[1]
		serverLevel.ServerName = perServerArray[2]
		serverLevel.CurStatus, err = strconv.Atoi(perServerArray[3])
		if err != nil {
			Utils.LogErr(err)
			return
		}
		serverLevel.MaxClient, err = strconv.Atoi(perServerArray[4])
		if err != nil {
			Utils.LogErr(err)
			return
		}
		logicServerList[i] = serverLevel
		logicServerMapList[i] = Model.ChangeServerLevelToMap(serverLevel)
	}
	Utils.LogInfo("parse config xml=%v\n", vsConfig)
}

/**
得到服务器本机Ip的配置
**/
func GetServerIp() string {
	return vsConfig.ServerIp
}

/**
得到逻辑服务器的地址数组
**/
func GetLogicServerList() []*Model.ServerLevel {
	return logicServerList
}

/**
得到逻辑服务器的地址map数组
**/
func GetLogicServerMapList() []map[string]string {
	return logicServerMapList
}

/**
重新载入配置文件
**/
func ReParseConfig() {
	ParseXml(ConfigFile)
}

/**
获得日志文件路径
**/
func GetLogOutFile() string {
	return logFile
}

/**
获得错误日志文件路径
**/
func GetLogErrOutFile() string {
	return logErrFile
}
