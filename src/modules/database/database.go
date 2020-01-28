package database

import (
	"INServer/src/common/dbobj"
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/uuid"
	"INServer/src/dao"
	"INServer/src/modules/node"
	"INServer/src/proto/data"
	"INServer/src/proto/db"
	"INServer/src/proto/msg"
	"fmt"

	"github.com/gogo/protobuf/proto"
)

var Instance *Database

type (
	Database struct {
		DB                *dbobj.DBObject
		roleSummary       map[string]*data.RoleSummaryData
		roleSummaryByName map[string]*data.RoleSummaryData
		players           map[string]*data.Player
		staticMaps        map[int32]map[int32]*data.MapData
	}
)

func New() *Database {
	d := new(Database)
	d.roleSummary = make(map[string]*data.RoleSummaryData)
	d.roleSummaryByName = make(map[string]*data.RoleSummaryData)
	d.players = make(map[string]*data.Player)
	d.staticMaps = make(map[int32]map[int32]*data.MapData)
	d.DB = dbobj.New()
	d.DB.Open(global.ServerConfig.DatabaseConfig.Database, global.DatabaseSchema)
	d.loadAllRoleSummaryData()
	d.loadAllStaticMapsData()
	return d
}

func (d *Database) Start() {
	node.Instance.Net.Listen(msg.Command_LD_CREATE_PLAYER_REQ, d.onCreatePlayerReq)
	node.Instance.Net.Listen(msg.Command_GD_LOAD_PLAYER_REQ, d.onLoadPlayerReq)
	node.Instance.Net.Listen(msg.Command_GD_RELEASE_PLAYER_NTF, d.onReleasePlayerNtf)
	node.Instance.Net.Listen(msg.Command_GD_CREATE_ROLE_REQ, d.onCreateRoleReq)
	node.Instance.Net.Listen(msg.Command_GD_LOAD_ROLE_REQ, d.onLoadRoleReq)
	node.Instance.Net.Listen(msg.Command_LOAD_STATIC_MAP_REQ, d.onLoadStaticMapReq)
	node.Instance.Net.Listen(msg.Command_SAVE_STATIC_MAP_REQ, d.onSaveStaticMapReq)
}

func (d *Database) onCreatePlayerReq(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.CreatePlayerResp{}
	defer node.Instance.Net.Responce(header, resp)
	message := &msg.CreatePlayerReq{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		logger.Debug(err)
		return
	}
	player := &data.Player{}
	serializedData, err := proto.Marshal(player)
	if err != nil {
		logger.Debug(err)
		return
	}
	dbplayer := &db.DBPlayer{
		UUID:           message.PlayerUUID,
		SerializedData: serializedData,
	}
	err = dao.PlayerInsert(d.DB, dbplayer)
	if err != nil {
		logger.Debug(err)
		return
	}
	resp.Success = true
}

func (d *Database) onLoadPlayerReq(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.LoadPlayerResp{}
	defer node.Instance.Net.Responce(header, resp)
	message := &msg.LoadPlayerReq{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		logger.Debug(err)
		return
	}
	player, ok := d.players[message.PlayerUUID]
	if ok {
		resp.Success = true
		resp.Player = player
	} else {
		dbplayer, err := dao.PlayerQuery(d.DB, message.PlayerUUID)
		if err != nil {
			logger.Debug(err)
			return
		}
		player := &data.Player{}
		err = proto.Unmarshal(dbplayer.SerializedData, player)
		if err != nil {
			logger.Debug(err)
			return
		}
		resp.Success = true
		resp.Player = player
		d.players[message.PlayerUUID] = player
	}
}
func (d *Database) onReleasePlayerNtf(header *msg.MessageHeader, buffer []byte) {
	message := &msg.ReleasePlayerNtf{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		logger.Debug(err)
		return
	}
	if _, ok := d.players[message.PlayerUUID]; ok {
		delete(d.players, message.PlayerUUID)
	}
}

func (d *Database) onCreateRoleReq(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.CreateRoleResp{}
	defer node.Instance.Net.Responce(header, resp)
	message := &msg.CreateRoleReq{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		logger.Debug(err)
		return
	}
	if player, ok := d.players[message.PlayerUUID]; ok {
		if _, ok := d.roleSummaryByName[message.RoleName]; ok {
			return
		}

		roleUUID := uuid.New()
		roleSummaryData := &data.RoleSummaryData{
			Name:     message.RoleName,
			Zone:     message.Zone,
			RoleUUID: roleUUID,
		}
		d.roleSummaryByName[message.RoleName] = roleSummaryData

		getStaticMapUUIDReq := &msg.GetStaticMapUUIDReq{
			ZoneID:      message.Zone,
			StaticMapID: 1,
		}
		getStaticMapUUIDResp := &msg.GetStaticMapUUIDResp{}
		err := node.Instance.Net.Request(msg.Command_GET_STATIC_MAP_UUID_REQ, getStaticMapUUIDReq, getStaticMapUUIDResp)
		if err != nil {
			logger.Error(err)
			return
		}
		if len(getStaticMapUUIDResp.StaticMapUUID) == 0 {
			logger.Error("没有找到这张地图")
			return
		}
		roleSummaryData.MapUUID = getStaticMapUUIDResp.StaticMapUUID
		summaryData, err := proto.Marshal(roleSummaryData)
		if err != nil {
			return
		}
		roleOnlineData := &data.RoleOnlineData{}
		onlineData, err := proto.Marshal(roleOnlineData)
		if err != nil {
			return
		}
		dbrole := &db.DBRole{
			UUID:        roleUUID,
			SummaryData: summaryData,
			OnlineData:  onlineData,
		}
		err = dao.RoleInsert(d.DB, dbrole)
		if err != nil {
			logger.Debug(err)
			return
		}
		player.RoleList = append(player.RoleList, roleSummaryData)
		serializedData, err := proto.Marshal(player)
		if err != nil {
			logger.Error(err)
			return
		}
		dbplayer := &db.DBPlayer{
			UUID:           message.PlayerUUID,
			SerializedData: serializedData,
		}
		err = dao.PlayerUpdate(d.DB, dbplayer)
		if err != nil {
			logger.Error(err)
			return
		}
		d.roleSummary[roleUUID] = roleSummaryData
		resp.Success = true
		resp.Role = roleSummaryData
	} else {
		logger.Debug("Must Create Player Before Create Role!")
	}
}

