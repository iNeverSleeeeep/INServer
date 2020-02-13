package balcony

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"INServer/src/services/node"
	"fmt"

	"github.com/gogo/protobuf/proto"
)

// Instance 月台单例
var Instance *Balcony

type (
	// Balcony 月台 负责角色进入游戏世界之前的一些逻辑
	Balcony struct {
	}
)

// New 构造月台实例
func New() *Balcony {
	b := new(Balcony)
	b.initMessageHandler()
	return b
}

// Start 服务启动
func (b *Balcony) Start() {

}

// Stop 服务停止
func (b *Balcony) Stop() {

}

func (b *Balcony) initMessageHandler() {
	node.Net.Listen(msg.CMD_ROLE_ENTER, b.HANDLE_ROLE_ENTER)
}

func (b *Balcony) HANDLE_ROLE_ENTER(header *msg.MessageHeader, buffer []byte) {
	roleEnterResp := &msg.RoleEnterResp{}
	defer node.Net.Responce(header, roleEnterResp)
	roleEnterReq := &msg.RoleEnterReq{}
	err := proto.Unmarshal(buffer, roleEnterReq)
	if err != nil {
		logger.Info(err)
		return
	}

	loadRoleReq := &msg.LoadRoleReq{
		RoleUUID: roleEnterReq.RoleUUID,
	}
	loadRoleResp := &msg.LoadRoleResp{}
	if err := node.Net.Request(msg.CMD_GD_LOAD_ROLE_REQ, loadRoleReq, loadRoleResp); err != nil {
		logger.Info(err)
		return
	}
	logger.Info(fmt.Sprintf("加载角色 %s", roleEnterReq.RoleUUID))
	roleEnterResp.Success = loadRoleResp.Success
	roleEnterResp.Role = loadRoleResp.Role
	if loadRoleResp.Success {
		getMapIDReq := &msg.GetMapIDReq{
			MapUUID: loadRoleResp.MapUUID,
		}
		getMapIDResp := &msg.GetMapIDResp{}
		err := node.Net.RequestServer(msg.CMD_GET_MAP_ID, getMapIDReq, getMapIDResp, loadRoleResp.WorldID)
		if err != nil {
			logger.Info(err)
			return
		}
		roleEnterResp.MapID = getMapIDResp.MapID
		roleEnterResp.WorldID = loadRoleResp.WorldID

		updateRoleAddressNTF := &msg.UpdateRoleAddressNTF{
			RoleUUID: roleEnterReq.RoleUUID,
			Address: &data.RoleAddress{
				Gate:  global.InvalidServerID,
				World: loadRoleResp.WorldID,
			},
		}
		node.Net.Notify(msg.CMD_UPDATE_ROLE_ADDRESS_NTF, updateRoleAddressNTF)

		roleEnterNTF := &msg.RoleEnterNTF{
			Gate: header.From,
			Role: loadRoleResp.Role,
		}
		node.Net.NotifyServer(msg.CMD_ROLE_ENTER, roleEnterNTF, loadRoleResp.WorldID)
	}
}
