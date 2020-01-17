package center

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/modules/etcmgr"
	"INServer/src/modules/innet"
	"INServer/src/proto/etc"
	"INServer/src/proto/msg"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
)

var Instance *Center

const (
	timeout = 2000 * 1E6 // 纳秒
)

type (
	ServerInfo struct {
		Info      msg.ServerInfo
		KeepAlive int64
	}
	Center struct {
		Net       *innet.INNet
		Servers   map[int32]*ServerInfo
		realZones []*etc.Zone
	}
)

func New() *Center {
	c := new(Center)
	c.Net = innet.New()
	c.Servers = make(map[int32]*ServerInfo)
	c.realZones = make([]*etc.Zone, 0)
	c.refreshRealZones()
	c.tickServerState()
	return c
}

func (c *Center) Start() {
	c.registerListeners()
	c.Net.Start()
	logger.Debug(fmt.Sprintf("Server Start Type:%s ID:%d", global.ServerType, global.ServerID))
}

func (c *Center) registerListeners() {
	c.Net.Listen(msg.Command_SERVER_STATE, c.onServerStateChange)
	c.Net.Listen(msg.Command_KEEP_ALIVE, c.onServerKeepAlive)
	c.Net.Listen(msg.Command_RELOAD_ETC_REQ, c.onReloadETCReq)
}

func (c *Center) onServerStateChange(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.ServerStateResp{}
	defer c.Net.Responce(header, resp)
	message := &msg.ServerStateReq{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		return
	}

	serverID := message.Info.ServerID
	state := message.Info.State

	c.Servers[serverID] = &ServerInfo{
		Info:      *message.Info,
		KeepAlive: time.Now().UnixNano(),
	}

	switch state {
	case msg.ServerState_Start:
		resp.ServerType = etcmgr.Instance.GetServerType(serverID)
		resp.ServerConfig = etcmgr.Instance.GetServerConfig(serverID)
		c.Net.AddServers([]*msg.ServerInfo{message.Info})
		resp.Servers = etcmgr.Instance.Servers()
		resp.Zones = c.realZones
		c.Net.ResetServer(message.Info)
		// TODO 每个游戏区都在哪个服务器上
		resp.ZoneLocations = nil
		break
	case msg.ServerState_Running:
		c.resetServer(c.Servers[serverID])
		if etcmgr.Instance.GetServerType(serverID) == global.WorldServer {
			c.refreshRealZones()
			c.pushZonesState()
		}
		break
	}

	logger.Debug(fmt.Sprintf("ServerState ID:%d Type:%s State:%s", serverID, etcmgr.Instance.GetServerType(serverID), state.String()))
	c.printServerState()
}

func (c *Center) onServerKeepAlive(header *msg.MessageHeader, buffer []byte) {
	message := &msg.KeepAlive{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		return
	}
	serverID := message.ServerID
	info, ok := c.Servers[serverID]
	if ok {
		realive := false
		info.KeepAlive = time.Now().UnixNano()
		if info.Info.State != msg.ServerState_Running {
			info.Info.State = msg.ServerState_Running
			realive = true
			if etcmgr.Instance.GetServerType(serverID) == global.WorldServer {
				c.refreshRealZones()
				c.pushZonesState()
			}
		}
		if realive {
			// 进入这里 证明服务器的状态由离线变为在线 这样的情况 我们重置一下所有的连接状态
			c.resetServer(info)
			logger.Debug(fmt.Sprintf("ServerState ID:%d Type:%s State:%s", serverID, etcmgr.Instance.GetServerType(serverID), msg.ServerState_Running.String()))
			c.printServerState()
		}
	}
}

func (c *Center) onReloadETCReq(header *msg.MessageHeader, buffer []byte) {
	defer c.Net.ResponceBytes(header, make([]byte, 1))
	dir, _ := os.Getwd()
	etcmgr.Instance.Load(dir + "/etc")
	c.pushETCUpdate()
}

