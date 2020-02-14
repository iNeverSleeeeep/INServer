package robot

// Instance 机器人服务单例
var Instance *Robot

// Robot 机器人 用于压测
type Robot struct {
}

// New 构造Robot服务
func New() *Robot {
	r := new(Robot)
	return r
}

// Start 启动Robot服务
func (r *Robot) Start() {

}
