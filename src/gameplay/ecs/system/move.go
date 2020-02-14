package system

import (
	"INServer/src/engine/extensions/vector3"
	"INServer/src/gameplay/ecs"
	"INServer/src/proto/data"
	"INServer/src/proto/engine"
)

type move struct {
}

func (m *move) Tick(dt float64, entities map[string]*ecs.Entity) {
	for _, entity := range entities {
		physics := entity.GetComponent(data.ComponentType_Physics).Physics
		if physics != nil {
			transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
			if transform != nil {
				stop(dt, entity, physics, transform)
				step(dt, transform.Position, vector3.Add(physics.RawSpeed, physics.PassiveSpeed))
			}
		}
	}
}

func stop(dt float64, entity *ecs.Entity, physics *data.PhysicsComponent, transform *data.TransformComponent) {
	move := entity.GetComponent(data.ComponentType_Move).Move
	if move != nil {
		dir := vector3.Minus(move.Destination, transform.Position)
		// 如果前进方向和目标方向相反 则停止移动
		if vector3.Dot(physics.RawSpeed, dir) < 0 {
			physics.RawSpeed = &engine.Vector3{}
			transform.Position = move.Destination
		}
	} else {
		physics.RawSpeed = &engine.Vector3{}
	}
}

func step(dt float64, pos *engine.Vector3, speed *engine.Vector3) bool {
	X, Y, Z := pos.X, pos.Y, pos.Z
	pos.X += dt * speed.X
	pos.Y += dt * speed.Y
	pos.Z += dt * speed.Z
	return X != pos.X || Y != pos.Y || Z != pos.Z
}
