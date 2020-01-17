package system

import "INServer/src/gameplay/ecs"

var (
	MoveSystem = &Move{}
)

type (
	ISystem interface {
		Tick(dt float32, entities []*ecs.Entity)
	}
)
