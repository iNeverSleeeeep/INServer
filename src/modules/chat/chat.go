package chat

import (
	"INServer/src/modules/node"
	"INServer/src/proto/msg"
)

type (
	Chat struct {
	}
)

func New() *Chat {
	c := new(Chat)

	return c
}

func (c *Chat) Start() {
	node.Instance.Net.Listen(msg.Command_CCHAT_CHAT, c.onClientChatMessage)
}

func (c *Chat) onClientChatMessage(eader *msg.MessageHeader, buffer []byte) {

}
