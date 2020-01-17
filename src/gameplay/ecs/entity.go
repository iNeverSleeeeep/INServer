package ecs

import "INServer/src/proto/data"

type (
	Entity struct {
		components []*data.Component
	}
)

func NewEntity(entityData *data.EntityData, entityType *data.EntityType) *Entity {
	e := new(Entity)
	e.components = make([]*data.Component, len(data.ComponentType_value))
	return e
}

func (e *Entity) AddComponent(component *data.Component) {
	e.components[component.Type] = component
}

func (e *Entity) RemoveComponent(componentType data.ComponentType) {
	e.components[componentType] = nil
}

func (e *Entity) GetComponent(componentType data.ComponentType) *data.Component {
	if e.components[componentType] == nil {
		return &data.Component{}
	} else {
		return e.components[componentType]
	}
}
