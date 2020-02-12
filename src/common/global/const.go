package global

const (
	// CenterID 中心服务器ID
	CenterID int32 = 0

	// InvalidServerID 无效服务器ID
	InvalidServerID int32 = -1
	// InvalidServer 无效服务器类型
	InvalidServer string = "InvalidServer"

	// CenterServer 中心服务器
	CenterServer string = "Center"
	// GateServer 门服务器
	GateServer string = "Gate"
	// LoginServer 登录服务器
	LoginServer string = "Login"
	// ChatServer 聊天服务器
	ChatServer string = "Chat"
	// GPSServer 定位服务器
	GPSServer string = "GPS"
	// WorldServer 世界服务器
	WorldServer string = "World"
	// DatabaseServer 数据库服务器
	DatabaseServer string = "Database"
	// WebServer Web服务器
	WebServer string = "Web"
	// BalconyServer 月台服务器
	BalconyServer string = "Balcony"

	// ZoneStateOpen 游戏区开启
	ZoneStateOpen string = "Open"
	// ZoneStateClosed 游戏区关闭
	ZoneStateClosed string = "Closed"

	// DatabaseSchema 数据库名
	DatabaseSchema string = "indb"

	// SERVER_ID_MAX 服务器ID最大值
	SERVER_ID_MAX int32 = 999
	// CERT_KEY_LEN 客户端登录秘钥长度
	CERT_KEY_LEN int = 10
)
