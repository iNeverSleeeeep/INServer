package logger

import (
	"INServer/src/common/global"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const (
	timeFormat string = "2006-01-02T15:04:05.999999+08:00"
)

var (
	day        = time.Now().Day()
	logHandler *log.Logger
	id         string

	skipstrings = []string{"/stack.go", "debug.Stack", "/protect", "panic", "logger"}

	// IsDebug 是否是调试模式运行的
	IsDebug = strings.Contains(os.Args[0], "__debug_bin")
)

// Setup 日志初始化
func Setup() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	now := time.Now()
	file := wd + "/log/" + fmt.Sprintf("%d-%d%02d%02d.log", int(global.CurrentServerID), now.Year(), now.Month(), now.Day())

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatalln(err)
	}
	logHandler = log.New(logFile, "", 0)
}

func format(level string, v ...interface{}) string {
	now := time.Now().Format(timeFormat)
	return fmt.Sprintf("{\"server\": {\"type\":\"%s\", \"id\":\"%d\"}, \"level\":\"%s\", \"@timestamp\":\"%s\", \"message\":\"%s\"}", global.CurrentServerType, global.CurrentServerID, level, now, fmt.Sprint(v...))
}

// Debug 调试日志
func Debug(v ...interface{}) {
	log.Println(v...)
	logHandler.Println(format("debug", v...))
	PrintStack()
}

// Info 普通信息日志
func Info(v ...interface{}) {
	log.Println(v...)
	logHandler.Println(format("info", v...))
}

// Error 错误信息日志
func Error(v ...interface{}) {
	log.Println(v...)
	logHandler.Println(format("error", v...))
	PrintStack()
}

// Fatal 严重错误日志
func Fatal(v ...interface{}) {
	log.Println(v...)
	logHandler.Println(format("fatal", v...))
	PrintStack()
}

// PrintStack 输出调用栈
func PrintStack() {
	stack := strings.Split(string(debug.Stack()), "\n")
	logstr := ""
	for _, str := range stack {
		skip := false
		for _, skipstr := range skipstrings {
			if strings.Contains(str, skipstr) {
				skip = true
				break
			}
		}
		if skip == false {
			logstr = logstr + str
		}
	}
	log.Println(logstr)
	logHandler.Println(format("stack", logstr))
}
