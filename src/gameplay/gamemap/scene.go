package gamemap

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/engine/extensions/vector3"
	"INServer/src/engine/grid"
	"INServer/src/gameplay/ecs"
	"INServer/src/proto/config"
	"INServer/src/proto/data"
	"INServer/src/proto/engine"
	"INServer/src/proto/msg"
	"INServer/src/services/node"
	"time"
)

type (
	Scene struct {
		masterMap *Map
		search    *grid.Grid
		Width     int32
		Height    int32
		entities  map[string]*ecs.Entity

		syncTime int64
	}
)

func NewScene(masterMap *Map, sceneConfig *config.Scene) *Scene {
	s := new(Scene)
	s.masterMap = masterMap
	s.Width = 1000
	s.Height = 1000
	s.search = grid.New(10, s.Width, s.Height)
	s.entities = make(map[string]*ecs.Entity)
	return s
}

func (s *Scene) Tick() {
	now := time.Now().UnixNano()
	if s.syncTime+global.NANO_PER_SECONE < now {
		s.syncEntitiesNTF()
	}
}

func (s *Scene) EntityEnter(uuid string, entity *ecs.Entity) {
	transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
	if transform != nil {
		s.search.Add(uuid, transform.Position)
	}
	s.entities[uuid] = entity
}

func (s *Scene) EntityLeave(uuid string, entity *ecs.Entity) {
	transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
	if transform != nil {
		s.search.Remove(uuid, transform.Position)
	}
	delete(s.entities, uuid)
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
				logger.Info(item.EntityUUID)
				//nearEntity := s.masterMap.GetEntity(item.UUID())
			}
		}
	}
}

func (s *Scene) onEntityStopMoveINF(entity *ecs.Entity, inf *msg.StopMoveINF) {
	transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
	if transform != nil {
		s.search.Move(entity.UUID(), transform.Position, inf.Position)
		transform.Position = inf.Position
		physics := entity.GetComponent(data.ComponentType_Physics).Physics
		if physics != nil {
			physics.RawSpeed = &engine.Vector3{}
		}
		move := entity.GetComponent(data.ComponentType_Move).Move
		if move != nil {
			move.Destination = inf.Position
		}
	}
}

func (s *Scene) syncEntitiesNTF() {
	for _, entity := range s.entities {
		transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
		if transform != nil {
			items := s.search.GetNearItems(transform.Position)
			if len(items) > 1 {
				nearEntitiesNTF := &msg.NearEntitiesNTF{
					Entities: items,
				}
				gate := global.RoleGateGetter.GetRoleGate(entity.UUID())
				node.Net.NotifyServer(msg.CMD_NEAR_ENTITIES_NTF, nearEntitiesNTF, gate)
			}
		}
	}
}
