package balcony

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
}
