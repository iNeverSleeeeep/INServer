package global

// 这里类里面是避免一些循环引用问题，所以使用接口来调用方法

// IRoleGateGetter 取得角色所在的门服务器
type IRoleGateGetter interface {
	GetRoleGate(uuid string) int32
}
