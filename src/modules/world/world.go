package world

import (
	"INServer/src/gameplay/gamemap"
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
	//for _, zoneConfig := range global.ServerConfig.WorldConfig.Zones {

	//}
}

func (w *World) GetMap(uuid string) *gamemap.Map {
	if result, ok := w.gameMaps[uuid]; ok {
		return result
	}
	return nil
}
