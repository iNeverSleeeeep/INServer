package system

import "INServer/src/gameplay/ecs"

var (
	MoveSystem    = &move{}
	PhysicsSystem = &physics{}
)

type (
	ISystem interface {
		Tick(dt float32, entities []*ecs.Entity)
	}
)

func Tick(dt float32, entities []*ecs.Entity) {
	MoveSystem.Tick(dt, entities)
	PhysicsSystem.Tick(dt, entities)
}
