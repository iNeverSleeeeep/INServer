package world

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/uuid"
	"INServer/src/gameplay/ecs"
	"INServer/src/gameplay/gamemap"
	"INServer/src/gameplay/matrix"
	"INServer/src/proto/config"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"INServer/src/services/node"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
)

var Instance *World

type (
	// World 是世界服务 管理角色在游戏中的主要逻辑
	World struct {
		gameMaps map[string]*gamemap.Map
		roles    map[string]*data.Role
		roleGate map[string]int32
		matrix   *matrix.Matrix
	}
)

func New() *World {
	w := new(World)
	w.gameMaps = make(map[string]*gamemap.Map)
	w.roles = make(map[string]*data.Role)
	w.roleGate = make(map[string]int32)
	w.matrix = matrix.New()
	w.initMessageHandler()
	global.RoleGateGetter = w
	return w
}

func (w *World) Start() {
	req := &msg.LoadStaticMapReq{}
	resp := &msg.LoadStaticMapResp{}
	for _, zoneConfig := range global.CurrentServerConfig.WorldConfig.Zones {
		for _, gameMapID := range zoneConfig.StaticMaps {
			logger.Info(fmt.Sprintf("加载地图 游戏区:%d-%s, 地图ID:%d", zoneConfig.ZoneID, zoneConfig.ZoneID, gameMapID))
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
				staticMap := gamemap.NewMap(&config.Map{}, mapData)
				w.gameMaps[mapData.MapUUID] = staticMap
				w.matrix.OnGamemapCreate(staticMap)
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

	go func() {
		for global.PendingExit == false {
			w.tick()
			time.Sleep(time.Millisecond * time.Duration(1000/int(global.CurrentServerConfig.WorldConfig.FPS)))
		}
	}()
}

func (w *World) Stop() {
	logger.Info("保存地图中...")
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

	logger.Info("保存玩家角色中...")
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

func (w *World) tick() {

}

func (w *World) initMessageHandler() {
	node.Net.Listen(msg.CMD_ROLE_ENTER, w.HANDLE_ROLE_ENTER)
	node.Net.Listen(msg.CMD_GET_MAP_ID, w.HANDLE_GET_MAP_ID)
	node.Net.Listen(msg.CMD_MOVE_INF, w.HANDLE_MOVE_INF)
	node.Net.Listen(msg.CMD_STOP_MOVE_INF, w.HANDLE_STOP_MOVE_INF)
	node.Net.Listen(msg.CMD_ROLE_LEAVE_REQ, w.HANDLE_ROLE_LEAVE_REQ)
}

// GetMap 根据UUID返回Map实例
func (w *World) GetMap(uuid string) *gamemap.Map {
	if result, ok := w.gameMaps[uuid]; ok {
		return result
	}
	return nil
}

func (w *World) HANDLE_ROLE_ENTER(header *msg.MessageHeader, buffer []byte) {
	roleEnterNTF := &msg.RoleEnterNTF{}
	err := proto.Unmarshal(buffer, roleEnterNTF)
	if err != nil {
		logger.Error(err)
		return
	}
	role := roleEnterNTF.Role
	w.roleGate[role.SummaryData.GetRoleUUID()] = roleEnterNTF.Gate

	if gameMap, ok := w.gameMaps[role.SummaryData.GetMapUUID()]; ok {
		w.roles[role.SummaryData.RoleUUID] = role
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

func (w *World) HANDLE_ROLE_LEAVE_REQ(header *msg.MessageHeader, buffer []byte) {
	req := &msg.RoleLeaveReq{}
	resp := &msg.RoleLeaveResp{}
	defer node.Net.Responce(header, resp)
	err := proto.Unmarshal(buffer, req)
	if err != nil {
		logger.Error(err)
		return
	}

	roles := make([]*data.Role, 0)
	for _, uuid := range req.Roles {
		if role, ok := w.roles[uuid]; ok {
			roles = append(roles, role)
			if gameMap, ok2 := w.gameMaps[role.SummaryData.GetMapUUID()]; ok2 {
				gameMap.EntityLeave(uuid)
			}
		} else {
			logger.Error("角色不在当前服务器 %s", uuid)
		}
		ntf := &msg.RemoveRoleAddressNTF{
			RoleUUID: uuid,
		}
		node.Net.Notify(msg.CMD_REMOVE_ROLE_ADDRESS_NTF, ntf)
	}

	saveRoleReq := &msg.SaveRoleReq{
		Roles: roles,
	}
	saveRoleResp := &msg.SaveRoleResp{}
	err = node.Net.Request(msg.CMD_SAVE_ROLE_REQ, saveRoleReq, saveRoleResp)
	if err != nil {
		logger.Error(err)
	} else {
		for _, uuid := range req.Roles {
			delete(w.roles, uuid)
		}
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

func (w *World) HANDLE_MOVE_INF(header *msg.MessageHeader, buffer []byte) {
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

func (w *World) HANDLE_STOP_MOVE_INF(header *msg.MessageHeader, buffer []byte) {
	inf := &msg.StopMoveINF{}
	err := proto.Unmarshal(buffer, inf)
	if err != nil {
		logger.Error(err)
		return
	}
	if role, ok := w.roles[header.RoleUUID]; ok {
		if gameMap, ok2 := w.gameMaps[role.SummaryData.GetMapUUID()]; ok2 {
			gameMap.OnRoleStopMoveINF(role, inf)
		}
	}
}

// GetRoleGate 取得角色所在的门服务器
func (w *World) GetRoleGate(uuid string) int32 {
	if gate, ok := w.roleGate[uuid]; ok {
		return gate
	}
	return global.InvalidServerID
}