func (d *Database) onLoadRoleReq(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.LoadRoleResp{}
	defer node.Instance.Net.Responce(header, resp)
	message := &msg.LoadRoleReq{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		logger.Error(err)
		return
	}

	if roleSummary, ok := d.roleSummary[message.RoleUUID]; ok {

		mapAddressReq := &msg.GetMapAddressReq{MapUUID: roleSummary.MapUUID}
		mapAddressResp := &msg.GetMapAddressResp{}
		err := node.Instance.Net.Request(msg.Command_GET_MAP_ADDRESS_REQ, mapAddressReq, mapAddressResp)
		if err != nil {
			logger.Error(err)
			return
		}
		if mapAddressResp.ServerID == global.InvalidServerID {
			logger.Error("玩家所在地图没有创建 MapUUID:"+roleSummary.MapUUID)
			return
		}

		onlineData, err := dao.RoleOnlineDataQuery(d.DB, message.RoleUUID)
		if err != nil {
			logger.Error(err)
			return
		}
		roleOnline := &data.RoleOnlineData{}
		err = proto.Unmarshal(onlineData, roleOnline)
		if err != nil {
			logger.Error(err)
			return
		}

		role := &data.Role{
			SummaryData: roleSummary,
			OnlineData:  roleOnline,
		}

		resp.Success = true

		node.Instance.Net.Notify(msg.Command_ROLE_ENTER, role)
	}
}

func (d *Database) onLoadStaticMapReq(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.LoadStaticMapResp{}
	defer node.Instance.Net.Responce(header, resp)
	message := &msg.LoadStaticMapReq{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		logger.Error(err)
		return
	}
	if maps, ok := d.staticMaps[message.ZoneID]; ok {
		if staticMap, ok2 := maps[message.StaticMapID]; ok2 {
			resp.Map = staticMap
		}
	}
}

func (d *Database) onSaveStaticMapReq(header *msg.MessageHeader, buffer []byte) {
	resp := &msg.SaveStaticMapResp{}
	defer node.Instance.Net.Responce(header, resp)
	req := &msg.SaveStaticMapReq{}
	err := proto.Unmarshal(buffer, req)
	if err != nil {
		logger.Error(err)
		return
	}
	staticMaps := make([]*db.DBStaticMap, 0)
	logger.Debug(fmt.Sprintf("onSaveStaticMapReq len:%d", len(req.StaticMaps)))
	for _, staticMap := range req.StaticMaps {
		serializedData, err := proto.Marshal(staticMap)
		if err != nil {
			logger.Error(err)
		} else {
			dbStaticMap := &db.DBStaticMap{}
			dbStaticMap.UUID = staticMap.MapUUID
			dbStaticMap.SerializedData = serializedData
			staticMaps = append(staticMaps, dbStaticMap)
		}
	}

	err = dao.BulkStaticMapUpdate(d.DB, staticMaps)
	if err == nil {
		resp.Success = true
	}
}

func (d *Database) loadAllRoleSummaryData() {
	roles := dao.AllRoleSummaryDataQuery(d.DB)
	for _, role := range roles {
		summary := &data.RoleSummaryData{}
		proto.Unmarshal(role.SummaryData, summary)
		d.roleSummary[role.UUID] = summary
	}

	for _, role := range d.roleSummary {
		d.roleSummaryByName[role.Name] = role
	}

	logger.Info("加载所有角色的摘要数据成功")
}

func (d *Database) loadAllStaticMapsData() {
	staticMaps := dao.AllStaticMapQuery(d.DB)
	for _, staticMap := range staticMaps {
		mapdata := &data.MapData{}
		proto.Unmarshal(staticMap.SerializedData, mapdata)
		if _, ok := d.staticMaps[staticMap.ZoneID]; ok == false {
			d.staticMaps[staticMap.ZoneID] = make(map[int32]*data.MapData)
		}
		d.staticMaps[staticMap.ZoneID][staticMap.MapID] = mapdata
	}

	logger.Info("加载所有静态地图数据成功")
}
