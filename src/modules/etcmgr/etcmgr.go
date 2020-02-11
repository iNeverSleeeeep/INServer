package etcmgr

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/etc"
	"INServer/src/proto/msg"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gogo/protobuf/proto"
)

var Instance *ETC

type (
	ETC struct {
		servers      []*etc.Server
		database     *etc.Database
		basic        *etc.BasicConfig
		zones        []*etc.Zone
		zoneLocation map[int32][]int32 // 每个游戏区都在哪些WorldServer里面
		ok           bool
	}
)

func New() *ETC {
	e := new(ETC)
	if global.CurrentServerID == global.CenterID {
		dir, err := os.Getwd()
		if err != nil {
			logger.Fatal(err)
			return nil
		}
		e.Load(dir + "/etc")
	}
	return e
}

func (e *ETC) Load(path string) {
	basic := e.loadBasic(path)
	if e.checkBasic(basic) == false {
		logger.Debug("ETC:加载basic失败")
		return
	}
	database := e.loadDatabase(path)
	if e.checkDatabase(database) == false {
		logger.Debug("ETC:加载database失败")
		return
	}
	zones := e.loadZones(path)
	if e.checkZones(zones) == false {
		logger.Debug("ETC:加载zones失败")
		return
	}
	servers := e.loadServers(path)
	if e.checkServers(servers) == false {
		logger.Debug("ETC:加载servers失败")
		return
	}
	if e.checkAllConfig(basic, database, zones, servers) == false {
		logger.Debug("ETC:检查失败")
		return
	}
	e.basic = basic
	e.database = database
	e.zones = zones
	e.servers = servers
	e.makeConfig()
	e.ok = true
}

func (e *ETC) checkAllConfig(basic *etc.BasicConfig, database *etc.Database, zones []*etc.Zone, servers []*etc.Server) bool {
	return true
}

func (e *ETC) makeConfig() {
	e.zoneLocation = make(map[int32][]int32)
	for _, svr := range e.servers {
		if svr.ServerType == global.WorldServer {
			for _, zone := range svr.ServerConfig.WorldConfig.Zones {
				_, ok := e.zoneLocation[zone.ZoneID]
				if ok == false {
					e.zoneLocation[zone.ZoneID] = make([]int32, 0)
				}
				find := false
				for _, serverID := range e.zoneLocation[zone.ZoneID] {
					if serverID == svr.ServerID {
						find = true
						break
					}
				}
				if find == false {
					e.zoneLocation[zone.ZoneID] = append(e.zoneLocation[zone.ZoneID], svr.ServerID)
				}
			}
		}
	}
}

func (e *ETC) loadServers(path string) []*etc.Server {
	f, err := os.Open(path + "/servers.json")
	if err != nil {
		logger.Debug(err)
		return nil
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Debug(err)
		return nil
	}
	servers := &etc.ServerList{}
	err = json.Unmarshal(buf, servers)
	if err != nil {
		logger.Debug(err)
	}
	return servers.Servers
}

func (e *ETC) checkServers([]*etc.Server) bool {
	return true
}

func (e *ETC) loadDatabase(path string) *etc.Database {
	f, err := os.Open(path + "/database.json")
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Fatal(err)
	}
	database := &etc.Database{}
	err = json.Unmarshal(buf, database)
	if err != nil {
		logger.Fatal(err)
	}
	return database
}

func (e *ETC) checkDatabase(*etc.Database) bool {
	return true
}

func (e *ETC) loadBasic(path string) *etc.BasicConfig {
	f, err := os.Open(path + "/basic.json")
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Fatal(err)
	}
	basic := &etc.BasicConfig{}
	err = json.Unmarshal(buf, basic)
	if err != nil {
		logger.Fatal(err)
	}
	return basic
}

func (e *ETC) checkBasic(*etc.BasicConfig) bool {
	return true
}

func (e *ETC) loadZones(path string) []*etc.Zone {
	f, err := os.Open(path + "/zones.json")
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Fatal(err)
	}
	zones := &etc.ZoneList{}
	err = json.Unmarshal(buf, zones)
	if err != nil {
		logger.Fatal(err)
	}
	return zones.Zones
}

func (e *ETC) checkZones([]*etc.Zone) bool {
	return true
}

func (e *ETC) GetServerType(serverID int32) string {
	if len(e.servers) >= int(serverID) {
		return e.servers[int(serverID)].ServerType
	}
	return global.InvalidServer
}

func (e *ETC) GetServerConfig(serverID int32) *etc.ServerConfig {
	if len(e.servers) >= int(serverID) {
		server := e.servers[int(serverID)]
		switch server.ServerType {
		case global.CenterServer:
			return server.ServerConfig
		case global.GateServer:
			return server.ServerConfig
		case global.LoginServer:
			if server.ServerConfig.LoginConfig.Database == nil {
				server.ServerConfig.LoginConfig.Database = e.database
			}
			return server.ServerConfig
		case global.DatabaseServer:
			if server.ServerConfig.DatabaseConfig.Database == nil {
				server.ServerConfig.DatabaseConfig.Database = e.database
			}
			return server.ServerConfig
		case global.WebServer:
			return server.ServerConfig
		case global.WorldServer:
			return server.ServerConfig
		case global.GPSServer:
			return server.ServerConfig
		default:
			logger.Debug("没有实现这种服务器类型:" + server.ServerType)
			return nil
		}
	}
	return nil
}

func (e *ETC) BasicConfig() *etc.BasicConfig {
	return e.basic
}

func (e *ETC) Database() *etc.Database {
	return e.database
}

func (e *ETC) Servers() []*etc.Server {
	return e.servers
}

func (e *ETC) Zones() []*etc.Zone {
	return e.zones
}

func (e *ETC) GetZoneLocation(zoneID int32) []int32 {
	if list, ok := e.zoneLocation[zoneID]; ok {
		return list
	}
	return nil
}

// OK 配置是否加载完成，对于center是读文件完成，对于node是center返回数据
func (e *ETC) OK() bool {
	return e.ok
}

func (e *ETC) HANDLE_ETC_SYNC_NTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.ETCSyncNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
		return
	}
	e.servers = ntf.ServerList.Servers
	e.database = ntf.Database
	e.basic = ntf.BasicConfig
	e.zones = ntf.ZoneList.Zones
	e.makeConfig()
	e.ok = true
}

func (e *ETC) GenerateETCSyncNTF() *msg.ETCSyncNTF {
	ntf := &msg.ETCSyncNTF{
		BasicConfig: e.BasicConfig(),
		Database:    e.Database(),
		ServerList: &etc.ServerList{
			Servers: e.Servers(),
		},
		ZoneList: &etc.ZoneList{
			Zones: e.Zones(),
		},
	}
	return ntf
}
