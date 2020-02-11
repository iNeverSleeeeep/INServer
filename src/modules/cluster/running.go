package cluster

import (
	"INServer/src/common/global"
	"INServer/src/modules/etcmgr"
	"INServer/src/proto/etc"
	"INServer/src/proto/msg"
)

var (
	gates    []int32
	database int32
	gps      int32
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

func refreshRunning() {
	gates = make([]int32, 0)
	database = global.InvalidServerID
	gps = global.InvalidServerID
	for serverID, info := range nodes {
		serverType := etcmgr.Instance.GetServerType(int32(serverID))
		if info.NodeState == msg.NodeState_Running {
			if serverType == global.GateServer {
				gates = append(gates, int32(serverID))
			} else if serverType == global.DatabaseServer {
				database = int32(serverID)
			} else if serverType == global.GPSServer {
				gps = int32(serverID)
			}
		}
	}
}

func refreshRunningZones() {
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

func init() {
	gates = make([]int32, 0)
	zones = make([]*etc.Zone, 0)
}
