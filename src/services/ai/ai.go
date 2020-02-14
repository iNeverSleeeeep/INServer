package ai

// Instance AI单例
var Instance *AI

// AI 负责游戏内怪物/NPC等的行为
type AI struct {
}

// New 构造AI服务
func New() *AI {
	a := new(AI)
	return a
}

// Start 启动AI服务
func (a *AI) Start() {

}
