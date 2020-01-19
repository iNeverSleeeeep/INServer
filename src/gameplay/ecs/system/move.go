package system

import "INServer/src/gameplay/ecs"

import "INServer/src/proto/data"

import "INServer/src/proto/engine"

import "INServer/src/modules/world"

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
				position := *transform.Position
				dirty := step(dt, transform.Position, physics.RawSpeed, physics.PassiveSpeed)
				if dirty {
					gamemap := world.Instance.GetMap(entity.RealTimeData().GetCurrentMapUUID())
					gamemap.SyncEntityPosition(entity.UUID(), &position)
				}
			}
		}
	}
}

func step(dt float32, pos *engine.Vector3, rspeed *engine.Vector3, pspeed *engine.Vector3) bool {
	X, Y, Z := pos.X, pos.Y, pos.Z
	pos.X += dt * (rspeed.X + pspeed.X)
	pos.Y += dt * (rspeed.Y + pspeed.Y)
	pos.Z += dt * (rspeed.Z + pspeed.Z)
	return X != pos.X || Y != pos.Y || Z != pos.Z
}
