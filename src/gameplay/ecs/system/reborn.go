package system

import (
	"INServer/src/gameplay/ecs"
	"INServer/src/proto/data"
	"time"
)

type reborn struct {
}

func (r *reborn) Tick(dt float64, entities map[string]*ecs.Entity) {
	now := time.Now().UnixNano()
	for _, entity := range entities {
		reborn := entity.GetComponent(data.ComponentType_Reborn).Reborn
		if reborn != nil && reborn.RebornType == data.RebornType_Auto && reborn.RebornTime > 0 {
			if reborn.RebornTime < now {
				reset(entity, reborn)
			}
		}
	}
}

func reset(entity *ecs.Entity, reborn *data.RebornComponent) {
	attribute := entity.GetComponent(data.ComponentType_Attribute).Attribute
	if attribute != nil {
		attribute.HP = attribute.MaxHP
	}
	transform := entity.GetComponent(data.ComponentType_Transofrm).Transform
	if transform != nil {
		transform.Position = reborn.Position
	}
}
