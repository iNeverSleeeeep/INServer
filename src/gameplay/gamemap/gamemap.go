package gamemap

import (
	"INServer/src/gameplay/ecs"
	"INServer/src/gameplay/ecs/system"
	"time"
	"INServer/src/proto/config"
)

type (
	Map struct {
		mapData    *data.MapData
		scenes     []*Scene

		firstScene *Scene
		entitiesMap   map[string]*ecs.Entity
		running bool
	}
)

func NewMap(mapConfig *config.Map) *Map {
	m := new(Map)
	m.scenes = make([]*Scene, 0)
	m.entitiesLocation = make(map[string]*Scene, 0)
	m.firstScene = NewScene(m, nil)
	return m
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
			dt := float64(now - lasttime) / float64(1E9)
			
			system.Tick(float32(dt), m.mapData.entities)
			
			time.Sleep(time.Millisecond * 33)
		}
	}
}

func (m *Map) EntityEnter(uuid string, entity *ecs.Entity) {
	m.entities[uuid] = entity
	m.firstScene.EntityEnter(uuid, entity)
}

func (m *Map) EntityLeave(uuid string) {
	if entity, ok := m.entitiesMap[uuid]; ok {
		m.firstScene.EntityLeave(uuid, entity)
		delete(m.entities, uuid)
	}
}

func (m *Map) Tick() {

}

func (m *Map) SyncEntityPosition(uuid string) {
	if entity, ok := m.entitiesMap[uuid]; ok {
		m.firstScene.SyncEntityPosition(uuid, entity)
	}
}
