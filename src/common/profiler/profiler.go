package profiler

import (
	"INServer/src/common/global"
	"expvar"
	"net/http"
	"strconv"
	"time"
)

var (
	messages   map[uint64]int64
	expvarport       = 13000
	enabled    bool  = false
	messageRT  int64 = 0
)

func BeginSampleMessage(id uint64) {
	if enabled {
		messages[id] = time.Now().UnixNano()
	}
}

func EndSampleMessage(id uint64) {
	if enabled {
		messageRT += time.Now().UnixNano() - messages[id]
		delete(messages, id)
	}
}

func getRT() interface{} {
	rt := messageRT
	messageRT = 0
	return rt
}

func getServerID() interface{} {
	return global.ServerID
}

func getServerType() interface{} {
	return global.ServerType
}

func Start() {
	expvar.Publish("Message RT", expvar.Func(getRT))
	expvar.Publish("ServerID", expvar.Func(getServerID))
	expvar.Publish("ServerType", expvar.Func(getServerType))
	go http.ListenAndServe(":"+strconv.Itoa(expvarport+int(global.ServerID)), nil)
	enabled = true
}

func init() {
	messages = make(map[uint64]int64)
}
