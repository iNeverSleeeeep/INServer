package global

import (
	"INServer/src/proto/etc"
)

var (
	// CurrentServerID 当前服务器ID
	CurrentServerID int32 = 0
	// CurrentServerType 当前服务器类型
	CurrentServerType string = InvalidServer
	// CurrentServerConfig 当前服务器配置
	CurrentServerConfig *etc.ServerConfig = nil

	// CenterIP 中心服默认IP
	CenterIP string = "127.0.0.1"

	// PendingExit 等待进程终止状态
	PendingExit bool
	// Exit 由命令行或消息引起的进程退出使用这个chan
	Exit chan bool

	// RoleGateGetter 根据角色UUID取得所在门服务器
	RoleGateGetter IRoleGateGetter
)
