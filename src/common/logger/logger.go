package logger

import (
	"INServer/src/common/global"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	timeFormat string = "2006-01-02T15:04:05.999999+08:00"
)

var (
	day        = time.Now().Day()
	logHandler *log.Logger
	id         string
)

func Setup() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	now := time.Now()
	file := wd + "/log/" + fmt.Sprintf("%d-%d%02d%02d.log", int(global.ServerID), now.Year(), now.Month(), now.Day())

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatalln(err)
	}
	logHandler = log.New(logFile, "", 0)
}

func format(level string, v ...interface{}) string {
	now := time.Now().Format(timeFormat)
	return fmt.Sprintf("{\"server\": {\"type\":\"%s\", \"id\":\"%d\"}, \"level\":\"%s\" \"@timestamp\":\"%s\", \"message\":\"%s\"}", global.ServerType, global.ServerID, level, now, fmt.Sprint(v...))
}

func Debug(v ...interface{}) {
	log.Println(v...)
	logHandler.Println(format("[debug]", v...))
}

func Info(v ...interface{}) {
	log.Println(v...)
	logHandler.Println(format("[info]", v...))
}

func Error(v ...interface{}) {
	log.Println(v...)
	logHandler.Println(format("[error]", v...))
}

func Fatal(v ...interface{}) {
	log.Println(v...)
	logHandler.Println(format("[fatal]", v...))
}
