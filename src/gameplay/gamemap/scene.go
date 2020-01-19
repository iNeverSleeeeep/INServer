package gamemap

import (
	"INServer/src/engine/grid"
	"INServer/src/gameplay/ecs"
	"INServer/src/proto/config"
	"INServer/src/proto/data"
	"INServer/src/proto/engine"
)

type (
	Scene struct {
		masterMap *Map
		search    *grid.Grid
	}
)

func NewScene(masterMap *Map, sceneConfig *config.Scene) *Scene {
	s := new(Scene)
	s.masterMap = masterMap
	s.search = grid.New(10, 1000, 1000)
	return s
}

func (s *Scene) EntityEnter(uuid string, entity *ecs.Entity) {
	transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
	if transform != nil {
		s.search.Add(uuid, transform.Position)
	}
}

func (s *Scene) EntityLeave(uuid string, entity *ecs.Entity) {
	transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
	if transform != nil {
		s.search.Remove(uuid, transform.Position)
	}
}

func (s *Scene) SyncEntityPosition(uuid string, entity *ecs.Entity, from *engine.Vector3) {
	transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
	if transform != nil {
		s.search.Move(uuid, from, transform.Position)
	}
}
