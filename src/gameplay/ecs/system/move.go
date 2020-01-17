package system

import "INServer/src/gameplay/ecs"

import "INServer/src/proto/data"

import "INServer/src/proto/engine"

type (
	move struct {
	}
)

func (m *move) Tick(dt float32, entities []*ecs.Entity) {
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

func step(dt float32, pos *engine.Vector3, rspeed *engine.Vector3, pspeed *engine.Vector3) {
	pos.X += dt * (rspeed.X + pspeed.X)
	pos.Y += dt * (rspeed.Y + pspeed.Y)
	pos.Z += dt * (rspeed.Z + pspeed.Z)
}
