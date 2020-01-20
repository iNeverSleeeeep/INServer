package system

import "INServer/src/gameplay/ecs"

var (
	MoveSystem    = &move{}
	PhysicsSystem = &physics{}
)

type (
	ISystem interface {
		Tick(dt float64, entities map[string]*ecs.Entity)
	}
)

func Tick(dt float64, entities map[string]*ecs.Entity) {
	MoveSystem.Tick(dt, entities)
	PhysicsSystem.Tick(dt, entities)
}
