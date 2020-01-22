package protect

import (
	"INServer/src/common/logger"
)

// CatchPanic 捕获异常
func CatchPanic() {
	if logger.IsDebug {
		return
	}
	if err := recover(); err != nil {
		logger.Fatal(err)
	}
}
