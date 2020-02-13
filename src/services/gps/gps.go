package gps

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"INServer/src/services/node"
	"fmt"

	"github.com/gogo/protobuf/proto"
)

var Instance *GPS

type (
	role struct {
		address *data.RoleAddress
		uuid    string
	}

	// GPS 定位服务器 可以查询每个地图和每个角色的位置
	GPS struct {
		maps       map[string]int32
		roles      map[string]*role
		staticmaps map[int32]map[int32]string
	}
)

// New 创建定位服务器
func New() *GPS {
	g := new(GPS)
	g.maps = make(map[string]int32)
	g.roles = make(map[string]*role)
	g.staticmaps = make(map[int32]map[int32]string)
	return g
}

// Start 启动定位服务器
func (g *GPS) Start() {
	g.initMessageHandler()
}

func (g *GPS) initMessageHandler() {
	node.Net.Listen(msg.CMD_UPDATE_ROLE_ADDRESS_NTF, g.HANDLE_UPDATE_ROLE_ADDRESS_NTF)
	node.Net.Listen(msg.CMD_REMOVE_ROLE_ADDRESS_NTF, g.HANDLE_REMOVE_ROLE_ADDRESS_NTF)
	node.Net.Listen(msg.CMD_UPDATE_MAP_ADDRESS_NTF, g.onUpdateMapAddressNTF)
	node.Net.Listen(msg.CMD_REMOVE_MAP_ADDRESS_NTF, g.onRemoveMapAddressNTF)
	node.Net.Listen(msg.CMD_GET_MAP_ADDRESS_REQ, g.onGetMapLocationReq)
	node.Net.Listen(msg.CMD_UPDATE_STATIC_MAP_UUID_NTF, g.onUpdateStaticMapUUIDNTF)
	node.Net.Listen(msg.CMD_GET_STATIC_MAP_UUID_REQ, g.onGetStaticMapUUIDReq)
}

func (g *GPS) HANDLE_UPDATE_ROLE_ADDRESS_NTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.UpdateRoleAddressNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
		return
	}
	if len(ntf.RoleUUID) > 0 {
		if _, ok := g.roles[ntf.RoleUUID]; ok == false {
			address := &data.RoleAddress{
				Gate:  global.InvalidServerID,
				World: global.InvalidServerID,
			}
			g.roles[ntf.RoleUUID] = &role{
				address: address,
			}

		}
		if ntf.Address.Gate != global.InvalidServerID {
			g.roles[ntf.RoleUUID].address.Gate = ntf.Address.Gate
		}
		if ntf.Address.World != global.InvalidServerID {
			g.roles[ntf.RoleUUID].address.World = ntf.Address.World
		}
	} else {
		logger.Error("Empty RoleUUID")
	}
}

func (g *GPS) HANDLE_REMOVE_ROLE_ADDRESS_NTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.RemoveRoleAddressNTF{}
	proto.Unmarshal(buffer, ntf)
	if r, ok := g.roles[ntf.RoleUUID]; ok {
		if _, ok := g.roles[r.uuid]; ok {
			delete(g.roles, r.uuid)
		}
	}
}

func (g *GPS) onUpdateMapAddressNTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.UpdateMapAddressNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
	} else {
		g.maps[ntf.MapUUID] = ntf.ServerID
		logger.Info(fmt.Sprintf("MapAddress UUID:%s ServerID:%d", ntf.MapUUID, ntf.ServerID))
	}
}

func (g *GPS) onUpdateStaticMapUUIDNTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.UpdateStaticMapUUIDNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
	} else {
		if _, ok := g.staticmaps[ntf.ZoneID]; ok == false {
			g.staticmaps[ntf.ZoneID] = make(map[int32]string)
		}
		g.staticmaps[ntf.ZoneID][ntf.StaticMapID] = ntf.StaticMapUUID
		logger.Info(fmt.Sprintf("StaticMap ZoneID:%d StaticMapID:%d UUID:%s", ntf.ZoneID, ntf.StaticMapID, ntf.StaticMapUUID))
	}
}

func (g *GPS) onGetStaticMapUUIDReq(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.GetStaticMapUUIDResp{}
	defer node.Net.Responce(header, resp)
	req := &msg.GetStaticMapUUIDReq{}
	err := proto.Unmarshal(buffer, req)
	if err != nil {
		logger.Error(err)
	} else {
		if maps, ok := g.staticmaps[req.ZoneID]; ok {
			if uuid, ok := maps[req.StaticMapID]; ok {
				resp.StaticMapUUID = uuid
			}
		}
	}
}

func (g *GPS) onRemoveMapAddressNTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.RemoveMapAddressNTF{}
	proto.Unmarshal(buffer, ntf)
	delete(g.maps, ntf.MapUUID)
}

func (g *GPS) onGetMapLocationReq(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.GetMapAddressResp{
		ServerID: global.InvalidServerID,
	}
	defer node.Net.Responce(header, resp)
	req := &msg.GetMapAddressReq{}
	err := proto.Unmarshal(buffer, req)
	if err != nil {
		logger.Error(err)
		return
	}
	if serverID, ok := g.maps[req.MapUUID]; ok {
		resp.ServerID = serverID
	}
}
