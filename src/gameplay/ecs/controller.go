package ecs

import (
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"INServer/src/services/node"

	"github.com/gogo/protobuf/proto"
)

type (
	// Controller 全部Controller的接口
	Controller interface {
		OnOtherMove(entity *Entity, ntf *msg.MoveNTF)
		OnNearEntities([]*msg.NearEntity)
	}

	// DummyController 占位 什么事情也不做
	DummyController struct {
		entity *Entity
	}
	// AIController 服务器控制
	AIController struct {
		entity *Entity
	}
	// PlayerController 玩家控制
	PlayerController struct {
		entity *Entity
	}
)

func initController(entity *Entity) {
	switch entity.entityType {
	case data.EntityType_MonsterEntity, data.EntityType_NPCEntity:
		entity.controller = &AIController{entity}
	case data.EntityType_RoleEntity:
		entity.controller = &PlayerController{entity}
	}
	if entity.controller == nil {
		entity.controller = &DummyController{entity}
	}
}

func (c *DummyController) OnOtherMove(entity *Entity, ntf *msg.MoveNTF) {

}
func (c *DummyController) OnNearEntities([]*msg.NearEntity) {

}
func (c *AIController) OnOtherMove(entity *Entity, ntf *msg.MoveNTF) {

}
func (c *AIController) OnNearEntities([]*msg.NearEntity) {

}
func (c *PlayerController) OnOtherMove(entity *Entity, ntf *msg.MoveNTF) {
	c.sendMessage(msg.CMD_MOVE_NTF, ntf)
}
func (c *PlayerController) OnNearEntities(items []*msg.NearEntity) {
	if len(items) > 1 {
		nearEntitiesNTF := &msg.NearEntitiesNTF{
			Entities: items,
		}
		node.Net.NotifyClient(msg.CMD_NEAR_ENTITIES_NTF, nearEntitiesNTF, c.entity.UUID())
	}
}

func (c *PlayerController) sendMessage(command msg.CMD, message proto.Message) {
	//gate := world.Instance.RoleGate(c.entity.UUID())
	//if gate != global.InvalidServerID {
	//	if buffer, err := proto.Marshal(message); err == nil {
	//		forward := &msg.ForwardPlayerMessage{}
	//	}
	//	node.Net.NotifyServer(msg.CMD_MOVE_NTF, message, gate)
	//}
}
