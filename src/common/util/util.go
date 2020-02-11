package util

import (
	"INServer/src/common/logger"
	"math/rand"
	"os"
	"reflect"
	"time"
	"unsafe"
)

// GetRandomString 取得随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// SetProcessName 设置linux下进程名称
func SetProcessName(name string) error {
	argv0str := (*reflect.StringHeader)(unsafe.Pointer(&os.Args[0]))
	argv0 := (*[1 << 30]byte)(unsafe.Pointer(argv0str.Data))[:argv0str.Len]

	n := copy(argv0, name)
	if n < len(argv0) {
		argv0[n] = 0
	}

	return nil
}

// Wait 等待函数返回true
func Wait(f func() bool, message string, d time.Duration) {
	for {
		if f() {
			break
		}
		logger.Info(message)
		time.Sleep(d)
	}
}
