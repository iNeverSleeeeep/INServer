package gate

import (
	"INServer/src/common/global"
	"INServer/src/common/util"
	"INServer/src/proto/data"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const INT_MAX int = int(^uint(0) >> 1)

type (
	session struct {
		conn    *net.TCPConn
		webconn *websocket.Conn
		info    *data.PlayerSessionInfo
		cert    *data.SessionCertData
	}
)

func newSession(conn *net.TCPConn, uuid string) *session {
	s := &session{
		conn: conn,
		info: &data.PlayerSessionInfo{
			UUID: uuid,
			Address: &data.PlayerAddress{
				Gate:   global.CurrentServerID,
				Entity: -1,
			},
			State: data.SessionState_Connected,
		},
		cert: &data.SessionCertData{
			Key:           util.GetRandomString(global.CERT_KEY_LEN),
			OutOfDateTime: generateCertOutOfDateTime(),
		},
	}
	return s
}

func (s *session) generateNewCertKey() {
	s.cert = &data.SessionCertData{
		Key:           util.GetRandomString(global.CERT_KEY_LEN),
		OutOfDateTime: generateCertOutOfDateTime(),
	}
	// TODO 发送到客户端
}

func generateCertOutOfDateTime() int64 {
	return time.Now().Unix() + global.CurrentServerConfig.GateConfig.OutOfDateTimeout
}
