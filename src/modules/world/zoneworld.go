package world

import "INServer/src/gameplay/gamemap"

type (
	ZoneWorld struct {
		gameMaps map[int32]*gamemap.Map
	}
)

func NewZoneWorld() *ZoneWorld {
	zw := new(ZoneWorld)
	zw.gameMaps = make(map[int32]*gamemap.Map)
	return zw
}