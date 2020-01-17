package protect

import (
	"INServer/src/common/logger"
	"os"
	"runtime/debug"
	"strings"
)

var skipstrings = []string{"/stack.go", "debug.Stack", "/protect", "panic"}
var isDebug = strings.Contains(os.Args[0], "__debug_bin")

func CatchPanic() {
	if isDebug {
		return
	}
	if err := recover(); err != nil {
		logger.Fatal(err)
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
		logger.Fatal(logstr)
	}
}
