package Utils

import (
	"log"
	"os"
)

var LogOut *log.Logger

var LogErrOut *log.Logger

/**日志文件**/
var files os.File

/**错误日志文件**/
var errFiles os.File

type logMsg struct {
	format string
	value  []interface{}
}

var logChan chan *logMsg

/**
初始化日志文件
**/
func InitLogOut(logFile string) error {
	//设置log文件
	files, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		return LogErr(err)
	}
	LogOut = log.New(files, "", 0)
	LogOut.SetFlags(log.Ldate | log.Ltime)
	//init log chan
	logChan = make(chan *logMsg, 100)
	go writeLog()
	return nil
}

/**
初始化错误日志文件
**/
func InitLogErrOut(logErrFile string) error {
	//设置logErr文件
	errFiles, err := os.OpenFile(logErrFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		return LogErr(err)
	}
	LogErrOut = log.New(errFiles, "", 0)
	LogErrOut.SetFlags(log.Ldate | log.Ltime)
	return nil
}

/**
关闭文件
**/
func DeferFiles() {
	files.Close()
	errFiles.Close()
}

/**
写到日志chan中
**/
func LogInfo(format string, info ...interface{}) {
	// if logChan != nil {
	// 	logChan <- &logMsg{format, info}
	// } else {
	// 	log.Printf(format, info...)
	// }
	if LogOut != nil {
		LogOut.Printf(format, info...)
	} else {
		log.Printf(format, info...)
	}
}

/**
写错误日志
**/
func LogErrInfo(format string, info ...interface{}) {
	if LogErrOut != nil {
		LogErrOut.Printf(format, info...)
	} else {
		log.Printf(format, info...)
	}
}

/**
从日志chan写到日志文件中
**/
func writeLog() {
	var logInfo *logMsg
	for {
		select {
		case logInfo = <-logChan:
			if LogOut != nil {
				LogOut.Printf((*logInfo).format, (*logInfo).value...)
			} else {
				log.Printf((*logInfo).format, (*logInfo).value...)
			}
		}
	}
}
