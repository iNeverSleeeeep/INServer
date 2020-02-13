package startup

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/util"
	"INServer/src/services/balcony"
	"INServer/src/services/center"
	"INServer/src/services/cluster"
	"INServer/src/services/database"
	"INServer/src/services/etcmgr"
	"INServer/src/services/gate"
	"INServer/src/services/gps"
	"INServer/src/services/login"
	"INServer/src/services/node"
	"INServer/src/services/web"
	"INServer/src/services/world"
	"fmt"
	"strings"
	"time"
)

// Run 服务器启动流程
func Run() {
	if global.CurrentServerID == global.CenterID {
		go startCenter()
	} else {
		go startNode()
	}
}

func startCenter() {
	global.CurrentServerType = global.CenterServer
	util.SetProcessName(fmt.Sprintf("%s@in-%d         ", strings.ToLower(global.CurrentServerType), global.CurrentServerID))
	etcmgr.Instance = etcmgr.New()
	global.CurrentServerConfig = etcmgr.Instance.GetServerConfig(global.CurrentServerID)
	center.Instance = center.New()
	center.Instance.Start()
}

func startNode() {
	etcmgr.Instance = etcmgr.New()
	node.Instance = node.New()
	for {
		if etcmgr.Instance.OK() {
			break
		} else {
			node.Instance.Net.SendNodeStartNTF()
			logger.Info("等待中心服启动完成...")
			time.Sleep(time.Second)
		}
	}

	global.CurrentServerType = etcmgr.Instance.GetServerType(global.CurrentServerID)
	global.CurrentServerConfig = etcmgr.Instance.GetServerConfig(global.CurrentServerID)
	util.SetProcessName(fmt.Sprintf("%s@in-%d          ", strings.ToLower(global.CurrentServerType), global.CurrentServerID))
	node.Instance.Prepare()
	startServer()
	node.Instance.Start()
	logger.Info(fmt.Sprintf("服务器启动完成 ID:%d Type:%s", global.CurrentServerID, global.CurrentServerType))
}

func startServer() {
	switch global.CurrentServerType {
	case global.GateServer:
		gate.Instance = gate.New()
		gate.Instance.Start()
		break
	case global.LoginServer:
		login.Instance = login.New()
		login.Instance.Start()
	case global.DatabaseServer:
		database.Instance = database.New()
		database.Instance.Start()
	case global.WebServer:
		web.Instance = web.New()
		web.Instance.Start()
	case global.WorldServer:
		util.Wait(func() bool {
			return cluster.RunningDatabase() != global.InvalidServerID
		}, "等待数据库服务器启动完成...", time.Second)
		util.Wait(func() bool {
			return cluster.RunningGPS() != global.InvalidServerID
		}, "等待定位服务器启动完成...", time.Second)
		world.Instance = world.New()
		world.Instance.Start()
	case global.GPSServer:
		gps.Instance = gps.New()
		gps.Instance.Start()
	case global.BalconyServer:
		balcony.Instance = balcony.New()
		balcony.Instance.Start()
	default:
		logger.Fatal("不支持的服务器类型:" + global.CurrentServerType)
	}
}
