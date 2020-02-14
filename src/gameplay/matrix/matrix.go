package matrix

import (
	"INServer/src/common/uuid"
	"INServer/src/gameplay/ecs"
	"INServer/src/gameplay/gamemap"
	"INServer/src/proto/data"
)

// Matrix 控制地图内实体的产生消亡
type Matrix struct {
}

// New 构造Matrix
func New() *Matrix {
	m := new(Matrix)
	return m
}

// OnGamemapCreate 地图创建时的初始化
func (m *Matrix) OnGamemapCreate(gameMap *gamemap.Map) {
	if len(gameMap.MapData().Entities) == 0 {
		for _, scene := range gameMap.Scenes() {
			for x := scene.Width / -2; x < scene.Width/2; x++ {
				for z := scene.Height / -2; z < scene.Height/2; z++ {
					entityUUID := uuid.New()
					entity := ecs.NewEntity(&data.EntityData{EntityUUID: entityUUID}, data.EntityType_MonsterEntity)
					if entity != nil {
						gameMap.EntityEnter(entity.UUID(), entity)
					}
				}
			}
		}
	}
}

// OnGamemapDestroy 地图销毁时的回收
func (m *Matrix) OnGamemapDestroy(gameMap *gamemap.Map) {

}

// TickGamemap Tick
func (m *Matrix) TickGamemap(gameMap *gamemap.Map) {

}
