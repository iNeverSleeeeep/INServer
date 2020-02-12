package chat

import (
	"INServer/src/proto/msg"
	"INServer/src/services/node"
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
	node.Instance.Net.Listen(msg.CMD_CCHAT_CHAT, c.onClientChatMessage)
}

func (c *Chat) onClientChatMessage(eader *msg.MessageHeader, buffer []byte) {

}
