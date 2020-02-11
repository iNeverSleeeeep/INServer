package gamemap

import (
	"INServer/src/common/logger"
	"INServer/src/engine/extensions/vector3"
	"INServer/src/engine/grid"
	"INServer/src/gameplay/ecs"
	"INServer/src/proto/config"
	"INServer/src/proto/data"
	"INServer/src/proto/engine"
	"INServer/src/proto/msg"
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

func (s *Scene) onEntityMoveINF(entity *ecs.Entity, inf *msg.MoveINF) {
	transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
	if transform != nil {
		s.search.Move(entity.UUID(), transform.Position, inf.Position)
		transform.Position = inf.Position
		physics := entity.GetComponent(data.ComponentType_Physics).Physics
		attribute := entity.GetComponent(data.ComponentType_Attribute).Attribute
		move := entity.GetComponent(data.ComponentType_Move).Move
		if physics != nil && attribute != nil && move != nil {
			move.Destination = inf.To
			physics.RawSpeed = vector3.Multiply(vector3.Normalize(vector3.Minus(inf.To, inf.Position)), float64(attribute.Speed))
			ntf := &msg.MoveNTF{}
			ntf.EntityUUID = entity.UUID()
			ntf.To = inf.To
			items := s.search.GetNearItems(inf.Position)
			for _, item := range items {
				logger.Info(item.UUID())
				//nearEntity := s.masterMap.GetEntity(item.UUID())
			}
		}
	}
}
