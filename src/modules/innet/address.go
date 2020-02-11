package innet

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/modules/cluster"
	"INServer/src/proto/msg"
	"fmt"
	"net"
	"strconv"
)

type (
	server struct {
		addr      *net.UDPAddr
		packageID uint64
		id        int32
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
		id:        global.CenterID,
	}
	return a
}

func (a *address) refresh() {
	servers := cluster.GetNodes()
	for serverID, info := range servers {
		if info.NodeAddress != nil && len(info.NodeAddress) > 0 {
			var svr *server
			var ok = false
			if svr, ok = a.servers[int32(serverID)]; ok == false {
				svr = &server{
					packageID: 0,
					id:        int32(serverID),
				}
				a.servers[svr.id] = svr
			}
			ip := &net.IPAddr{IP: info.NodeAddress}
			addr := &net.UDPAddr{IP: ip.IP, Port: serverID + recvport, Zone: ip.Zone}
			svr.addr = addr
		}
	}
}

func (a *address) resetServer(serverID int32) {
	if svr, ok := a.servers[serverID]; ok {
		svr.packageID = 0
	}
}

func (a *address) getByCommand(command msg.CMD) *server {
	switch command {
	case msg.CMD_SERVER_STATE, msg.CMD_KEEP_ALIVE, msg.CMD_RELOAD_ETC_REQ:
		return a.center
	case msg.CMD_LD_CREATE_PLAYER_REQ,
		msg.CMD_GD_LOAD_PLAYER_REQ,
		msg.CMD_GD_RELEASE_PLAYER_NTF,
		msg.CMD_GD_CREATE_ROLE_REQ,
		msg.CMD_GD_LOAD_ROLE_REQ,
		msg.CMD_LOAD_STATIC_MAP_REQ,
		msg.CMD_SAVE_STATIC_MAP_REQ:
		serverID := cluster.RunningDatabase()
		if serverID != global.InvalidServerID {
			if svr, ok := a.servers[serverID]; ok {
				return svr
			}
		}
		logger.Error(fmt.Sprintf("没有找到Database服务器"))
		break
	case msg.CMD_GET_MAP_ADDRESS_REQ,
		msg.CMD_GET_STATIC_MAP_UUID_REQ,
		msg.CMD_UPDATE_MAP_ADDRESS_NTF,
		msg.CMD_UPDATE_PLAYER_ADDRESS_NTF,
		msg.CMD_UPDATE_STATIC_MAP_UUID_NTF:
		serverID := cluster.RunningGPS()
		if serverID != global.InvalidServerID {
			if svr, ok := a.servers[serverID]; ok {
				return svr
			}
		}
		logger.Error(fmt.Sprintf("没有找到GPS服务器"))
		break
	case msg.CMD_ROLE_ENTER:
		break
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
