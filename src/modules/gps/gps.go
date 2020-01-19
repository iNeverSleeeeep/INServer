package gps

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/modules/node"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"github.com/gogo/protobuf/proto"
)

type (
	player struct {
		address  *data.PlayerAddress
		roleUUID string
	}

	GPS struct {
		maps    map[string]int32
		players map[string]*player
		roles   map[string]string
	}
)

func New() *GPS {
	g := new(GPS)
	g.maps = make(map[string]int32)
	g.players = make(map[string]*player)
	g.roles = make(map[string]string)
	return g
}

func (g *GPS) InitMessageHandler() {
	node.Instance.Net.Listen(msg.Command_UPDATE_PLAYER_ADDRESS_NTF, g.onUpdatePlayerAddressNTF)
	node.Instance.Net.Listen(msg.Command_REMOVE_PLAYER_ADDRESS_NTF, g.onRemovePlayerAddressNTF)
	node.Instance.Net.Listen(msg.Command_UPDATE_MAP_ADDRESS_NTF, g.onUpdateMapAddressNTF)
	node.Instance.Net.Listen(msg.Command_REMOVE_MAP_ADDRESS_NTF, g.onRemoveMapAddressNTF)
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
		g.players[ntf.PlayerUUID].address.Gate = ntf.Address.Gate
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
	proto.Unmarshal(buffer, ntf)
	g.maps[ntf.MapUUID] = ntf.ServerID
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
	proto.Unmarshal(buffer, req)
	if serverID, ok := g.maps[req.MapUUID]; ok {
		resp.ServerID = serverID
	}
}
