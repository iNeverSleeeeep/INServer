package gate

import (
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"INServer/src/services/node"

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
	role, ok := Instance.roles[UUID]
	if ok == false {
		role = newSession(nil, UUID)
		role.info.State = data.SessionState_Offline
		Instance.roles[UUID] = role
	}
	role.cert.Key = message.Cert.Key
	role.cert.OutOfDateTime = generateCertOutOfDateTime()
}
