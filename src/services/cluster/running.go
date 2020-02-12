package cluster

import (
	"INServer/src/common/global"
	"INServer/src/modules/etcmgr"
	"INServer/src/proto/etc"
	"INServer/src/proto/msg"
)

var (
	running  []int32
	gates    []int32
	database int32
	gps      int32
	balcony  int32
	zones    []*etc.Zone
)

// RunningGates 处于Running状态下的门服务器
func RunningGates() []int32 {
	return gates
}

// RunningDatabase 处于Running状态下的数据库服务器
func RunningDatabase() int32 {
	return database
}

// RunningGPS 处于Running状态下的定位服务器
func RunningGPS() int32 {
	return gps
}

// RunningBalcony 处于Runing状态下的月台服务器
func RunningBalcony() int32 {
	return balcony
}

// RefreshRunning 刷新活跃中的服务器
func RefreshRunning() {
	running = make([]int32, 0)
	gates = make([]int32, 0)
	database = global.InvalidServerID
	gps = global.InvalidServerID
	balcony = global.InvalidServerID
	for serverID, info := range nodes {
		serverType := etcmgr.Instance.GetServerType(int32(serverID))
		if info.NodeState == msg.NodeState_Running {
			running = append(running, int32(serverID))
			if serverType == global.GateServer {
				gates = append(gates, int32(serverID))
			} else if serverType == global.DatabaseServer {
				database = int32(serverID)
			} else if serverType == global.GPSServer {
				gps = int32(serverID)
			} else if serverType == global.BalconyServer {
				balcony = int32(serverID)
			}
		}
	}
}

// RefreshRunningZones 刷新游戏区状态
func RefreshRunningZones() {
	zones = make([]*etc.Zone, 0)
	for _, zone := range etcmgr.Instance.Zones() {
		realZone := &etc.Zone{}
		realZone.ZoneID = zone.ZoneID
		realZone.Name = zone.Name
		realZone.State = zone.State
		if zone.State != global.ZoneStateClosed {
			servers := etcmgr.Instance.GetZoneLocation(zone.ZoneID)
			for _, serverID := range servers {
				if GetNodeState(serverID) != msg.NodeState_Running {
					realZone.State = global.ZoneStateClosed
					break
				}
			}
		}
		zones = append(zones, realZone)
	}
}

// RunningCount 运行中服务器数量
func RunningCount() int {
	return len(running)
}

func init() {
	running = make([]int32, 0)
	gates = make([]int32, 0)
	database = global.InvalidServerID
	gps = global.InvalidServerID
	balcony = global.InvalidServerID
	zones = make([]*etc.Zone, 0)
}
