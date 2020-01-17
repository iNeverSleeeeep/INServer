package innet

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/msg"
	"net"
	"strconv"
)

type (
	server struct {
		addr      *net.UDPAddr
		info      *msg.ServerInfo
		packageID uint64
	}

	address struct {
		innet   *INNet
		servers map[int32]*server
		center  *server
	}
)

func newAddress(innet *INNet) *address {
	a := new(address)
	a.innet = innet
	a.servers = make(map[int32]*server)

	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:"+strconv.Itoa(int(global.CenterID)+recvport))
	if err != nil {
		logger.Fatal(err)
	}
	a.center = &server{
		addr:      addr,
		packageID: 0,
		info: &msg.ServerInfo{
			ServerID: global.CenterID,
			Address:  nil,
			State:    msg.ServerState_Running,
		},
	}
	return a
}

func (a *address) addServerList(servers []*msg.ServerInfo) {
	for _, serverToAdd := range servers {
		if serverExist, ok := a.servers[serverToAdd.ServerID]; ok {
			serverExist.info = serverToAdd
		} else {
			addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:"+strconv.Itoa(int(serverToAdd.ServerID)+recvport))
			if err != nil {
				logger.Debug(err)
				continue
			}
			a.servers[serverToAdd.ServerID] = &server{
				addr:      addr,
				info:      serverToAdd,
				packageID: 0,
			}
		}
		// 服务器离线了 要重置的一些状态
		if serverToAdd.State == msg.ServerState_Offline {
			a.resetServer(serverToAdd.ServerID)
			a.innet.receiver.resetServer(serverToAdd.ServerID)
			a.innet.retry.resetServer(serverToAdd.ServerID)
		}
	}
}

func (a *address) resetServer(serverID int32) {
	if svr, ok := a.servers[serverID]; ok {
		svr.packageID = 0
	}
}

func (a *address) getByCommand(command msg.Command) *server {
	switch command {
	case msg.Command_SERVER_STATE, msg.Command_KEEP_ALIVE, msg.Command_RELOAD_ETC_REQ:
		return a.center
	case msg.Command_LD_CREATE_PLAYER_REQ, msg.Command_GD_LOAD_PLAYER_REQ, msg.Command_GD_RELEASE_PLAYER_NTF, msg.Command_GD_CREATE_ROLE_REQ, msg.Command_GD_LOAD_ROLE_REQ:
		if svr, ok := a.servers[a.innet.database]; ok {
			return svr
		}
	}
	return nil
}

func (a *address) getByServerID(serverID int32) *server {
	if serverID == global.CenterID {
		return a.center
	} else if info, ok := a.servers[serverID]; ok {
		return info
	}
	return nil
}
