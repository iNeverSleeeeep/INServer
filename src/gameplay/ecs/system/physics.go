package system

import (
	"INServer/src/gameplay/ecs"
	"INServer/src/proto/data"
)

type physics struct {
}

func (m *physics) Tick(dt float64, entities map[string]*ecs.Entity) {
	for _, entity := range entities {
		physics := entity.GetComponent(data.ComponentType_Physics).Physics
		if physics != nil {
			slowdown(dt, physics)
		}
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func slowdown(dt float64, physics *data.PhysicsComponent) {
	a := 1 / physics.Mass
	physics.PassiveSpeed.X -= max(a*dt, physics.PassiveSpeed.X)
	physics.PassiveSpeed.Y -= max(a*dt, physics.PassiveSpeed.Y)
	physics.PassiveSpeed.Z -= max(a*dt, physics.PassiveSpeed.Z)
}
