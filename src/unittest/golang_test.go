package unittest

import (
	"INServer/src/common/logger"
	"strconv"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t1 := strconv.FormatInt(time.Now().UnixNano(), 10)
	t2 := strconv.FormatInt(time.Now().UnixNano(), 10)
	logger.Debug(t1)
	logger.Debug(t2)
}
