package gamemap

import (
	"INServer/src/gameplay/ecs"
	"INServer/src/proto/config"
)

type (
	Map struct {
		scenes     []*Scene
		firstScene *Scene
		entities   map[string]*ecs.Entity
	}
)

func NewMap(mapConfig *config.Map) *Map {
	m := new(Map)
	m.scenes = make([]*Scene, 0)
	m.firstScene = NewScene(m, nil)
	return m
}

func (m *Map) EntityEnter(uuid string, entity *ecs.Entity) {
	m.entities[uuid] = entity
	m.firstScene.EntityEnter(uuid, entity)
}

func (m *Map) EntityLeave(uuid string, entity *ecs.Entity) {
	m.firstScene.EntityLeave(uuid, entity)
	delete(m.entities, uuid)
}
