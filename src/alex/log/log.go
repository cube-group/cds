// 日志封装
package log

import (
	"time"
	"fmt"
	"os"
	"strings"
	"strconv"
	"alex/utils"
)

const (
	LEVEL_DEBUG = "DEBUG"
	LEVEL_INFO = "INFO"
	LEVEL_WARN = "WARN"
	LEVEL_ERROR = "ERROR"
	LEVEL_FATAL = "FATAL"
)

// 应用名称
var logApp string = ""
// 是否支持debug日志
var logDebug bool = false
// 文件指针
var logFile *os.File
// Logger类
type Logger struct{}
// 全局Logger
var Log *Logger

func NewLogger(app string, path string, locate string, debug bool) *Logger {
	if Log == nil {
		Log = new(Logger)
		Log.Init(app, path, locate, debug)
	}
	return Log
}

// 初始化日志系统
// app app名称
// path 日志文件路径
// 时区设置
// debug 是否为debug模式
func (l *Logger)Init(app string, path string, locate string, debug bool) {
	if locate == "" {
		locate = "Asia/Shanghai"
	}
	_, err := time.LoadLocation(locate)
	if err != nil {
		panic("log locate error")
	}
	if logApp != "" {
		panic("log app is exsits")
	}
	if path == "" || !utils.FileIsExist(path) {
		panic(utils.StringTrace(path, "log path is nil"))
	}

	path = utils.StringJoin("", path, "/", utils.GetFormatYmd(), ".txt")
	f1, err1 := os.OpenFile(path, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0666)
	if err1 != nil && os.IsNotExist(err1) {
		f2, err2 := os.Create(path)
		if err2 != nil {
			panic(err2.Error())
		} else {
			logFile = f2
		}
	} else {
		logFile = f1
	}

	logApp = app
	logDebug = debug
}

// Info日志
func (l *Logger)Info(key string, value ... interface{}) error {
	return logAppend(LEVEL_INFO, key, value...)
}

// Debug日志
func (l *Logger)Debug(key string, value ... interface{}) error {
	if logDebug == false {
		return nil
	}
	return logAppend(LEVEL_DEBUG, key, value...)
}

// Warn日志
func (l *Logger)Warn(key string, value ... interface{}) error {
	return logAppend(LEVEL_WARN, key, value...)
}

// Error日志
func (l *Logger)Error(key string, value ... interface{}) error {
	return logAppend(LEVEL_ERROR, key, value...)
}

// Fatal日志
func (l *Logger)Fatal(key string, value ... interface{}) error {
	return logAppend(LEVEL_FATAL, key, value...)
}

// 向日志文件中追加日志
// logType 日志的类型
// key 日志的所属文件名称或者功能标志
// requestId http连续性请求标志
func logAppend(logType string, key string, value ... interface{}) error {
	requestId := utils.GetShortUUID()

	var stringArr []string
	stringArr = append(stringArr, logGetFormatString(utils.GetFormatYmdHis()))
	stringArr = append(stringArr, logGetFormatString(strconv.Itoa(os.Getpid())))
	stringArr = append(stringArr, logGetFormatString(key))
	stringArr = append(stringArr, logGetFormatString(logType))
	stringArr = append(stringArr, fmt.Sprintf("(%s)", requestId))
	jsonString, _ := utils.JsonEncode(value)
	stringArr = append(stringArr, jsonString)

	stringContent := strings.Join(stringArr, "")
	logFileContent := stringContent + "\n"
	fmt.Println(stringContent)
	_, err := logFile.Write([]byte(logFileContent))
	return err
}

// 生成标准日志结构体
func logGetFormatString(value string) string {
	return fmt.Sprintf("[%s]", value)
}