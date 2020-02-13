package world

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/uuid"
	"INServer/src/gameplay/ecs"
	"INServer/src/gameplay/gamemap"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"INServer/src/services/node"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
)

var Instance *World

type (
	World struct {
		gameMaps map[string]*gamemap.Map
		roles    map[string]*data.Role
		roleGate map[string]int32
	}
)

func New() *World {
	w := new(World)
	w.gameMaps = make(map[string]*gamemap.Map)
	w.roles = make(map[string]*data.Role)
	w.roleGate = make(map[string]int32)
	w.initMessageHandler()
	return w
}

func (w *World) Start() {
	req := &msg.LoadStaticMapReq{}
	resp := &msg.LoadStaticMapResp{}
	for _, zoneConfig := range global.CurrentServerConfig.WorldConfig.Zones {
		for _, gameMapID := range zoneConfig.StaticMaps {
			req.ZoneID = zoneConfig.ZoneID
			req.StaticMapID = gameMapID
			resp.Reset()
			var mapData *data.MapData
			err := node.Net.Request(msg.CMD_LOAD_STATIC_MAP_REQ, req, resp)
			if err != nil {
				logger.Error(err)
			} else if resp.Map == nil {
				mapData = &data.MapData{}
				mapData.MapID = gameMapID
				mapData.MapUUID = uuid.New()
				mapData.LastTickTime = time.Now().UnixNano()
			} else {
				mapData = resp.Map
			}

			if mapData != nil {
				staticMap := gamemap.NewMap(nil, mapData)
				w.gameMaps[mapData.MapUUID] = staticMap
				updateMapAddress := &msg.UpdateMapAddressNTF{
					MapUUID:  mapData.MapUUID,
					ServerID: global.CurrentServerID,
				}
				node.Net.Notify(msg.CMD_UPDATE_MAP_ADDRESS_NTF, updateMapAddress)
				udpateStatcMapUUID := &msg.UpdateStaticMapUUIDNTF{
					ZoneID:        zoneConfig.ZoneID,
					StaticMapID:   gameMapID,
					StaticMapUUID: mapData.MapUUID,
				}
				node.Net.Notify(msg.CMD_UPDATE_STATIC_MAP_UUID_NTF, udpateStatcMapUUID)
			}
		}
	}
}

func (w *World) Stop() {
	staticMaps := make([]*data.MapData, 0)
	for _, gamemap := range w.gameMaps {
		staticMaps = append(staticMaps, gamemap.MapData())
	}
	saveStaticMapReq := &msg.SaveStaticMapReq{
		StaticMaps: staticMaps,
	}
	saveStaticMapResp := &msg.SaveStaticMapResp{}
	err := node.Net.Request(msg.CMD_SAVE_STATIC_MAP_REQ, saveStaticMapReq, saveStaticMapResp)
	if err != nil {
		logger.Error(err)
	}

	roles := make([]*data.Role, 0)
	for _, role := range w.roles {
		roles = append(roles, role)
	}

	saveRoleReq := &msg.SaveRoleReq{
		Roles: roles,
	}
	saveRoleResp := &msg.SaveRoleResp{}
	err = node.Net.Request(msg.CMD_SAVE_ROLE_REQ, saveRoleReq, saveRoleResp)
	if err != nil {
		logger.Error(err)
	}
}

func (w *World) initMessageHandler() {
	node.Net.Listen(msg.CMD_ROLE_ENTER, w.onRoleEnterNTF)
	node.Net.Listen(msg.CMD_GET_MAP_ID, w.HANDLE_GET_MAP_ID)
	node.Net.Listen(msg.CMD_MOVE_INF, w.onRoleMoveINF)
}

// GetMap 根据UUID返回Map实例
func (w *World) GetMap(uuid string) *gamemap.Map {
	if result, ok := w.gameMaps[uuid]; ok {
		return result
	}
	return nil
}

func (w *World) onRoleEnterNTF(header *msg.MessageHeader, buffer []byte) {
	roleEnterNTF := &msg.RoleEnterNTF{}
	err := proto.Unmarshal(buffer, roleEnterNTF)
	if err != nil {
		logger.Error(err)
		return
	}
	role := roleEnterNTF.Role
	w.roleGate[role.SummaryData.GetRoleUUID()] = roleEnterNTF.Gate

	if gameMap, ok := w.gameMaps[role.SummaryData.GetMapUUID()]; ok {
		entity := ecs.NewEntity(role.OnlineData.EntityData, data.EntityType_RoleEntity)
		gameMap.EntityEnter(role.SummaryData.RoleUUID, entity)

		ntf := &msg.UpdateRoleAddressNTF{
			RoleUUID: role.SummaryData.RoleUUID,
			Address: &data.RoleAddress{
				Gate:  global.InvalidServerID,
				World: global.CurrentServerID,
			},
		}
		node.Net.Notify(msg.CMD_UPDATE_ROLE_ADDRESS_NTF, ntf)
	} else {
		logger.Error(fmt.Sprintf("角色进入失败，地图不存在 role:%s map:%s", role.SummaryData.RoleUUID, role.SummaryData.MapUUID))
	}
}

func (w *World) HANDLE_GET_MAP_ID(header *msg.MessageHeader, buffer []byte) {
	req := &msg.GetMapIDReq{}
	resp := &msg.GetMapIDResp{}
	defer node.Net.Responce(header, resp)
	err := proto.Unmarshal(buffer, req)
	if err != nil {
		logger.Error(err)
		return
	}

	if gameMap, ok := w.gameMaps[req.MapUUID]; ok {
		resp.MapID = gameMap.MapData().MapID
	} else {
		logger.Error("这个地图不在当前服务器")
	}
}

func (w *World) onRoleMoveINF(header *msg.MessageHeader, buffer []byte) {
	inf := &msg.MoveINF{}
	err := proto.Unmarshal(buffer, inf)
	if err != nil {
		logger.Error(err)
		return
	}
	if role, ok := w.roles[header.RoleUUID]; ok {
		if gameMap, ok2 := w.gameMaps[role.SummaryData.GetMapUUID()]; ok2 {
			gameMap.OnRoleMoveINF(role, inf)
		}
	}
}

func (w *World) RoleGate(uuid string) int32 {
	if gate, ok := w.roleGate[uuid]; ok {
		return gate
	}
	return global.InvalidServerID
}
