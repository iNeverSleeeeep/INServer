package finalize

import (
	"INServer/src/common/global"
	"INServer/src/modules/world"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Stop 关服
var Stop chan bool

// Wait 等待结束
func Wait() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	for {
		stopped := false
		select {
		case <-Stop:
			break
		case sig := <-sigs:
			if sig.String() == "interrupt" {
				stopped = true
				stopNode()
			}
			break
		}
		if stopped {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func stopNode() {
	stopServer()
}

func stopServer() {
	switch global.CurrentServerType {
	case global.WorldServer:
		world.Instance.Stop()
		break
	case global.DatabaseServer:
		<-Stop
		break
	case global.CenterServer:
		<-Stop
		break
	}
}

func init() {
	Stop = make(chan bool)
}
