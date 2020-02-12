package main

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/profiler"
	"INServer/src/lifetime/finalize"
	"INServer/src/lifetime/startup"
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

var serverID = flag.Int("id", -1, "本服务器ID(范围0~65535)")
var centerIP = flag.String("center", "127.0.0.1", "中心服务器IP")

func main() {
	runtime.GOMAXPROCS(1)
	flag.Parse()
	global.CurrentServerID = int32(*serverID)
	global.CenterIP = *centerIP

	if global.CurrentServerID == -1 {
		log.Fatalln("必须使用参数(-id ?)指定本服务器ID")
	} else if global.CurrentServerID > global.SERVER_ID_MAX || global.CurrentServerID < 0 {
		log.Fatalln("服务器ID范围0~999")
	}

	go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", 8000+global.CurrentServerID), nil)

	logger.Setup()

	profiler.Start()

	startup.Run()

	finalize.Wait()

	logger.Info(fmt.Sprintf("%d-%s Shut Down!", global.CurrentServerID, global.CurrentServerType))
}
