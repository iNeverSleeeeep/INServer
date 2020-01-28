package world

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/uuid"
	"INServer/src/gameplay/gamemap"
	"INServer/src/modules/node"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"time"
)

var Instance *World

type (
	World struct {
		gameMaps map[string]*gamemap.Map
	}
)

func New() *World {
	w := new(World)
	w.gameMaps = make(map[string]*gamemap.Map)
	return w
}

func (w *World) Start() {
	req := &msg.LoadStaticMapReq{}
	resp := &msg.LoadStaticMapResp{}
	for _, zoneConfig := range global.ServerConfig.WorldConfig.Zones {
		for _, gameMapID := range zoneConfig.StaticMaps {
			req.ZoneID = zoneConfig.ZoneID
			req.StaticMapID = gameMapID
			resp.Reset()
			var mapData *data.MapData
			err := node.Instance.Net.Request(msg.Command_LOAD_STATIC_MAP_REQ, req, resp)
			if err != nil {
				logger.Error(err)
			} else if resp.Map == nil {
				mapData = &data.MapData{}
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
					ServerID: global.ServerID,
				}
				node.Instance.Net.Notify(msg.Command_UPDATE_MAP_ADDRESS_NTF, updateMapAddress)
				udpateStatcMapUUID := &msg.UpdateStaticMapUUIDNTF{
					ZoneID:        zoneConfig.ZoneID,
					StaticMapID:   gameMapID,
					StaticMapUUID: mapData.MapUUID,
				}
				node.Instance.Net.Notify(msg.Command_UPDATE_STATIC_MAP_UUID_NTF, udpateStatcMapUUID)
			}
		}
	}
}

func (w *World) Stop() {
	staticMaps := make([]*data.MapData, 0)
	req := &msg.SaveStaticMapReq{
		StaticMaps: staticMaps,
	}
	for _, gamemap := range w.gameMaps {
		staticMaps = append(staticMaps, gamemap.MapData())
	}
	resp := &msg.SaveStaticMapResp{}
	err := node.Instance.Net.Request(msg.Command_SAVE_STATIC_MAP_REQ, req, resp)
	if err != nil {
		logger.Error(err)
	}
}

func (w *World) GetMap(uuid string) *gamemap.Map {
	if result, ok := w.gameMaps[uuid]; ok {
		return result
	}
	return nil
}
