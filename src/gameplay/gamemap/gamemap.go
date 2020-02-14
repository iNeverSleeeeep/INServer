package gamemap

import (
	"INServer/src/gameplay/ecs"
	"INServer/src/gameplay/ecs/system"
	"INServer/src/proto/config"
	"INServer/src/proto/data"
	"INServer/src/proto/engine"
	"INServer/src/proto/msg"
	"time"
)

type (
	ComponentValueCache struct {
		position engine.Vector3
	}

	Map struct {
		mapData   *data.MapData
		mapConfig *config.Map
		scenes    []*Scene

		firstScene  *Scene
		entitiesMap map[string]*ecs.Entity
		running     bool
	}
)

func NewMap(mapConfig *config.Map, mapData *data.MapData) *Map {
	m := new(Map)
	m.scenes = make([]*Scene, 0)
	m.firstScene = NewScene(m, nil)
	m.scenes = append(m.scenes, m.firstScene)
	m.mapData = mapData
	m.mapConfig = mapConfig
	m.entitiesMap = make(map[string]*ecs.Entity)
	if m.mapData == nil {
		m.mapData = &data.MapData{
			MapID:    mapConfig.MapID,
			MapUUID:  mapData.MapUUID,
			Entities: make([]*data.EntityData, 0),
		}
	}
	return m
}

func (m *Map) MapData() *data.MapData {
	return m.mapData
}

func (m *Map) MapConfig() *config.Map {
	return m.mapConfig
}

func (m *Map) Scenes() []*Scene {
	return m.scenes
}

func (m *Map) GetEntity(uuid string) *ecs.Entity {
	if entity, ok := m.entitiesMap[uuid]; ok {
		return entity
	}
	return nil
}

func (m *Map) Start() {
	m.running = true
	if m.mapData.LastTickTime == 0 {
		m.mapData.LastTickTime = time.Now().UnixNano()
	}
	go func() {
		for m.running {
			lasttime := m.mapData.LastTickTime
			now := time.Now().UnixNano()
			dt := float64(now-lasttime) / float64(1E9)

			m.tickSystems(dt)

			time.Sleep(time.Millisecond * 33)
		}
	}()
}

func (m *Map) EntityEnter(uuid string, entity *ecs.Entity) {
	m.mapData.Entities = append(m.mapData.Entities, entity.EntityData())
	m.entitiesMap[uuid] = entity
	m.firstScene.EntityEnter(uuid, entity)
}

func (m *Map) EntityLeave(uuid string) {
	if entity, ok := m.entitiesMap[uuid]; ok {
		m.firstScene.EntityLeave(uuid, entity)
		delete(m.entitiesMap, uuid)
		for index, entityData := range m.mapData.Entities {
			if entityData.EntityUUID == uuid {
				m.mapData.Entities = append(m.mapData.Entities[:index], m.mapData.Entities[index+1:]...)
				break
			}
		}
	}
}

func (m *Map) Tick() {
	for _, scene := range m.scenes {
		scene.Tick()
	}
}

func (m *Map) tickSystems(dt float64) {
	cachedValues := make(map[string]*ComponentValueCache)
	for uuid, entity := range m.entitiesMap {
		cache := &ComponentValueCache{}
		cachedValues[uuid] = cache
		transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
		if transform != nil {
			cache.position = *transform.Position
		}
	}
	system.Tick(dt, m.entitiesMap)
	for uuid, entity := range m.entitiesMap {
		transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
		if transform != nil {
			m.firstScene.SyncEntityPosition(uuid, entity, &cachedValues[uuid].position)
		}
	}
}

// OnRoleMoveINF 响应角色移动
func (m *Map) OnRoleMoveINF(role *data.Role, inf *msg.MoveINF) {
	entity := m.GetEntity(role.SummaryData.RoleUUID)
	if entity != nil {
		m.firstScene.onEntityMoveINF(entity, inf)
	}
}

// OnRoleStopMoveINF 响应角色停止移动
func (m *Map) OnRoleStopMoveINF(role *data.Role, inf *msg.StopMoveINF) {
	entity := m.GetEntity(role.SummaryData.RoleUUID)
	if entity != nil {
		m.firstScene.onEntityStopMoveINF(entity, inf)
	}
}
