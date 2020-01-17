package gamemap

import (
	"INServer/src/gameplay/ecs"
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
			dt := float32(now - lasttime) / float32(1E6)
			

			time.Sleep(time.Millisecond * 33)
		}
	}
}

func (m *Map) EntityEnter(uuid string, entity *ecs.Entity) {
	m.entities[uuid] = entity
	m.firstScene.EntityEnter(uuid, entity)
}

func (m *Map) EntityLeave(uuid string, entity *ecs.Entity) {
	m.firstScene.EntityLeave(uuid, entity)
	delete(m.entities, uuid)
}

func (m *Map) Tick() {

}
