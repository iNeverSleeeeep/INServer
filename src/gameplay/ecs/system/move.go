package system

import (
	"INServer/src/gameplay/ecs"
	"INServer/src/proto/data"
	"INServer/src/proto/engine"
)

type (
	move struct {
	}
)

func (m *move) Tick(dt float64, entities map[string]*ecs.Entity) {
	for _, entity := range entities {
		physics := entity.GetComponent(data.ComponentType_Physics).Physics
		if physics != nil {
			transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
			if transform != nil {
				step(dt, transform.Position, physics.RawSpeed, physics.PassiveSpeed)
			}
		}
	}
}

func step(dt float64, pos *engine.Vector3, rspeed *engine.Vector3, pspeed *engine.Vector3) bool {
	X, Y, Z := pos.X, pos.Y, pos.Z
	pos.X += dt * (rspeed.X + pspeed.X)
	pos.Y += dt * (rspeed.Y + pspeed.Y)
	pos.Z += dt * (rspeed.Z + pspeed.Z)
	return X != pos.X || Y != pos.Y || Z != pos.Z
}
