package world

import (
	"INServer/src/common/global"
	"INServer/src/gameplay/gamemap"
	"INServer/src/modules/node"
	"INServer/src/proto/msg"
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
	for _, zoneConfig := range global.ServerConfig.WorldConfig.Zones {
		for _, gameMap := range zoneConfig.StaticMaps {
			node.Instance.Net.Request(msg.Command_LOAD_MAP_REQ)
			gamemap.NewMap()
		}
	}
}

func (w *World) GetMap(uuid string) *gamemap.Map {
	if result, ok := w.gameMaps[uuid]; ok {
		return result
	}
	return nil
}
