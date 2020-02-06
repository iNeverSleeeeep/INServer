package ecs

import (
	"INServer/src/proto/data"
	"INServer/src/proto/engine"
)

type (
	// Entity 游戏实体
	Entity struct {
		entityData *data.EntityData
		entityType data.EntityType
		controller Controller
	}
)

// NewEntity 构造实体
func NewEntity(entityData *data.EntityData, entityType data.EntityType) *Entity {
	e := new(Entity)
	e.entityData = entityData
	e.entityType = entityType
	initController(e)
	return e
}

// AddComponent 添加组件
func (e *Entity) AddComponent(component *data.Component) {
	e.entityData.Components[component.Type] = component
}

// RemoveComponent 移除组件
func (e *Entity) RemoveComponent(componentType data.ComponentType) {
	e.entityData.Components[componentType] = nil
}

// GetComponent 取得组件
func (e *Entity) GetComponent(componentType data.ComponentType) *data.Component {
	if e.entityData.Components[componentType] == nil {
		return &data.Component{}
	}
	return e.entityData.Components[componentType]
}

// RealTimeData 实时数据
func (e *Entity) RealTimeData() *data.EntityRealtimeData {
	return e.entityData.RealTimeData
}

// EntityData 实体数据
func (e *Entity) EntityData() *data.EntityData {
	return e.entityData
}

// UUID 返回UUID
func (e *Entity) UUID() string {
	return e.entityData.EntityUUID
}

// InitComponents 根据实体类型初始化组件
func InitComponents(entityType data.EntityType) []*data.Component {
	components := make([]*data.Component, len(data.ComponentType_name))
	switch entityType {
	default:
		components[data.ComponentType_Invalid] = &data.Component{
			Type: data.ComponentType_Invalid,
		}
		components[data.ComponentType_Transofrm] = &data.Component{
			Type: data.ComponentType_Transofrm,
			Transform: &data.TransformComponent{
				Position: &engine.Vector3{},
				Rotation: &engine.Quaternion{},
			},
		}
		components[data.ComponentType_Physics] = &data.Component{
			Type: data.ComponentType_Physics,
			Physics: &data.PhysicsComponent{
				Mass:         100,
				RawSpeed:     &engine.Vector3{},
				PassiveSpeed: &engine.Vector3{},
			},
		}
		components[data.ComponentType_Attribute] = &data.Component{
			Type: data.ComponentType_Attribute,
			Attribute: &data.AttributeComponent{
				Speed: 10,
				HP:    100,
			},
		}
		break
	}
	return components
}