func (c *Center) tickServerState() {
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			now := time.Now().UnixNano()
			serverlist := make([]*msg.ServerInfo, 0)
			for _, info := range c.Servers {
				if info.Info.State == msg.ServerState_Running && info.KeepAlive+timeout < now {
					info.Info.State = msg.ServerState_Offline
					serverlist = append(serverlist, &info.Info)
					logger.Debug(fmt.Sprintf("ServerState ID:%d Type:%s State:%s", info.Info.ServerID, etcmgr.Instance.GetServerType(info.Info.ServerID), msg.ServerState_Offline.String()))
					c.printServerState()
					c.refreshRealZones()
					c.pushZonesState()
				}
			}
			if len(serverlist) > 0 {
				c.pushServerList(serverlist)
			}
		}
	}()
}

func (c *Center) pushServerList(servers []*msg.ServerInfo) {
	serverlist := &msg.ServerInfoList{
		Servers: servers,
	}
	for _, info := range c.Servers {
		if info.Info.State != msg.ServerState_Offline {
			c.Net.NotifyServer(msg.Command_SERVER_LIST_CHANGED, serverlist, info.Info.ServerID)
		}
	}
}

func (c *Center) resetServer(server *ServerInfo) {
	serverlist := make([]*msg.ServerInfo, 0)
	for _, info := range c.Servers {
		serverlist = append(serverlist, &info.Info)
	}
	c.Net.NotifyServer(msg.Command_SERVER_RESET, &msg.ServerInfoList{Servers: serverlist}, server.Info.ServerID)

	serverlist = make([]*msg.ServerInfo, 1)
	serverlist[0] = &server.Info
	c.pushServerList(serverlist)
}

func (c *Center) pushETCUpdate() {
	updateETC := &msg.UpdateETC{}
	updateETC.Servers = etcmgr.Instance.Servers()
	updateETC.Zones = c.realZones
	for _, info := range c.Servers {
		if info.Info.State != msg.ServerState_Offline {
			updateETC.ServerConfig = etcmgr.Instance.GetServerConfig(info.Info.ServerID)
			c.Net.NotifyServer(msg.Command_UPDATE_ETC_NTF, updateETC, info.Info.ServerID)
		}
	}
}

func (c *Center) pushZonesState() {
	updateETC := &msg.UpdateETC{}
	updateETC.Zones = c.realZones
	for _, info := range c.Servers {
		if info.Info.State != msg.ServerState_Offline {
			c.Net.NotifyServer(msg.Command_UPDATE_ETC_NTF, updateETC, info.Info.ServerID)
		}
	}
}

func (c *Center) refreshRealZones() {
	zones := make([]*etc.Zone, 0)
	for _, zone := range etcmgr.Instance.Zones() {
		realZone := &etc.Zone{}
		realZone.ZoneID = zone.ZoneID
		realZone.Name = zone.Name
		realZone.State = zone.State
		if zone.State != global.ZoneStateClosed {
			servers := etcmgr.Instance.GetZoneLocation(zone.ZoneID)
			for _, serverID := range servers {
				svr, ok := c.Servers[serverID]
				if ok == false {
					realZone.State = global.ZoneStateClosed
					break
				} else if svr.Info.State != msg.ServerState_Running {
					realZone.State = global.ZoneStateClosed
					break
				}
			}
		}
		zones = append(zones, realZone)
	}
	c.realZones = zones
}

func (c *Center) printServerState() {
	states := "[ServerState] "
	for _, svr := range etcmgr.Instance.Servers() {
		if svr.ServerID == 0 {
			continue
		}
		running := false
		start := false
		if server, ok := c.Servers[svr.ServerID]; ok {
			running = server.Info.State == msg.ServerState_Running
			start = server.Info.State == msg.ServerState_Start
		}
		if running {
			states = states + fmt.Sprintf("%d%s ", svr.ServerID, "O")
		} else if start {
			states = states + fmt.Sprintf("%d%s ", svr.ServerID, "?")
		} else {
			states = states + fmt.Sprintf("%d%s ", svr.ServerID, "X")
		}
	}
	logger.Debug(states)
}
