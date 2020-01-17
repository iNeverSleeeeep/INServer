package system

import "INServer/src/gameplay/ecs"

var (
	MoveSystem = &Move{}
)

type (
	ISystem interface {
		Tick(entities []*ecs.Entity)
	}
)
