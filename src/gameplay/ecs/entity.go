package ecs

import "INServer/src/proto/data"

type (
	Entity struct {
		entityData *data.EntityData
	}
)

func NewEntity(entityData *data.EntityData, entityType *data.EntityType) *Entity {
	e := new(Entity)
	return e
}

func (e *Entity) AddComponent(component *data.Component) {
	e.entityData.Components[component.Type] = component
}

func (e *Entity) RemoveComponent(componentType data.ComponentType) {
	e.entityData.Components[componentType] = nil
}

func (e *Entity) GetComponent(componentType data.ComponentType) *data.Component {
	if e.entityData.Components[componentType] == nil {
		return &data.Component{}
	} else {
		return e.entityData.Components[componentType]
	}
}

func (e *Entity) RealTimeData() *data.EntityRealtimeData {
	return e.entityData.RealTimeData
}

func (e *Entity) UUID() string {
	return e.entityData.EntityUUID
}
