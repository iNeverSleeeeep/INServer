package matrix

import (
	"INServer/src/common/uuid"
	"INServer/src/gameplay/ecs"
	"INServer/src/gameplay/gamemap"
	"INServer/src/proto/data"
	"INServer/src/proto/engine"
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
			for x := scene.Width / -2; x < scene.Width/2; x = x + 10 {
				for z := scene.Height / -2; z < scene.Height/2; z = z + 10 {
					entityUUID := uuid.New()
					components := ecs.InitComponents(data.EntityType_MonsterEntity)
					components[data.ComponentType_Transofrm].Transform.Position = &engine.Vector3{
						X: float64(x),
						Y: 0,
						Z: float64(z),
					}
					entityData := &data.EntityData{
						EntityUUID: entityUUID,
						Components: components,
						RealTimeData: &data.EntityRealtimeData{
							LastStaticMapUUID: gameMap.MapData().MapUUID,
							CurrentMapUUID:    gameMap.MapData().MapUUID,
						},
					}
					entity := ecs.NewEntity(entityData, data.EntityType_MonsterEntity)
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
