package gps

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/modules/node"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"github.com/gogo/protobuf/proto"
	"fmt"
)

var Instance *GPS

type (
	player struct {
		address  *data.PlayerAddress
		roleUUID string
	}

	// GPS 定位服务器 可以查询每个地图和每个角色的位置
	GPS struct {
		maps       map[string]int32
		players    map[string]*player
		roles      map[string]string
		staticmaps map[int32]map[int32]string
	}
)

// New 创建定位服务器
func New() *GPS {
	g := new(GPS)
	g.maps = make(map[string]int32)
	g.players = make(map[string]*player)
	g.roles = make(map[string]string)
	g.staticmaps = make(map[int32]map[int32]string)
	return g
}

// Start 启动定位服务器
func (g *GPS) Start() {
	g.initMessageHandler()
}

func (g *GPS) initMessageHandler() {
	node.Instance.Net.Listen(msg.Command_UPDATE_PLAYER_ADDRESS_NTF, g.onUpdatePlayerAddressNTF)
	node.Instance.Net.Listen(msg.Command_REMOVE_PLAYER_ADDRESS_NTF, g.onRemovePlayerAddressNTF)
	node.Instance.Net.Listen(msg.Command_UPDATE_MAP_ADDRESS_NTF, g.onUpdateMapAddressNTF)
	node.Instance.Net.Listen(msg.Command_REMOVE_MAP_ADDRESS_NTF, g.onRemoveMapAddressNTF)
	node.Instance.Net.Listen(msg.Command_GET_MAP_ADDRESS_REQ, g.onGetMapLocationReq)
	node.Instance.Net.Listen(msg.Command_UPDATE_STATIC_MAP_UUID_NTF, g.onUpdateStaticMapUUIDNTF)
	node.Instance.Net.Listen(msg.Command_GET_STATIC_MAP_UUID_REQ, g.onGetStaticMapUUIDReq)
}

func (g *GPS) onUpdatePlayerAddressNTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.UpdatePlayerAddressNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
		return
	}
	if len(ntf.PlayerUUID) > 0 {
		if _, ok := g.players[ntf.PlayerUUID]; ok == false {
			address := &data.PlayerAddress{
				Gate:   global.InvalidServerID,
				Entity: global.InvalidServerID,
			}
			g.players[ntf.PlayerUUID] = &player{
				address: address,
			}

		}
		if ntf.Address.Gate != global.InvalidServerID {
			g.players[ntf.PlayerUUID].address.Gate = ntf.Address.Gate
		}
		if ntf.Address.Entity != global.InvalidServerID {
			g.players[ntf.PlayerUUID].address.Entity = ntf.Address.Entity
		}
		if len(ntf.RoleUUID) > 0 {
			g.players[ntf.PlayerUUID].roleUUID = ntf.RoleUUID
			g.roles[ntf.RoleUUID] = ntf.PlayerUUID
		}
	} else {
		logger.Error("Empty PlayerUUID")
	}
}

func (g *GPS) onRemovePlayerAddressNTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.RemovePlayerAddressNTF{}
	proto.Unmarshal(buffer, ntf)
	if p, ok := g.players[ntf.PlayerUUID]; ok {
		if _, ok := g.roles[p.roleUUID]; ok {
			delete(g.roles, p.roleUUID)
		}
		delete(g.players, ntf.PlayerUUID)
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
	defer node.Instance.Net.Responce(header, resp)
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
	defer node.Instance.Net.Responce(header, resp)
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
