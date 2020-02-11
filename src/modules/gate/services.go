package gate

import (
	"INServer/src/modules/node"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"

	"github.com/golang/protobuf/proto"
)

type (
	services struct {
	}
)

func newServices() *services {
	s := new(services)
	return s
}

func (s *services) start() {
	node.Instance.Net.Listen(msg.CMD_SESSION_CERT_NTF, s.onSessionCert)
}

func (s *services) onSessionCert(header *msg.MessageHeader, buffer []byte) {
	message := &msg.LoginToGate{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		return
	}
	UUID := message.Cert.UUID
	player, ok := Instance.players[UUID]
	if ok == false {
		player = newSession(nil, UUID)
		player.info.State = data.SessionState_Offline
		Instance.players[UUID] = player
	}
	player.cert.Key = message.Cert.Key
	player.cert.OutOfDateTime = generateCertOutOfDateTime()
}
